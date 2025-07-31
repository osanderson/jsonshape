package jsonshape

import (
	"math"
	"strconv"
	"strings"
	"unicode"
)

type Options struct {
	Matchers []FormatMatcher
}

// TODO - json parser orders the keys.. need to use custom parsing if we want to preserve original key order
func Sanitize(value interface{}, opts Options) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(v))
		for k, val := range v {
			out[k] = Sanitize(val, opts)
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(v))
		for i, item := range v {
			out[i] = Sanitize(item, opts)
		}
		return out
	case string:
		return sanitizeString(v, opts)
	case float64:
		return sanitizeNumber(v)
	case bool:
		return false
	case nil:
		return nil
	default:
		return nil
	}
}

func sanitizeString(s string, opts Options) string {
	for _, m := range opts.Matchers {
		if m.Match(s) {
			return m.Redact(s)
		}
	}

	var result strings.Builder

	const runeWhitelist = "-: "

	for _, r := range s {
		switch {
		case unicode.IsLetter(r):
			result.WriteRune('A')
		case unicode.IsDigit(r):
			result.WriteRune('9')
		case strings.ContainsRune(runeWhitelist, r):
			result.WriteRune(r)
		default:
			result.WriteRune('?')
		}
	}

	return result.String()
}

func sanitizeNumber(n float64) float64 {
	isNegative := n < 0
	abs := math.Abs(n)
	str := strconv.FormatFloat(abs, 'f', -1, 64)

	var result strings.Builder

	if isNegative {
		result.WriteByte('-')
	}

	for _, ch := range str {
		switch ch {
		case '.':
			result.WriteByte('.')
		default:
			result.WriteByte('9')
		}
	}

	out, _ := strconv.ParseFloat(result.String(), 64)
	// TODO - error check?

	return out
}
