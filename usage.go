package command

// Usage provides documentation for command.
type Usage struct {
	// A one-line description of the command, shown in lists of commands
	Short string `yaml:"short"`

	// A generic representation of the command-line arguments, without any options, e.g. "<arg1> <arg2>"
	Use string `yaml:"use,omitempty"`

	// A longer description shown in the help for a single command
	Long string `yaml:"long,omitempty"`

	// Examples of command line invocation
	Examples []string `yaml:"examples,omitempty"`
}
