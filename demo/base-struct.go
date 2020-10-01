package main

import (
	"fmt"
	"os"

	"melato.org/export/command"
)

/** Example of single struct that defines a command with various options
 */
type Demo struct {
	command.Base
	S     string   `name:"s" usage:"demo string"`
	T     bool     `usage:"boolean flag"`
	F     bool     `usage:"boolean flag"`
	N     int      `usage:"int flag"`
	AB    string   `name:"ab" usage:"long name"`
	Home  string   `name:"home" default:"HOME env variable" usage:"field with default tag"`
	Homes []string `name:"homes" default:"HOME env variable" usage:"slice with default tag"`
	M     string   `name:"m,multi" usage:"multiple names"`
	/** non-flag */
	P      int      `name:""` // an explicit empty name excludes this flag
	Sarray []string `name:"a" usage:"multiple values"`
}

/** Initialize values */
func (t *Demo) Init() error {
	t.S = "s-example"
	t.T = true
	t.N = 7
	t.AB = "aa"
	t.Sarray = []string{"a", "b"}
	return nil
}

func (t *Demo) Configured() error {
	if t.Home == "" {
		t.Home, _ = os.LookupEnv("HOME")
	}

	return nil
}

func (t *Demo) Run(args []string) error {
	fmt.Println("Demo", t.S, t.T, t.F, t.N, t.Sarray)
	fmt.Println("Home:", t.Home)
	fmt.Println("Homes:", t.Homes)
	fmt.Println("Args:")
	for _, arg := range args {
		fmt.Println(arg)
	}
	return nil
}

func (t *Demo) Usage() *command.Usage {
	return &command.Usage{
		Use:   "[arg]...",
		Short: "flags of various types",
		Long: `Demostrates:
- Defining a Command with flags in a single struct
- Flags of various types, including arrays
- Specifying flag names and usage through tags
- Init() - initialization of default values
`,
		Example: `demo1.go -t=false -f -n 3 -a x -a y
demo1.go -h`,
	}
}

func main() {
	command.Main(&Demo{})
}
