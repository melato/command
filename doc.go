/*
Package command is used to create command line interfaces (CLI),
using reflection to define command flags (options)
from the public (exported) fields of any Go struct type.

It also uses reflection to map other command line arguments (after flags)
to the arguments of the command's Run method.

It supports scalar fields (string, int, ...).

It also supports struct or struct pointer fields, whose fields are also added as flags.

The optional "name" and "usage" field tags are used to set the flag name and usage,
or to exclude a field from being used as a flag.
See demo.App for an example.

A command has a hierarchy of sub-commands.  Each sub-command can have additional flags.

Flag default values cab be specified in an optional Init() method.

Flag validation can be performed in an optional Configured() method.

command uses the Go flags package for command-line processing.
*/
package command
