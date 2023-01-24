A Go command line interface that uses reflection to define command flags from the exported fields of any user-specified struct.

# Example

```
package main

import (
	_ "embed"

	"fmt"

	"melato.org/command"
	"melato.org/command/usage"
)

type App struct {
	S string `name:"s" usage:"example flag"`
}

func (t *App) Init() error {
	t.S = "hello"
	return nil
}

func (t *App) Run() {
	fmt.Printf("s=%s\n", t.S)
}

//go:embed usage.yaml
var usageData []byte

func main() {
	var cmd command.SimpleCommand

	var app App
	cmd.Command("run").Flags(&app).RunFunc(app.Run)

	usage.Apply(&cmd, usageData)
	command.Main(&cmd)
}
```

# Features
- A flag can be any primitive Go type, an alias of a primitive type, struct, pointer to struct, slice of primitive type
- nested commands
- flag names and usage are specified by go tag comments.  If there are no comments, a default name is used
- struct fields can be excluded from flags.
- reduces dependencies from application code.  There is nothing to subclass.  Implementation of Init() and Configured() is optional
- command help (short, long, usage, examples, flag description)
- command help can be specified from yaml data
- command functions can have a variety of signatures and are called by reflection,
automatically converting command line string arguments to the appropriate function argument types
