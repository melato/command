# command

A Go command line interface that uses reflection to define command flags from the exported fields of any user-specified struct.

**Example**

```
package main

import (
	"fmt"

	"melato.org/command"
)

type Hello struct {
	Prefix string
}

// Init is optional.  Used to initialize flags
func (t *Hello) Init() error {
	t.Prefix = "hello "
	return nil
}

func (t *Hello) Hello(args []string) {
	for _, name := range args {
		fmt.Println(t.Prefix + name)
	}
}

func main() {
	cmd := &command.SimpleCommand{}
	hello := &Hello{}
	cmd.Flags(hello).RunFunc(hello.Hello).Use("<name>...").Short("add a greeting to a name")
	command.Main(cmd)
}

# go run hello.go world 
# go run hello.go -h
# go run hello.go --prefix "hey " world 
```
**Features**
- A flag can be any primitive Go type, an alias of a primitive type, struct, pointer to struct, slice of primitive type
- nested commands
- flag names and usage are specified by go tag comments.  If there are no comments, a default name is used
- struct fields can be excluded from flags.
- reduces dependencies from application code.  There is nothing to subclass.  Implementation of Init() and Configured() is optional
- command help (short, long, usage, examples, flag description)
- command help can be specified from yaml data
- command functions can have a variety of signatures and are called by reflection,
automatically converting command line string arguments to the appropriate function argument types
