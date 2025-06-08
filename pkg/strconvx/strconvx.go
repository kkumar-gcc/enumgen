package strconvx

import (
	"strconv"
)

// ParseBool converts a string to a boolean.
// It is a direct wrapper around strconv.ParseBool, accepting "1", "t", "T", "true", "TRUE", "True", "0", "f", "F", "false", "FALSE", "False".
func ParseBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// ToBool converts a string to a boolean, returning a default value if the conversion fails.
func ToBool(s string, defaultValue bool) bool {
	val, err := strconv.ParseBool(s)
	if err != nil {
		return defaultValue
	}
	return val
}

// ParseInt64 converts a string to an int64.
func ParseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64) // base 10, 64-bit
}

// ToInt64 converts a string to an int64, returning a default value if the conversion fails.
func ToInt64(s string, defaultValue int64) int64 {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return defaultValue
	}
	return val
}

// ToInt converts a string to an int, returning a default value if the conversion fails.
// Note: This may truncate values on 32-bit systems if the number exceeds the 32-bit integer range.
func ToInt(s string, defaultValue int) int {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return defaultValue
	}
	return int(val)
}

// ParseFloat64 converts a string to a float64.
func ParseFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// ToFloat64 converts a string to a float64, returning a default value if the conversion fails.
func ToFloat64(s string, defaultValue float64) float64 {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultValue
	}
	return val
}

// Unquote removes surrounding quotes from a string.
// If the string is not quoted, it returns the original string.
// This is useful for parsing JSON strings or similar formats.
//
//	Unquote("\"foo\"") // "foo"
//	Unquote("foo") // "foo"
//	Unquote("\"foo") // "\"foo"
//	Unquote("foo\"") // "foo\""
func Unquote(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}
