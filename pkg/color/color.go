package color

import "fmt"

type Color string

const (
	Reset     Color = "\033[0m"
	Red       Color = "\033[31m"
	Yellow    Color = "\033[33m"
	Blue      Color = "\033[34m"
	Magenta   Color = "\033[35m"
	Cyan      Color = "\033[36m"
	White     Color = "\033[97m"
	BoldColor Color = "\033[1m"
)

func (c Color) Sprint(s ...any) string {
	return string(c) + fmt.Sprint(s...) + string(Reset)
}

func (c Color) Sprintf(format string, s ...any) string {
	return string(c) + fmt.Sprintf(format, s...) + string(Reset)
}
