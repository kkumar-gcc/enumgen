package rules

import (
	"fmt"

	"github.com/kkumar-gcc/enumgen/src/ast"
	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
	"github.com/kkumar-gcc/enumgen/src/errors"
	"github.com/kkumar-gcc/enumgen/src/token"
)

type TypeCompatibilityRule struct{}

func NewTypeCompatibilityRule() *TypeCompatibilityRule {
	return &TypeCompatibilityRule{}
}

func (r *TypeCompatibilityRule) Name() string {
	return "TypeCompatibilityRule"
}

func (r *TypeCompatibilityRule) Check(ctx *compiler.Context, node ast.Node) []compiler.Issue {
	enumDef, ok := node.(*ast.EnumDefinition)
	if !ok {
		return nil
	}

	declared := r.declaredTypes(enumDef)
	used := make(map[string]struct{})
	var issues []compiler.Issue

	for _, member := range enumDef.Members {
		switch expr := member.Value.(type) {
		case *ast.KeyValueExpr:
			issues = append(issues, r.checkKeyValue(expr, declared, used)...) // key:value pairs

		case *ast.BasicLit:
			issues = append(issues, r.checkLiteralMember(expr, declared, used)...) // single literal

		default:
			issues = append(issues, r.newError(member.Pos(),
				"unsupported enum member value type",
				"use a basic literal or key-value expression matching declared types"))
		}
	}

	return issues
}

func (r *TypeCompatibilityRule) declaredTypes(def *ast.EnumDefinition) []string {
	if def.TypeSpec == nil {
		return nil
	}
	names := make([]string, len(def.TypeSpec.Types))
	for i, t := range def.TypeSpec.Types {
		names[i] = t.Name.Name
	}
	return names
}

func (r *TypeCompatibilityRule) checkKeyValue(expr *ast.KeyValueExpr, declared []string, used map[string]struct{}) []compiler.Issue {
	pos := expr.Pos()
	if len(declared) != 2 {
		return []compiler.Issue{r.newError(pos,
			"enum with key-value members must declare exactly two types",
			"declare two types for key-value enum")}
	}
	var issues []compiler.Issue

	issues = append(issues, r.checkLiteral(expr.Key, declared[0], expr.Key.Pos(), fmt.Sprintf("key literal must be type %s", declared[0]), fmt.Sprintf("use literal type %s", declared[0]))...)
	issues = append(issues, r.checkLiteral(expr.Value, declared[1], expr.Value.Pos(), fmt.Sprintf("value literal must be type %s", declared[1]), fmt.Sprintf("use literal type %s", declared[1]))...)

	if lit, ok := expr.Key.(*ast.BasicLit); ok {
		if _, seen := used[lit.Value]; seen {
			issues = append(issues, r.newError(lit.Pos(),
				"duplicate enum key literal",
				"ensure each key literal is unique"))
		} else {
			used[lit.Value] = struct{}{}
		}
	}
	return issues
}

func (r *TypeCompatibilityRule) checkLiteralMember(lit *ast.BasicLit, declared []string, used map[string]struct{}) []compiler.Issue {
	pos := lit.Pos()
	if len(declared) != 1 {
		return []compiler.Issue{r.newError(pos,
			"enum with simple literals must declare exactly one type",
			"declare one type for simple-literal enum")}
	}
	var issues []compiler.Issue

	issues = append(issues, r.checkLiteral(lit, declared[0], pos, fmt.Sprintf("literal must be type %s", declared[0]), fmt.Sprintf("use literal type %s", declared[0]))...)

	if _, seen := used[lit.Value]; seen {
		issues = append(issues, r.newError(pos,
			"duplicate enum literal",
			"ensure each literal is unique"))
	} else {
		used[lit.Value] = struct{}{}
	}
	return issues
}

func (r *TypeCompatibilityRule) checkLiteral(expr ast.Expr, expectedType string, exprPos token.Position, msg, fix string) []compiler.Issue {
	lit, ok := expr.(*ast.BasicLit)
	if !ok {
		return []compiler.Issue{r.newError(exprPos, msg, fix)}
	}
	actual := literalType(lit.Kind)
	if actual != expectedType {
		return []compiler.Issue{r.newError(exprPos,
			fmt.Sprintf("literal type '%s' does not match expected '%s'", actual, expectedType),
			fix)}
	}
	return nil
}

func (r *TypeCompatibilityRule) newError(pos token.Position, msg, fix string) compiler.Issue {
	return compiler.Issue{
		Position: pos,
		Message:  msg,
		Fix:      fix,
		RuleName: r.Name(),
		Severity: errors.SeverityError,
	}
}

func literalType(kind token.Token) string {
	switch kind {
	case token.INT:
		return "int"
	case token.FLOAT:
		return "float"
	case token.CHAR, token.STRING:
		return "string"
	case token.TRUE, token.FALSE:
		return "bool"
	default:
		return "unknown"
	}
}
