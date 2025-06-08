package golang

type Style string

const (
	// StyleStandard generates a standard Go enum with const values of a custom type
	StyleStandard Style = "standard"

	// StyleUnknown is used for unrecognized styles
	StyleUnknown Style = "unknown"
)

func (s Style) String() string {
	return string(s)
}

func ParseStyle(style string) Style {
	switch style {
	case string(StyleStandard):
		return StyleStandard
	default:
		return StyleUnknown
	}
}
