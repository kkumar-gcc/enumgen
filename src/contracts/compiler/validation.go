package compiler

import "github.com/kkumar-gcc/enumgen/src/token"

type ValidationResult struct {
	Errors   []Issue
	Warnings []Issue
}

type Issue struct {
	Position token.Position
	Message  string
	Fix      string
	RuleName string
	Severity Severity
}

type Severity int

const (
	SeverityError Severity = iota
	SeverityWarning
	SeverityInfo
	SeverityFatal
)
