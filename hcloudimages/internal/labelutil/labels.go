package labelutil

import "fmt"

func Merge(a, b map[string]string) map[string]string {
	result := make(map[string]string, len(a)+len(b))

	for k, v := range a {
		result[k] = v
	}
	for k, v := range b {
		result[k] = v
	}

	return result
}

func Selector(labels map[string]string) string {
	selector := make([]byte, 0, 64)
	separator := ""

	for k, v := range labels {
		selector = fmt.Appendf(selector, "%s%s=%s", separator, k, v)

		// Do not print separator on first element
		separator = ","
	}

	return string(selector)
}
