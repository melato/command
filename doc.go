// Package command imlements a command line interface that uses reflection
// to define command flags (options) from the fields of any user-specified struct.
//
// If a field is a struct or a pointer to a struct, its fields are also added as flags, and so on.
//
// The optional "name" and "usage" field tags are used to set the flag name and usage.  go doc demo.App for an example.
//
// The flag default value is any non-zero existing flag value, which can be set from an optional Init() method.
//
// A command has a hierarchy of sub-commands.  Each sub-command can have additional flags.
//
// At each level, optional Init(), Configured(), and Close() methods can do initialization, validation, and cleanup.
//
// command uses the Go flags package for command-line processing.
package command
