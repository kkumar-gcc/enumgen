package golang

const (
	OptionPackage          = "package"
	OptionGenerateStringer = "generate_stringer"
	OptionGenerateJSON     = "generate_json"
	OptionPrefixEnumName   = "prefix_enum_name"
	OptionGenerateMap      = "generate_map"
	OptionEnumStyle        = "enum_style"
)

type OptionDef struct {
	Key          string
	DefaultValue string
	HelpText     string
}

var allOptions = []OptionDef{
	{
		Key:          OptionPackage,
		DefaultValue: "main",
		HelpText:     "The name of the Go package for the generated file.",
	},
	{
		Key:          OptionGenerateStringer,
		DefaultValue: "true",
		HelpText:     "If true, generates a String() method mapping the enum value to its name.",
	},
	{
		Key:          OptionGenerateJSON,
		DefaultValue: "true",
		HelpText:     "If true, generates MarshalJSON and UnmarshalJSON methods.",
	},
	{
		Key:          OptionPrefixEnumName,
		DefaultValue: "false",
		HelpText:     "If true, prefixes member names with the enum type name (e.g., ColorRED).",
	},
	{
		Key:          OptionGenerateMap,
		DefaultValue: "true",
		HelpText:     "If true, generates a map of all enum members.",
	},
	{
		Key:          OptionEnumStyle,
		DefaultValue: "standard",
		HelpText:     "The style of the generated enum code ('standard', 'iota', etc.).",
	},
}

var (
	defaultOptions map[string]string
	optionHelp     map[string]string
)

func init() {
	defaultOptions = make(map[string]string)
	optionHelp = make(map[string]string)

	for _, opt := range allOptions {
		defaultOptions[opt.Key] = opt.DefaultValue
		if opt.HelpText != "" {
			optionHelp[opt.Key] = opt.HelpText
		}
	}
}
