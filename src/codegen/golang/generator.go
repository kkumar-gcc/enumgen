package golang

import (
	"bytes"
	"fmt"
	"github.com/kkumar-gcc/enumgen/src/version"
	"maps"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/kkumar-gcc/enumgen/pkg/strconvx"
	"github.com/kkumar-gcc/enumgen/src/codegen/contracts"
	"github.com/kkumar-gcc/enumgen/src/codegen/golang/types"
	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
)

var _ contracts.Generator = (*Generator)(nil)

type Generator struct {
	templates       map[Style]*template.Template
	valueFormatters map[string]types.ValueFormatter
}

func New() (*Generator, error) {
	templates, err := LoadTemplates(defaultTemplates, templatesFS)
	if err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	return &Generator{
		templates:       templates,
		valueFormatters: defaultValueFormatters(), // Renamed
	}, nil
}

func (g *Generator) Name() string {
	return "Go"
}

func (g *Generator) Language() string {
	return "go"
}

func (g *Generator) DefaultOptions() map[string]string {
	return defaultOptions
}

func (g *Generator) OptionHelp() string {
	sb := strings.Builder{}
	sb.WriteString("Available options for " + g.Name() + " code generation:\n")
	for key, value := range g.DefaultOptions() {
		help := optionHelp[key]
		if help == "" {
			sb.WriteString(fmt.Sprintf("  - %s (default: %s)\n", key, value))
			continue
		}
		sb.WriteString(fmt.Sprintf("  - %s: %s (default: %s)\n", key, help, value))
	}
	return sb.String()
}

func (g *Generator) Generate(module compiler.IRModule, options map[string]string) ([]*compiler.OutputFile, error) {
	opts := g.DefaultOptions()
	maps.Copy(opts, options)

	files := make([]*compiler.OutputFile, 0, len(module.Enums()))
	for _, enum := range module.Enums() {
		fileName := generateFileName(enum.Name())
		filePath := filepath.Join(fileName)

		code, err := g.generateEnum(enum, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to generate code for enum '%s': %w", enum.Name(), err)
		}

		files = append(files, &compiler.OutputFile{
			Path: filePath,
			Body: code,
		})
	}

	return files, nil
}

func (g *Generator) generateEnum(enum compiler.IREnumDefinition, options map[string]string) ([]byte, error) {
	enumStyle, _ := options[OptionEnumStyle]
	templateName := ParseStyle(enumStyle)

	data, err := g.prepareTemplateData(enum, options)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare data: %w", err)
	}

	tmpl, ok := g.templates[templateName]
	if !ok {
		return nil, fmt.Errorf("template for style '%s' not found", templateName)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.Bytes(), nil
}

func (g *Generator) prepareTemplateData(enum compiler.IREnumDefinition, options map[string]string) (*TemplateData, error) {
	valueType := enum.ValueType()
	if valueType == nil {
		return nil, fmt.Errorf("enum '%s' has no value type defined", enum.Name())
	}
	valueFormatter := g.getValueFormatter(valueType.String())
	if valueFormatter == nil {
		return nil, fmt.Errorf("unsupported value type '%s' for enum '%s'", valueType.String(), enum.Name())
	}

	keyFormatter := valueFormatter
	if keyType := enum.KeyType(); keyType != nil {
		keyFormatter = g.getValueFormatter(keyType.String())
		if keyFormatter == nil {
			return nil, fmt.Errorf("unsupported key type '%s' for enum '%s'", keyType.String(), enum.Name())
		}
	}

	members := make([]TemplateMember, 0, len(enum.Members()))
	for i, member := range enum.Members() {
		var keyIR, valueIR compiler.IRValue

		if kv, ok := member.Value().(compiler.IRKeyValue); ok {
			keyIR = kv.Key()
			valueIR = kv.Value()
		} else {
			keyIR = member.Value()
			valueIR = member.Value()
		}

		formattedKey, err := keyFormatter.FormatMemberValue(keyIR, member.Name(), i)
		if err != nil {
			return nil, fmt.Errorf("error formatting key for member '%s': %w", member.Name(), err)
		}

		formattedValue, err := valueFormatter.FormatMemberValue(valueIR, member.Name(), i)
		if err != nil {
			return nil, fmt.Errorf("error formatting value for member '%s': %w", member.Name(), err)
		}

		members = append(members, TemplateMember{
			Name:  member.Name(),
			Doc:   member.Doc(),
			Key:   formattedKey,
			Value: formattedValue,
		})
	}

	return &TemplateData{
		EDLVersion:       version.Version,
		ToolVersion:      Version,
		Package:          options[OptionPackage],
		EnumName:         enum.Name(),
		EnumDoc:          enum.Doc(),
		KeyType:          keyFormatter.GoTypeName(),
		ValueType:        valueFormatter.GoTypeName(),
		KeyZeroValue:     keyFormatter.ZeroValue(),
		ValueZeroValue:   valueFormatter.ZeroValue(),
		Members:          members,
		GenerateStringer: strconvx.ToBool(options[OptionGenerateStringer], false),
		GenerateJSON:     strconvx.ToBool(options[OptionGenerateJSON], false),
		PrefixEnumName:   strconvx.ToBool(options[OptionPrefixEnumName], false),
		GenerateMap:      strconvx.ToBool(options[OptionGenerateMap], false),
	}, nil
}

func (g *Generator) getValueFormatter(enumType string) types.ValueFormatter {
	return g.valueFormatters[enumType]
}

func defaultValueFormatters() map[string]types.ValueFormatter {
	return map[string]types.ValueFormatter{
		"char":    &types.CharFormatter{},
		"string":  &types.StringFormatter{},
		"int":     &types.IntFormatter{ConcreteGoType: "int"},
		"int8":    &types.IntFormatter{ConcreteGoType: "int8"},
		"int16":   &types.IntFormatter{ConcreteGoType: "int16"},
		"int32":   &types.IntFormatter{ConcreteGoType: "int32"},
		"rune":    &types.IntFormatter{ConcreteGoType: "int32"},
		"int64":   &types.IntFormatter{ConcreteGoType: "int64"},
		"uint":    &types.IntFormatter{ConcreteGoType: "uint"},
		"uint8":   &types.IntFormatter{ConcreteGoType: "uint8"},
		"byte":    &types.IntFormatter{ConcreteGoType: "uint8"},
		"uint16":  &types.IntFormatter{ConcreteGoType: "uint16"},
		"uint32":  &types.IntFormatter{ConcreteGoType: "uint32"},
		"uint64":  &types.IntFormatter{ConcreteGoType: "uint64"},
		"float":   &types.FloatFormatter{ConcreteGoType: "float32"},
		"float32": &types.FloatFormatter{ConcreteGoType: "float32"},
		"float64": &types.FloatFormatter{ConcreteGoType: "float64"},
		"bool":    &types.BoolFormatter{},
	}
}

func generateFileName(enumName string) string {
	return fmt.Sprintf("%s_gen.go", strings.ToLower(enumName))
}
