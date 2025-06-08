package strconvx

import (
	"strconv"
	"time"
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

// A list of common time layouts to try when parsing.
var commonTimeLayouts = []string{
	time.RFC3339,
	"2006-01-02T15:04:05", // ISO 8601 without timezone
	"2006-01-02 15:04:05",
	"2006-01-02",
	"02-Jan-2006",
	"01/02/2006",
	time.RFC822,
}

// ParseTime converts a string to a time.Time by trying a list of common layouts.
func ParseTime(s string) (time.Time, error) {
	for _, layout := range commonTimeLayouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, &time.ParseError{Layout: "multiple", Value: s, Message: ": could not parse time"}
}

// ToTime converts a string to a time.Time, returning a default value if parsing fails.
func ToTime(s string, defaultValue time.Time) time.Time {
	val, err := ParseTime(s)
	if err != nil {
		return defaultValue
	}
	return val
}
