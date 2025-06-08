package golang

import (
	"embed"
	"fmt"
	"text/template"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

var defaultTemplates = map[Style]string{
	StyleStandard: "templates/standard.go.tmpl",
	StyleUnknown:  "templates/standard.go.tmpl",
}

type TemplateMember struct {
	Name  string
	Doc   string
	Key   any
	Value any
}

type TemplateData struct {
	EDLVersion       string
	ToolVersion      string
	Package          string
	EnumName         string
	EnumDoc          string
	KeyType          string
	ValueType        string
	KeyZeroValue     any
	ValueZeroValue   any
	Members          []TemplateMember
	GenerateStringer bool
	GenerateJSON     bool
	PrefixEnumName   bool
	GenerateMap      bool
}

// LoadTemplates loads the Go templates from the embedded filesystem
// and returns a map of Style to *template.Template.
func LoadTemplates(templates map[Style]string, fs embed.FS) (map[Style]*template.Template, error) {
	loadedTemplates := make(map[Style]*template.Template)

	for style, path := range templates {
		data, err := fs.ReadFile(path)
		if err != nil {
			return nil, err
		}

		tmpl, err := template.New(style.String()).Parse(string(data))
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", path, err)
		}

		loadedTemplates[style] = tmpl
	}

	return loadedTemplates, nil
}
