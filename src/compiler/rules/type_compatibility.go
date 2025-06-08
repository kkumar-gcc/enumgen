package rules

import (
	"fmt"
	goconst "go/constant"
	gotoken "go/token"
	"math"
	"math/big"

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
	if err := isFitsInTypeRange(lit, expectedType); err != nil {
		return []compiler.Issue{r.newError(exprPos,
			fmt.Sprintf("literal %s is not a valid %s: %v", lit.Value, expectedType, err),
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

func makeUntypedConst(lit *ast.BasicLit) (goconst.Value, error) {
	var goKind gotoken.Token
	switch lit.Kind {
	case token.INT:
		goKind = gotoken.INT
	case token.FLOAT:
		goKind = gotoken.FLOAT
	case token.STRING:
		goKind = gotoken.STRING
	case token.CHAR:
		goKind = gotoken.CHAR
	case token.TRUE, token.FALSE:
		return goconst.MakeBool(lit.Kind == token.TRUE), nil
	default:
		return nil, fmt.Errorf("unsupported literal kind %s", lit.Kind)
	}
	return goconst.MakeFromLiteral(lit.Value, goKind, 0), nil
}

func intRange(bitSize int) (min *big.Int, max *big.Int) {
	one := big.NewInt(1)
	two := big.NewInt(2)
	pow := new(big.Int).Exp(two, big.NewInt(int64(bitSize-1)), nil)
	max = new(big.Int).Sub(pow, one) // 2^(b-1)-1
	min = new(big.Int).Neg(pow)      // -2^(b-1)
	return
}

func uintRange(bitSize int) (min *big.Int, max *big.Int) {
	two := big.NewInt(2)
	pow := new(big.Int).Exp(two, big.NewInt(int64(bitSize)), nil) // 2^bitSize
	max = new(big.Int).Sub(pow, big.NewInt(1))                    // 2^b - 1
	min = big.NewInt(0)
	return
}

func floatRange(bitSize int) (max *big.Float, minNeg *big.Float) {
	if bitSize == 32 {
		m := new(big.Float).SetFloat64(math.MaxFloat32)
		return m, new(big.Float).Neg(m)
	}
	m := new(big.Float).SetFloat64(math.MaxFloat64)
	return m, new(big.Float).Neg(m)
}

var (
	intTypeToBitSize = map[string]int{
		"int8":  8,
		"int16": 16,
		"int32": 32,
		"int64": 64,
	}

	uintTypeToBitSize = map[string]int{
		"uint8":  8,
		"uint16": 16,
		"uint32": 32,
		"uint64": 64,
	}

	floatTypeToBitSize = map[string]int{
		"float32": 32,
		"float64": 64,
	}
)

// isFitsInTypeRange returns nil if lit can be represented exactly as expectedType.
func isFitsInTypeRange(lit *ast.BasicLit, expectedType string) error {
	litStr := lit.Value
	val, err := makeUntypedConst(lit)
	if err != nil {
		return fmt.Errorf("cannot parse literal %q: %v", litStr, err)
	}

	switch expectedType {
	case "int8", "int16", "int32", "int64":
		if val.Kind() != goconst.Int {
			return fmt.Errorf("literal %s is not an integer", litStr)
		}

		bitSize := intTypeToBitSize[expectedType]

		minInt, maxInt := intRange(bitSize)
		minConst := goconst.Make(minInt)
		maxConst := goconst.Make(maxInt)

		isTooSmall := goconst.Compare(val, gotoken.LSS, minConst)
		isTooLarge := goconst.Compare(val, gotoken.GTR, maxConst)

		if isTooSmall || isTooLarge {
			return fmt.Errorf(
				"integer %s out of range for %s (min: %s, max: %s)",
				litStr,
				expectedType,
				minInt.String(),
				maxInt.String(),
			)
		}

		return nil

	case "uint8", "uint16", "uint32", "uint64":
		if val.Kind() != goconst.Int {
			return fmt.Errorf("literal %s is not an integer", litStr)
		}

		if goconst.Sign(val) < 0 {
			return fmt.Errorf("integer %s is negative, cannot assign to %s", litStr, expectedType)
		}

		bitSize := uintTypeToBitSize[expectedType]
		_, maxBigInt := uintRange(bitSize)

		maxConst := goconst.Make(maxBigInt)

		if goconst.Compare(val, gotoken.GTR, maxConst) {
			return fmt.Errorf("integer %s out of range for %s (max: %s)", litStr, expectedType, maxBigInt.String())
		}

		return nil

	case "float32", "float64":
		if val.Kind() != goconst.Float && val.Kind() != goconst.Int {
			return fmt.Errorf("literal %s is not a numeric value", litStr)
		}

		bitSize := floatTypeToBitSize[expectedType]
		maxFloat, minNegFloat := floatRange(bitSize)

		maxConst := goconst.Make(maxFloat)
		minConst := goconst.Make(minNegFloat)

		isTooLarge := goconst.Compare(val, gotoken.GTR, maxConst)
		isTooSmall := goconst.Compare(val, gotoken.LSS, minConst)

		if isTooSmall || isTooLarge {
			return fmt.Errorf(
				"floating-point literal %s out of range for %s",
				litStr,
				expectedType,
			)
		}

		return nil

	case "true", "false", "bool":
		if lit.Kind != token.TRUE && lit.Kind != token.FALSE {
			return fmt.Errorf("literal %s is not a boolean literal", litStr)
		}

		return nil

	case "string":
		if lit.Kind != token.STRING {
			return fmt.Errorf("literal %s is not a string literal", litStr)
		}

		return nil

	case "char":
		if lit.Kind != token.CHAR {
			return fmt.Errorf("literal %s is not a character literal", litStr)
		}

		if len(lit.Value) != 3 || lit.Value[0] != '\'' || lit.Value[2] != '\'' || lit.Value[1] < 0 {
			return fmt.Errorf("literal %s is not a valid character literal", litStr)
		}

		return nil

	default:
		return fmt.Errorf("unknown or unsupported target type %q", expectedType)
	}
}
