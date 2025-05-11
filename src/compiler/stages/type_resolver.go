package stages

import (
	"fmt"

	"github.com/kkumar-gcc/enumgen/src/ast"
	"github.com/kkumar-gcc/enumgen/src/compiler/types"
	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
	"github.com/kkumar-gcc/enumgen/src/errors"
)

type TypeResolver struct{}

func NewTypeResolver() *TypeResolver {
	return &TypeResolver{}
}

func (r *TypeResolver) Name() string {
	return "TypeResolver"
}

func (r *TypeResolver) Process(ctx *compiler.Context) error {
	if ctx == nil {
		return fmt.Errorf("context is nil")
	}

	if ctx.Types == nil {
		ctx.Types = types.NewRegistry()
	}

	stringType := ctx.Types.LookupType("string")
	if stringType == nil {
		stringType = types.NewType(compiler.TypePrimitive, "string", nil)
		err := ctx.Types.RegisterType(stringType)
		if err != nil {
			return fmt.Errorf("failed to register string type: %w", err)
		}
	}

	for _, decl := range ctx.AST.Declarations {
		enumDecl, ok := decl.(*ast.EnumDefinition)
		if !ok {
			continue
		}

		enumName := enumDecl.Name.Name
		enumSymbol := ctx.Symbols.LookupEnum(enumName)
		if enumSymbol == nil {
			ctx.Errors.Add(&errors.CompilationError{
				Pos:      enumDecl.Name.Pos(),
				Msg:      fmt.Sprintf("enum %s not found in symbol table", enumName),
				Severity: errors.SeverityWarning,
				Stage:    r.Name(),
				Filename: ctx.SourcePath,
			})
			continue
		}

		enumType := types.NewType(compiler.TypeEnum, enumName, enumDecl)
		if enumDecl.TypeSpec != nil {
			if len(enumDecl.TypeSpec.Types) == 0 {
				ctx.Errors.Add(&errors.CompilationError{
					Pos:      enumDecl.TypeSpec.Pos(),
					Msg:      "value type is required in enum type specification",
					Severity: errors.SeverityError,
					Stage:    r.Name(),
					Filename: ctx.SourcePath,
				})
				continue
			}

			if len(enumDecl.TypeSpec.Types) > 1 {
				// First type is the key type in map-like enums
				keyTypeRef := enumDecl.TypeSpec.Types[0]
				keyTypeName := keyTypeRef.Name.Name
				keyType := r.resolveTypeRef(ctx, keyTypeRef)
				if keyType == nil {
					ctx.Errors.Add(&errors.CompilationError{
						Pos:      keyTypeRef.Pos(),
						Msg:      fmt.Sprintf("unknown key type: %s", keyTypeName),
						Severity: errors.SeverityError,
						Stage:    r.Name(),
						Filename: ctx.SourcePath,
					})
					continue
				}
				enumType.SetKeyType(keyType)

				// Second type is the value type
				valueTypeRef := enumDecl.TypeSpec.Types[1]
				valueTypeName := valueTypeRef.Name.Name
				valueType := r.resolveTypeRef(ctx, valueTypeRef)
				if valueType == nil {
					ctx.Errors.Add(&errors.CompilationError{
						Pos:      valueTypeRef.Pos(),
						Msg:      fmt.Sprintf("unknown value type: %s", valueTypeName),
						Severity: errors.SeverityError,
						Stage:    r.Name(),
						Filename: ctx.SourcePath,
					})
					continue
				}
				enumType.SetValueType(valueType)
			} else {
				valueTypeRef := enumDecl.TypeSpec.Types[0]
				valueTypeName := valueTypeRef.Name.Name
				valueType := r.resolveTypeRef(ctx, valueTypeRef)
				if valueType == nil {
					ctx.Errors.Add(&errors.CompilationError{
						Pos:      valueTypeRef.Pos(),
						Msg:      fmt.Sprintf("unknown value type: %s", valueTypeName),
						Severity: errors.SeverityError,
						Stage:    r.Name(),
						Filename: ctx.SourcePath,
					})
					continue
				}
				enumType.SetValueType(valueType)
			}
		} else {
			enumType.SetValueType(stringType)
		}

		enumSymbol.Type = enumType
		if err := ctx.Types.RegisterType(enumType); err != nil {
			ctx.Errors.Add(&errors.CompilationError{
				Pos:      enumDecl.Pos(),
				Msg:      fmt.Sprintf("failed to register enum type: %s", err),
				Severity: errors.SeverityError,
				Stage:    r.Name(),
				Filename: ctx.SourcePath,
			})
		}
	}

	return nil
}

func (r *TypeResolver) resolveTypeRef(ctx *compiler.Context, typeRef *ast.TypeRef) compiler.Type {
	if typeRef == nil {
		return nil
	}

	typeName := typeRef.Name.Name
	resolvedType := ctx.Types.LookupType(typeName)

	if resolvedType == nil && types.IsPrimitiveType(typeName) {
		primitiveType := types.NewType(compiler.TypePrimitive, typeName, nil)
		if err := ctx.Types.RegisterType(primitiveType); err == nil {
			resolvedType = primitiveType
		}
	}

	return resolvedType
}
