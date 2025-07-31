package jsonshape

import (
	"regexp"
)

type FormatMatcher struct {
	Name   string
	Match  func(string) bool
	Redact func(string) string
}

var DefaultMatchers = []FormatMatcher{
	{
		Name: "email",
		Match: func(s string) bool {
			return regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`).MatchString(s)
		},
		Redact: func(_ string) string {
			return "user@example.com"
		},
	},
	{
		Name: "phone",
		Match: func(s string) bool {
			return regexp.MustCompile(`^\+?\d{8,15}$`).MatchString(s)
		},
		Redact: func(_ string) string {
			// TODO - may want to maintain original length (smae q for email)
			return "+1234567890"
		},
	},
	{
		Name: "isoDate",
		Match: func(s string) bool {
			return regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d{3})?Z$`).MatchString(s)
		},
		Redact: func(_ string) string {
			return "1970-01-D01T00:00:00.000Z"
		},
	},
}

// TODO - do we want to support other ISO Date variants?
