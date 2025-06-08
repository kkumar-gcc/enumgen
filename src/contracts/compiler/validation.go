package compiler

import (
	"fmt"
	"strings"

	"github.com/kkumar-gcc/enumgen/pkg/color"
	"github.com/kkumar-gcc/enumgen/src/errors"
	"github.com/kkumar-gcc/enumgen/src/token"
)

type ValidationResult struct {
	Errors   []Issue
	Warnings []Issue
}

type Issue struct {
	Position token.Position
	Message  string
	Fix      string
	RuleName string
	Severity errors.Severity
	Filename string
}

func (r *ValidationResult) HasErrors() bool {
	return len(r.Errors) > 0
}

func (r *ValidationResult) HasWarnings() bool {
	return len(r.Warnings) > 0
}

func (r *ValidationResult) FormatWarnings() string {
	if len(r.Warnings) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, issue := range r.Warnings {
		sb.WriteString(issue.Format())
		sb.WriteString("\n")
	}
	return sb.String()
}

func (r *ValidationResult) FormatErrors() string {
	if len(r.Errors) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, issue := range r.Errors {
		sb.WriteString(issue.Format())
		sb.WriteString("\n")
	}
	return sb.String()
}

func (r *ValidationResult) String() string {
	var sb strings.Builder
	for _, issue := range r.Errors {
		sb.WriteString(issue.String())
		sb.WriteString("\n")
	}
	for _, issue := range r.Warnings {
		sb.WriteString(issue.String())
		sb.WriteString("\n")
	}
	return sb.String()
}

func (r Issue) String() string {
	return fmt.Sprintf("%s: %s: %s", r.Position.String(), r.Severity, r.Message)
}

func (r Issue) Format() string {
	var sb strings.Builder

	var severityPrinter func(...any) string
	if r.Severity == errors.SeverityError {
		severityPrinter = color.Error
	} else {
		severityPrinter = color.Warning
	}

	header := fmt.Sprintf("%s: %s: %s",
		color.Bold(r.Position.String()),
		severityPrinter(r.Severity),
		r.Message,
	)
	sb.WriteString(header)

	if r.RuleName != "" {
		sb.WriteString(fmt.Sprintf(" %s", color.Rule(fmt.Sprintf("[%s]", r.RuleName))))
	}

	if r.Fix != "" {
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("  %s", color.Hint("hint: "+r.Fix)))
	}

	return sb.String()
}
