package rules

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/kkumar-gcc/enumgen/src/ast"
	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
	"github.com/kkumar-gcc/enumgen/src/errors"
)

type MemberNamingRule struct {
	strict bool
}

func NewMemberNamingRule(strict bool) *MemberNamingRule {
	return &MemberNamingRule{strict: strict}
}

func (r *MemberNamingRule) Name() string {
	return "EnumNaming"
}

func (r *MemberNamingRule) Check(ctx *compiler.Context, node ast.Node) []compiler.Issue {
	issues := make([]compiler.Issue, 0)

	if enumNode, ok := node.(*ast.EnumDefinition); ok {
		enumName := enumNode.Name.Name
		isEnumExported := len(enumName) > 0 && unicode.IsUpper(rune(enumName[0]))

		for _, member := range enumNode.Members {
			memberName := member.Name.Name
			pos := member.Name.Pos()

			isMemberExported := len(memberName) > 0 && unicode.IsUpper(rune(memberName[0]))
			if isEnumExported != isMemberExported {
				var msg string
				var fix string
				if isEnumExported {
					msg = fmt.Sprintf("unexported member %s of exported enum %s should be exported", memberName, enumName)
					fix = fmt.Sprintf("Rename to %s%s", strings.ToUpper(memberName[:1]), memberName[1:])
				} else {
					msg = fmt.Sprintf("exported member %s of unexported enum %s should be unexported", memberName, enumName)
					fix = fmt.Sprintf("Rename to %s%s", strings.ToLower(memberName[:1]), memberName[1:])
				}

				issues = append(issues, compiler.Issue{
					Position: pos,
					Message:  msg,
					Fix:      fix,
					RuleName: r.Name(),
					Severity: errors.SeverityError,
					Filename: ctx.SourcePath,
				})
			}

			if strings.Contains(memberName, "_") {
				severity := errors.SeverityWarning
				if r.strict {
					severity = errors.SeverityError
				}

				parts := strings.Split(memberName, "_")
				for i := range parts {
					if len(parts[i]) > 0 {
						if i == 0 && !isEnumExported {
							parts[i] = strings.ToLower(string(parts[i][0])) + parts[i][1:]
						} else {
							parts[i] = strings.ToUpper(string(parts[i][0])) + parts[i][1:]
						}
					}
				}
				suggestedName := strings.Join(parts, "")

				issues = append(issues, compiler.Issue{
					Position: pos,
					Message:  fmt.Sprintf("member name %s should not contain underscores", memberName),
					Fix:      fmt.Sprintf("Consider renaming to %s", suggestedName),
					RuleName: r.Name(),
					Severity: severity,
					Filename: ctx.SourcePath,
				})
			}

			if r.strict && strings.Contains(memberName, "_") {
				hasMixedCase := false
				for _, r := range memberName {
					if r != '_' && unicode.IsUpper(r) {
						hasMixedCase = true
						break
					}
				}

				if hasMixedCase {
					issues = append(issues, compiler.Issue{
						Position: pos,
						Message:  fmt.Sprintf("member name %s should not mix capitalization with underscores", memberName),
						Fix:      "Choose either an all_lowercase_with_underscores or CamelCase naming style",
						RuleName: r.Name(),
						Severity: errors.SeverityError,
						Filename: ctx.SourcePath,
					})
				}
			}

			if r.strict {
				if r.isDuplicateMemberNameAcrossEnums(ctx, enumNode, memberName) {
					issues = append(issues, compiler.Issue{
						Position: pos,
						Message:  fmt.Sprintf("member name %s is already used in another enum and may cause confusion", memberName),
						Fix:      fmt.Sprintf("Consider using a more specific name such as %s%s", enumName, memberName),
						RuleName: r.Name(),
						Severity: errors.SeverityError,
						Filename: ctx.SourcePath,
					})
				}
			}
		}
	}

	return issues
}

func (r *MemberNamingRule) isDuplicateMemberNameAcrossEnums(ctx *compiler.Context, currentEnum *ast.EnumDefinition, memberName string) bool {
	for _, decl := range ctx.AST.Declarations {
		if otherEnum, ok := decl.(*ast.EnumDefinition); ok && otherEnum != currentEnum {
			for _, member := range otherEnum.Members {
				if member.Name.Name == memberName {
					return true
				}
			}
		}
	}
	return false
}
