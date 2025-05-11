package stages

import (
	"github.com/kkumar-gcc/enumgen/src/contracts/compiler"
)

type Validator struct {
	rules []compiler.Rule
}

func NewValidator(rules []compiler.Rule) *Validator {
	return &Validator{
		rules: rules,
	}
}

func (r *Validator) Name() string {
	return "Validation"
}

func (r *Validator) Process(ctx *compiler.Context) error {
	ctx.Validations = compiler.ValidationResult{
		Warnings: []compiler.Issue{},
		Errors:   []compiler.Issue{},
	}

	for _, decl := range ctx.AST.Declarations {
		for _, rule := range r.rules {
			issues := rule.Check(ctx, decl)
			for _, issue := range issues {
				if issue.Severity >= compiler.SeverityError {
					ctx.Validations.Errors = append(ctx.Validations.Errors, issue)
				} else {
					ctx.Validations.Warnings = append(ctx.Validations.Warnings, issue)
				}
			}
		}
	}

	return nil
}
