package errors

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kkumar-gcc/enumgen/pkg/color"
	"github.com/kkumar-gcc/enumgen/src/token"
)

type Severity int

const (
	SeverityWarning Severity = iota
	SeverityError
	SeverityFatal
	SeverityInfo
)

func (s Severity) String() string {
	switch s {
	case SeverityWarning:
		return "warning"
	case SeverityError:
		return "error"
	case SeverityFatal:
		return "fatal"
	case SeverityInfo:
		return "info"
	default:
		return "unknown"
	}
}

type CompilationError struct {
	Pos      token.Position
	Msg      string
	Fix      string
	Severity Severity
	Stage    string
	Filename string
}

func (r *CompilationError) Error() string {
	return fmt.Sprintf("%s: %s: %s", r.Pos.String(), r.Severity, r.Msg)
}

func (r *CompilationError) Format() string {
	var sb strings.Builder

	var severityPrinter func(...any) string
	switch r.Severity {
	case SeverityError, SeverityFatal:
		severityPrinter = color.Error
	case SeverityWarning:
		severityPrinter = color.Warning
	default:
		severityPrinter = color.Info
	}

	header := fmt.Sprintf("%s: %s: %s",
		color.Bold(r.Pos.String()),
		severityPrinter(r.Severity),
		r.Msg,
	)
	sb.WriteString(header)

	if r.Stage != "" {
		sb.WriteString(fmt.Sprintf(" %s", color.Rule(fmt.Sprintf("[%s]", r.Stage))))
	}

	if r.Fix != "" {
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("  %s", color.Hint("hint: "+r.Fix)))
	}

	return sb.String()
}

type ErrorList []*CompilationError

func (l *ErrorList) Add(err *CompilationError) {
	*l = append(*l, err)
}

func (l ErrorList) HasFatal() bool {
	for _, err := range l {
		if err.Severity == SeverityFatal {
			return true
		}
	}
	return false
}

func (l ErrorList) HasErrors() bool {
	for _, err := range l {
		if err.Severity == SeverityError || err.Severity == SeverityFatal {
			return true
		}
	}
	return false
}

func (l ErrorList) IsEmpty() bool {
	return len(l) == 0
}

func (l ErrorList) SortByPosition() {
	sort.Slice(l, func(i, j int) bool {
		if l[i].Filename != l[j].Filename {
			return l[i].Filename < l[j].Filename
		}

		if l[i].Pos.Line != l[j].Pos.Line {
			return l[i].Pos.Line < l[j].Pos.Line
		}

		return l[i].Pos.Column < l[j].Pos.Column
	})
}

func (l ErrorList) FilterBySeverity(severity Severity) ErrorList {
	var filtered ErrorList
	for _, err := range l {
		if err.Severity == severity {
			filtered = append(filtered, err)
		}
	}
	return filtered
}

func (l ErrorList) Format() string {
	if len(l) == 0 {
		return color.Success("✔ Validation successful.")
	}

	l.SortByPosition()
	var sb strings.Builder

	for i, err := range l {
		sb.WriteString(err.Format())
		if i < len(l)-1 {
			sb.WriteString("\n")
		}
	}

	summary := l.Summary()
	if summary != "" {
		sb.WriteString("\n")
		sb.WriteString(color.Rule(summary))
	}

	return sb.String()
}

func (l ErrorList) Summary() string {
	counts := make(map[Severity]int)
	for _, err := range l {
		counts[err.Severity]++
	}

	var parts []string
	if count := counts[SeverityFatal]; count > 0 {
		parts = append(parts, color.Error(fmt.Sprintf("%d fatal", count)))
	}
	if count := counts[SeverityError]; count > 0 {
		parts = append(parts, color.Error(fmt.Sprintf("%d errors", count)))
	}
	if count := counts[SeverityWarning]; count > 0 {
		parts = append(parts, color.Warning(fmt.Sprintf("%d warnings", count)))
	}
	if count := counts[SeverityInfo]; count > 0 {
		parts = append(parts, color.Info(fmt.Sprintf("%d info", count)))
	}

	if len(parts) == 0 {
		return color.Success("No issues found.")
	}

	prefix := color.Error("✖ Found")
	if !l.HasErrors() && l.HasFatal() {
		prefix = color.Warning("✖ Found")
	}

	return fmt.Sprintf("%s %s.", prefix, strings.Join(parts, " and "))
}

func (l ErrorList) Error() string {
	if len(l) == 0 {
		return "no errors"
	}

	l.SortByPosition()

	var messages []string
	for _, err := range l {
		messages = append(messages, err.Error())
	}

	return strings.Join(messages, "\n")
}
