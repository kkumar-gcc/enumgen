package codegen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kkumar-gcc/enumgen/src/codegen/contracts"
)

type Registry struct {
	generators map[string]contracts.Generator
}

func NewRegistry() *Registry {
	return &Registry{
		generators: make(map[string]contracts.Generator),
	}
}

func (r *Registry) Register(generator contracts.Generator) {
	r.generators[generator.Language()] = generator
}

func (r *Registry) Get(language string) (contracts.Generator, error) {
	lang := strings.ToLower(language)

	if generator, ok := r.generators[lang]; ok {
		return generator, nil
	}

	return nil, fmt.Errorf("unsupported language: %s", language)
}

func (r *Registry) Languages() []string {
	languages := make([]string, 0, len(r.generators))

	for lang := range r.generators {
		languages = append(languages, lang)
	}

	sort.Strings(languages)
	return languages
}

func (r *Registry) PrintLanguageOptions() string {
	var sb strings.Builder

	sb.WriteString("Available languages:\n")

	languages := r.Languages()
	for _, lang := range languages {
		generator, _ := r.Get(lang)

		sb.WriteString(fmt.Sprintf("- %s (%s)\n", generator.Name(), lang))

		options := generator.GetDefaultOptions()
		if len(options) > 0 {
			sb.WriteString("  Options:\n")
			sb.WriteString(generator.OptionHelp())
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
