package color

import "fmt"

type Printer struct {
	Theme   Theme
	Enabled bool
}

func NewPrinter(theme Theme, enabled bool) *Printer {
	return &Printer{
		Theme:   theme,
		Enabled: enabled,
	}
}

var DefaultPrinter = NewPrinter(DefaultTheme(), true)

func (r *Printer) Error(s ...any) string {
	if !r.Enabled {
		return fmt.Sprint(s...)
	}
	return r.Theme.Error.Sprint(s...)
}

func (r *Printer) Warning(s ...any) string {
	if !r.Enabled {
		return fmt.Sprint(s...)
	}
	return r.Theme.Warning.Sprint(s...)
}

func (r *Printer) Info(s ...any) string {
	if !r.Enabled {
		return fmt.Sprint(s...)
	}
	return r.Theme.Info.Sprint(s...)
}

func (r *Printer) Hint(s ...any) string {
	if !r.Enabled {
		return fmt.Sprint(s...)
	}
	return r.Theme.Hint.Sprint(s...)
}

func (r *Printer) Rule(s ...any) string {
	if !r.Enabled {
		return fmt.Sprint(s...)
	}
	return r.Theme.Rule.Sprint(s...)
}

func (r *Printer) Bold(s ...any) string {
	if !r.Enabled {
		return fmt.Sprint(s...)
	}
	return r.Theme.Bold.Sprint(s...)
}

func (r *Printer) Success(s ...any) string {
	if !r.Enabled {
		return fmt.Sprint(s...)
	}
	return r.Theme.Success.Sprint(s...)
}

func Error(s ...any) string {
	return DefaultPrinter.Error(s...)
}

func Warning(s ...any) string {
	return DefaultPrinter.Warning(s...)
}

func Info(s ...any) string {
	return DefaultPrinter.Info(s...)
}

func Hint(s ...any) string {
	return DefaultPrinter.Hint(s...)
}

func Rule(s ...any) string {
	return DefaultPrinter.Rule(s...)
}

func Bold(s ...any) string {
	return DefaultPrinter.Bold(s...)
}

func Success(s ...any) string {
	return DefaultPrinter.Success(s...)
}
