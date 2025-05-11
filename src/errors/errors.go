package errors

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kkumar-gcc/enumgen/src/token"
)

type Severity int

const (
	SeverityWarning Severity = iota
	SeverityError
	SeverityFatal
	SeverityInfo
)

type CompilationError struct {
	Pos      token.Position
	Msg      string
	Fix      string
	Severity Severity
	Stage    string
	Filename string
}

func (e *CompilationError) Error() string {
	prefix := "warning"
	if e.Severity == SeverityError {
		prefix = "error"
	} else if e.Severity == SeverityFatal {
		prefix = "fatal error"
	} else if e.Severity == SeverityInfo {
		prefix = "info"
	}

	location := e.Pos.String()
	if e.Filename != "" && !strings.Contains(location, e.Filename) {
		location = fmt.Sprintf("%s:%s", e.Filename, location)
	}

	message := fmt.Sprintf("%s: %s: %s", location, prefix, e.Msg)

	if e.Fix != "" {
		message += fmt.Sprintf("\n   └─ Suggestion: %s", e.Fix)
	}

	return message
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

func (l ErrorList) Summary() string {
	errorCount := 0
	warningCount := 0
	infoCount := 0
	fatalCount := 0

	for _, err := range l {
		switch err.Severity {
		case SeverityError:
			errorCount++
		case SeverityWarning:
			warningCount++
		case SeverityInfo:
			infoCount++
		case SeverityFatal:
			fatalCount++
		}
	}

	var parts []string
	if fatalCount > 0 {
		parts = append(parts, fmt.Sprintf("%d fatal", fatalCount))
	}
	if errorCount > 0 {
		parts = append(parts, fmt.Sprintf("%d errors", errorCount))
	}
	if warningCount > 0 {
		parts = append(parts, fmt.Sprintf("%d warnings", warningCount))
	}
	if infoCount > 0 {
		parts = append(parts, fmt.Sprintf("%d info messages", infoCount))
	}

	if len(parts) == 0 {
		return "No issues found"
	}

	return fmt.Sprintf("Found %s", strings.Join(parts, ", "))
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
