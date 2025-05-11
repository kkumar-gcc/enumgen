package stages

import (
	"fmt"

	"github.com/kkumar-gcc/enumgen/src/ast"
	"github.com/kkumar-gcc/enumgen/src/compiler/symbols"
	"github.com/kkumar-gcc/enumgen/src/compiler/types"
	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
	"github.com/kkumar-gcc/enumgen/src/errors"
	"github.com/kkumar-gcc/enumgen/src/token"
)

type SymbolCollector struct {
}

func NewSymbolCollector() *SymbolCollector {
	return &SymbolCollector{}
}

func (r *SymbolCollector) Name() string {
	return "SymbolCollector"
}

func (r *SymbolCollector) Process(ctx *compiler.Context) error {
	ctx.Symbols = symbols.NewTable()
	ctx.Types = types.NewRegistry()

	for _, decl := range ctx.AST.Declarations {
		switch d := decl.(type) {
		case *ast.EnumDefinition:
			r.processEnum(ctx, d)
		}
	}

	return nil
}

func (r *SymbolCollector) processEnum(ctx *compiler.Context, enumDef *ast.EnumDefinition) {
	enumName := enumDef.Name.Name
	enumPos := enumDef.Name.Pos()

	if existing := ctx.Symbols.LookupEnum(enumName); existing != nil {
		ctx.Errors.Add(&errors.CompilationError{
			Pos:      enumPos,
			Msg:      fmt.Sprintf("duplicate enum name %s, previously defined at %v", enumName, existing.Pos),
			Severity: errors.SeverityError,
			Stage:    r.Name(),
			Filename: ctx.SourcePath,
		})
		return
	}

	docString := ""
	if enumDef.Doc != nil && len(enumDef.Doc.List) > 0 {
		docString = enumDef.Doc.String()
	}

	enumSymbol := &compiler.Symbol{
		Name:      enumName,
		Kind:      compiler.SymbolEnum,
		Node:      enumDef,
		Pos:       enumPos,
		Docstring: docString,
	}

	if err := ctx.Symbols.Define(enumSymbol); err != nil {
		ctx.Errors.Add(&errors.CompilationError{
			Pos:      enumPos,
			Msg:      fmt.Sprintf("duplicate enum name: %s", err),
			Severity: errors.SeverityError,
			Stage:    r.Name(),
			Filename: ctx.SourcePath,
		})
		return
	}

	enumScope := ctx.Symbols.EnterScope()
	enumSymbol.Scope = enumScope

	seenMembers := make(map[string]token.Position)

	for _, member := range enumDef.Members {
		memberName := member.Name.Name
		memberPos := member.Name.Pos()

		if prevPos, exists := seenMembers[memberName]; exists {
			ctx.Errors.Add(&errors.CompilationError{
				Pos: memberPos,
				Msg: fmt.Sprintf("duplicate member name %s in enum %s, previously defined at %v",
					memberName, enumName, prevPos),
				Severity: errors.SeverityError,
				Stage:    r.Name(),
				Filename: ctx.SourcePath,
			})
			continue
		}

		seenMembers[memberName] = memberPos

		docString := ""
		if member.Doc != nil && len(member.Doc.List) > 0 {
			docString = member.Doc.String()
		}

		memberSymbol := &compiler.Symbol{
			Name:      memberName,
			Kind:      compiler.SymbolEnumMember,
			Node:      member,
			Pos:       memberPos,
			Docstring: docString,
		}

		if err := enumScope.Define(memberSymbol); err != nil {
			ctx.Errors.Add(&errors.CompilationError{
				Pos:      memberPos,
				Msg:      fmt.Sprintf("duplicate member name: %s", err),
				Severity: errors.SeverityError,
				Stage:    r.Name(),
				Filename: ctx.SourcePath,
			})
		}
	}

	ctx.Symbols.ExitScope()
}
