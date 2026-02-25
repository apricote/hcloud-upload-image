package ui

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
)

// Developed with guidance from golang docs:
// https://github.com/golang/example/blob/32022caedd6a177a7717aa8680cbe179e1045935/slog-handler-guide/README.md

const (
	// ANSI escape codes
	reset = "\033[0m"
	bold  = "\033[1m"

	// ANSI 256-color codes (closest matches to the hex colors)
	// Primary (Purple #7D56F4) -> 99
	primaryColor = "\033[38;5;99m"
	// Success (Green #04B575) -> 36
	successColor = "\033[38;5;36m"
	// Warning (Amber #FFAA00) -> 214
	warningColor = "\033[38;5;214m"
	// Error (Red #FF5F87) -> 204
	errorColor = "\033[38;5;204m"
	// Muted (Gray #626262) -> 241
	mutedColor = "\033[38;5;241m"
	// Info (White #FAFAFA) -> 231
	infoColor = "\033[38;5;231m"
)

var (
	// Style functions that apply ANSI codes
	stepHeaderStyle = func(s string) string { return bold + primaryColor + s + reset }
	infoStyle       = func(s string) string { return bold + infoColor + s + reset }
	successStyle    = func(s string) string { return bold + successColor + s + reset }
	debugStyle      = func(s string) string { return mutedColor + s + reset }
	warnStyle       = func(s string) string { return bold + warningColor + s + reset }
	errorStyle      = func(s string) string { return bold + errorColor + s + reset }
	attrStyle       = func(s string) string { return mutedColor + s + reset }

	// Icons
	stepIcon     = "▶"
	successIcon  = "✓"
	warningIcon  = "⚠"
	errorIcon    = "✗"
	cleanupIcon  = "🧹"
	skipIcon     = "⏭"
	completeIcon = "✓"
)

type Handler struct {
	opts      HandlerOptions
	goas      []groupOrAttrs
	mu        *sync.Mutex
	out       io.Writer
	seenAttrs map[string]string // Track seen attribute key-value pairs to reduce noise
}

// HandlerOptions are a subset of [slog.HandlerOptions] that are implemented for the UI handler.
type HandlerOptions struct {
	// Level reports the minimum record level that will be logged.
	// The handler discards records with lower levels.
	// If Level is nil, the handler assumes LevelInfo.
	// The handler calls Level.Level for each record processed;
	// to adjust the minimum level dynamically, use a LevelVar.
	Level slog.Leveler

	// ReplaceAttr is called to rewrite each non-group attribute before it is logged.
	// The attribute's value has been resolved (see [Value.Resolve]).
	// If ReplaceAttr returns a zero Attr, the attribute is discarded.
	//
	// The built-in attributes with keys "time", "level", "source", and "msg"
	// are passed to this function, except that time is omitted
	// if zero, and source is omitted if AddSource is false.
	//
	// The first argument is a list of currently open groups that contain the
	// Attr. It must not be retained or modified. ReplaceAttr is never called
	// for Group attributes, only their contents. For example, the attribute
	// list
	//
	//     Int("a", 1), Group("g", Int("b", 2)), Int("c", 3)
	//
	// results in consecutive calls to ReplaceAttr with the following arguments:
	//
	//     nil, Int("a", 1)
	//     []string{"g"}, Int("b", 2)
	//     nil, Int("c", 3)
	//
	// ReplaceAttr can be used to change the default keys of the built-in
	// attributes, convert types (for example, to replace a `time.Time` with the
	// integer seconds since the Unix epoch), sanitize personal information, or
	// remove attributes from the output.
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
}

// groupOrAttrs holds either a group name or a list of [slog.Attr].
type groupOrAttrs struct {
	group string      // group name if non-empty
	attrs []slog.Attr // attrs if non-empty
}

var _ slog.Handler = &Handler{}

