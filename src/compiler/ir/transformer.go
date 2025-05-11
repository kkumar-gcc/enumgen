package ir

import (
	"github.com/kkumar-gcc/enumgen/src/ast"
	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
	"github.com/kkumar-gcc/enumgen/src/token"
)

type Transformer struct {
	ctx         *compiler.Context
	currentEnum compiler.IREnumDefinition
}

func NewTransformer(ctx *compiler.Context) *Transformer {
	return &Transformer{
		ctx: ctx,
	}
}

func (t *Transformer) VisitFile(node *ast.File) any {
	module := NewModule(t.ctx.SourcePath, "")
	module.SetSource(string(t.ctx.SourceCode))

	var enums []compiler.IREnumDefinition
	for _, decl := range node.Declarations {
		if enumDecl, ok := decl.(*ast.EnumDefinition); ok {
			if irEnum := t.VisitEnum(enumDecl); irEnum != nil {
				enums = append(enums, irEnum.(compiler.IREnumDefinition))
			}
		}
	}

	module.SetEnums(enums)
	t.ctx.IRModule = module
	return module
}

func (t *Transformer) VisitEnum(node *ast.EnumDefinition) any {
	var doc string
	if node.Doc != nil {
		doc = node.Doc.String()
	}

	// Determine types for the enum
	var keyType, valueType compiler.Type
	if node.TypeSpec != nil && len(node.TypeSpec.Types) > 0 {
		if len(node.TypeSpec.Types) == 1 {
			// Single type enum (basic literal values)
			valueType = t.resolveType(node.TypeSpec.Types[0].Name.Name)
		} else if len(node.TypeSpec.Types) >= 2 {
			// Key-value enum
			keyType = t.resolveType(node.TypeSpec.Types[0].Name.Name)
			valueType = t.resolveType(node.TypeSpec.Types[1].Name.Name)
		}
	}

	// Default to string if no type is specified
	if valueType == nil {
		valueType = t.resolveType("string")
	}

	// Convert members
	var members []compiler.IREnumMember
	for _, member := range node.Members {
		if irMember := t.VisitMember(member); irMember != nil {
			members = append(members, irMember.(compiler.IREnumMember))
		}
	}

	enum := NewEnumDefinition(
		node.Name.Name,
		doc,
		members,
		valueType,
		keyType,
		node.Pos(),
		node,
	)

	t.currentEnum = enum
	return enum
}

func (t *Transformer) VisitMember(node *ast.MemberDefinition) any {
	var doc string
	if node.Doc != nil {
		doc = node.Doc.String()
	}

	var value compiler.IRValue
	if node.Value != nil {
		if irValue := t.VisitValue(node.Value); irValue != nil {
			value = irValue.(compiler.IRValue)
		}
	}

	return NewEnumMember(
		node.Name.Name,
		doc,
		value,
		node.Pos(),
		node,
	)
}

func (t *Transformer) VisitValue(node ast.Expr) any {
	switch v := node.(type) {
	case *ast.BasicLit:
		return t.VisitLiteral(v)
	case *ast.KeyValueExpr:
		return t.VisitKeyValue(v)
	default:
		return nil
	}
}

func (t *Transformer) VisitLiteral(node *ast.BasicLit) any {
	var typeInfo compiler.Type
	switch node.Kind {
	case token.INT:
		typeInfo = t.resolveType("int")
	case token.FLOAT:
		typeInfo = t.resolveType("float")
	case token.STRING:
		typeInfo = t.resolveType("string")
	case token.TRUE, token.FALSE:
		typeInfo = t.resolveType("bool")
	default:
		typeInfo = t.resolveType("unknown")
	}

	return NewLiteral(
		node.Kind,
		node.Value,
		node.Pos(),
		typeInfo,
	)
}

func (t *Transformer) VisitKeyValue(node *ast.KeyValueExpr) any {
	var key, value compiler.IRValue

	if irKey := t.VisitValue(node.Key); irKey != nil {
		key = irKey.(compiler.IRValue)
	}

	if irValue := t.VisitValue(node.Value); irValue != nil {
		value = irValue.(compiler.IRValue)
	}

	return NewKeyValue(
		key,
		value,
		node.Pos(),
	)
}

func (t *Transformer) resolveType(name string) compiler.Type {
	if t.ctx.Types != nil {
		return t.ctx.Types.LookupType(name)
	}
	return nil
}

func (t *Transformer) Transform() compiler.IRModule {
	file := t.ctx.AST
	return t.VisitFile(file).(compiler.IRModule)
}
