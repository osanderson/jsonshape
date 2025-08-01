package jsonshape

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

type Options struct {
	Matchers []FormatMatcher
}

// TODO - json parser orders the keys.. need to use custom parsing if we want to preserve original key order
func Sanitize(value interface{}, opts Options) (interface{}, error) {
	var err error

	switch v := value.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(v))
		for k, val := range v {
			out[k], err = Sanitize(val, opts)
			if err != nil {
				return nil, fmt.Errorf("Sanitize error: %w", err)
			}
		}
		return out, nil
	case []interface{}:
		out := make([]interface{}, len(v))
		for i, item := range v {
			out[i], err = Sanitize(item, opts)
			if err != nil {
				return nil, fmt.Errorf("Sanitize error: %w", err)
			}
		}
		return out, nil
	case string:
		return sanitizeString(v, opts), nil
	case float64:
		return sanitizeNumber(v)
	case bool:
		return false, nil
	case nil:
		return nil, nil
	default:
		return nil, nil
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

func sanitizeNumber(n float64) (float64, error) {
	var err error

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

	out, err := strconv.ParseFloat(result.String(), 64)
	if err != nil {
		return 0, fmt.Errorf("ParseFloat error: %w", err)
	}

	return out, nil
}