func NewHandler(out io.Writer, opts *HandlerOptions) *Handler {
	h := &Handler{
		out:       out,
		mu:        &sync.Mutex{},
		seenAttrs: make(map[string]string),
	}
	if opts != nil {
		h.opts = *opts
	}
	if h.opts.Level == nil {
		h.opts.Level = slog.LevelInfo
	}
	return h
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *Handler) Handle(_ context.Context, record slog.Record) error {
	buf := make([]byte, 0, 512)

	// Extract structured attributes for styling decisions
	var logType, status, event string
	var step int

	record.Attrs(func(a slog.Attr) bool {
		switch a.Key {
		case "log.type":
			logType = a.Value.String()
		case "status":
			status = a.Value.String()
		case "event":
			event = a.Value.String()
		case "step":
			step = int(a.Value.Int64())
		}
		return true
	})

	// Determine the style and icon based on structured attributes
	var styleFunc func(string) string
	var icon string
	var prefix string

	switch logType {
	case "step":
		styleFunc = stepHeaderStyle
		icon = stepIcon
		prefix = fmt.Sprintf("Step %d: ", step)
	case "cleanup":
		if record.Level == slog.LevelWarn {
			styleFunc = warnStyle
			icon = warningIcon
		} else if event == "skip" {
			styleFunc = debugStyle
			icon = skipIcon
		} else {
			styleFunc = debugStyle
			icon = cleanupIcon
		}
	case "result":
		switch status {
		case "success":
			styleFunc = successStyle
			icon = successIcon
		case "failure":
			styleFunc = errorStyle
			icon = errorIcon
		default:
			styleFunc = infoStyle
		}
	case "event":
		if event == "complete" {
			styleFunc = successStyle
			icon = completeIcon
		} else {
			styleFunc = infoStyle
		}
	default:
		// Fall back to level-based styling
		switch record.Level {
		case slog.LevelWarn:
			styleFunc = warnStyle
			icon = warningIcon
		case slog.LevelError:
			styleFunc = errorStyle
			icon = errorIcon
		case slog.LevelInfo:
			styleFunc = infoStyle
		case slog.LevelDebug:
			styleFunc = debugStyle
		default:
			styleFunc = infoStyle
		}
	}

	// Render the main message with styling
	if icon != "" {
		buf = fmt.Appendf(buf, "%s ", icon)
	}

	msg := prefix + record.Message
	styledMsg := styleFunc(msg)
	buf = append(buf, []byte(styledMsg)...)

	// Add attributes if present (filtering out our special attributes)
	if record.NumAttrs() > 0 {
		// Attributes from [WithGroup] and [WithAttrs] calls
		goas := h.goas
		if record.NumAttrs() == 0 {
			for len(goas) > 0 && goas[len(goas)-1].group != "" {
				goas = goas[:len(goas)-1]
			}
		}
		group := ""
		for _, goa := range goas {
			if goa.group != "" {
				group = goa.group
			} else {
				for _, a := range goa.attrs {
					if !h.isSpecialAttr(a.Key) {
						buf = h.appendAttr(buf, group, a)
					}
				}
			}
		}

		record.Attrs(func(a slog.Attr) bool {
			if !h.isSpecialAttr(a.Key) {
				buf = h.appendAttr(buf, group, a)
			}
			return true
		})
	}

	buf = append(buf, '\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.out.Write(buf)
	return err
}

// isSpecialAttr returns true if the attribute is used for styling and shouldn't be displayed
func (h *Handler) isSpecialAttr(key string) bool {
	return key == "log.type" || key == "step" || key == "step.name" ||
		key == "event" || key == "status"
}

func (h *Handler) appendAttr(buf []byte, group string, a slog.Attr) []byte {
	a.Value = a.Value.Resolve()

	if h.opts.ReplaceAttr != nil {
		a = h.opts.ReplaceAttr([]string{group}, a)
	}

	// No-op if null attr
	if a.Equal(slog.Attr{}) {
		return buf
	}

	// Build the full key (with group if present)
	fullKey := a.Key
	if group != "" {
		fullKey = group + "." + a.Key
	}

	// Check if we've seen this attribute with the same value before
	valueStr := a.Value.String()
	if prevValue, seen := h.seenAttrs[fullKey]; seen && prevValue == valueStr {
		// Skip attributes that haven't changed
		return buf
	}

	// Update the seen attributes map
	h.seenAttrs[fullKey] = valueStr

	if group != "" {
		group += "."
	}

	// Format the attribute with styling
	var attrStr string
	switch a.Value.Kind() {
	case slog.KindString:
		attrStr = fmt.Sprintf(" %s%s=%q", group, a.Key, a.Value)
	case slog.KindAny:
		if err, ok := a.Value.Any().(error); ok {
			attrStr = fmt.Sprintf(" %s%s=%q", group, a.Key, err.Error())
		} else {
			attrStr = fmt.Sprintf(" %s%s=%s", group, a.Key, a.Value)
		}
	default:
		attrStr = fmt.Sprintf(" %s%s=%s", group, a.Key, a.Value)
	}

	// Apply styling to the attribute
	styledAttr := attrStyle(attrStr)
	buf = append(buf, []byte(styledAttr)...)

	return buf
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	return h.withGroupOrAttrs(groupOrAttrs{attrs: attrs})
}

func (h *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return h.withGroupOrAttrs(groupOrAttrs{group: name})
}

func (h *Handler) withGroupOrAttrs(goa groupOrAttrs) *Handler {
	h2 := *h
	h2.goas = make([]groupOrAttrs, len(h.goas)+1)
	copy(h2.goas, h.goas)
	h2.goas[len(h2.goas)-1] = goa
	return &h2
}
