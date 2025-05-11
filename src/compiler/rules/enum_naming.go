package rules

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/kkumar-gcc/enumgen/src/ast"
	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
)

type EnumNaming struct {
	strict bool
}

func NewEnumNamingRule(strict bool) *EnumNaming {
	return &EnumNaming{
		strict: strict,
	}
}

func (r *EnumNaming) Name() string {
	return "EnumNaming"
}

func (r *EnumNaming) Check(_ *compiler.Context, node ast.Node) []compiler.Issue {
	issues := make([]compiler.Issue, 0)

	enumNode, ok := node.(*ast.EnumDefinition)
	if !ok {
		return issues
	}

	name := enumNode.Name.Name
	pos := enumNode.Name.Pos()

	if len(name) > 0 && !unicode.IsUpper(rune(name[0])) {
		issues = append(issues, compiler.Issue{
			Position: pos,
			Message:  fmt.Sprintf("enum name %s must begin with uppercase letter", name),
			Fix:      fmt.Sprintf("Rename to %s%s", strings.ToUpper(name[:1]), name[1:]),
			RuleName: r.Name(),
			Severity: compiler.SeverityError,
		})
	}

	if strings.Contains(name, "_") {
		severity := compiler.SeverityWarning
		if r.strict {
			severity = compiler.SeverityError
		}

		parts := strings.Split(name, "_")
		for i := range parts {
			if len(parts[i]) > 0 {
				parts[i] = strings.ToUpper(string(parts[i][0])) + parts[i][1:]
			}
		}
		suggestedName := strings.Join(parts, "")

		issues = append(issues, compiler.Issue{
			Position: pos,
			Message:  fmt.Sprintf("enum name %s should not contain underscores", name),
			Fix:      fmt.Sprintf("Consider renaming to %s", suggestedName),
			RuleName: r.Name(),
			Severity: severity,
		})
	}

	return issues
}
