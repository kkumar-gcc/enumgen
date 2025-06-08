package color

type Theme struct {
	Error   Color
	Warning Color
	Info    Color
	Hint    Color
	Rule    Color
	Bold    Color
	Success Color
}

func DefaultTheme() Theme {
	return Theme{
		Error:   Red,
		Warning: Yellow,
		Info:    Blue,
		Hint:    Cyan,
		Rule:    Magenta,
		Bold:    BoldColor,
		Success: Blue,
	}
}
