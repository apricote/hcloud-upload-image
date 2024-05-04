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
	ansiClear      = "\033[0m"
	ansiBold       = "\033[1m"
	ansiBoldYellow = "\033[1;93m"
	ansiBoldRed    = "\033[1;31m"
	ansiThinGray   = "\033[2;37m"
)

type Handler struct {
	opts HandlerOptions
	goas []groupOrAttrs
	mu   *sync.Mutex
	out  io.Writer
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
		out: out,
		mu:  &sync.Mutex{},
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

	formattingPrefix := ""

	switch record.Level {
	case slog.LevelInfo:
		formattingPrefix = ansiBold
	case slog.LevelWarn:
		// Bold + Yellow
		formattingPrefix = ansiBoldYellow
	case slog.LevelError:
		// Bold + Red
		formattingPrefix = ansiBoldRed
	}

	// Print main message in formatted text
	buf = fmt.Appendf(buf, "%s%s%s", formattingPrefix, record.Message, ansiClear)

	// Add attributes in thin gray
	buf = fmt.Append(buf, ansiThinGray)

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
				buf = h.appendAttr(buf, group, a)
			}
		}
	}

	record.Attrs(func(a slog.Attr) bool {
		buf = h.appendAttr(buf, group, a)
		return true
	})

	buf = fmt.Appendf(buf, "%s\n", ansiClear)

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.out.Write(buf)
	return err
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

	if group != "" {
		group += "."
	}

	switch a.Value.Kind() {
	case slog.KindString:
		buf = fmt.Appendf(buf, " %s%s=%q", group, a.Key, a.Value)
	case slog.KindAny:
		if err, ok := a.Value.Any().(error); ok {
			buf = fmt.Appendf(buf, " %s%s=%q", group, a.Key, err.Error())
		} else {
			buf = fmt.Appendf(buf, " %s%s=%s", group, a.Key, a.Value)
		}
	default:
		buf = fmt.Appendf(buf, " %s%s=%s", group, a.Key, a.Value)
	}

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
