package main

import (
	"fmt"

	"melato.org/export/command"
)

type Extra struct {
	E string `name:"e" usage:"flag defined in another struct"`
}

type Excluded struct {
	X string `name:"x" usage:"not a flag in Demo, but could be a flag somewhere else"`
}

/** Our root object */
type Demo struct {
	S        string    `name:"s" usage:"a flag"`
	Extra    *Extra    // flags from another struct
	Excluded *Excluded `name:""` // empty name excludes this from flags
}

func (t *Demo) Init() error {
	t.S = "s-example"
	// The Extra flags are detected after Init(), so it's fine to initialize t.Extra here
	t.Extra = &Extra{E: "xtra"}
	return nil
}

func (t *Demo) One(args []string) error {
	fmt.Println("One Demo=", t)
	fmt.Println("Args:")
	for _, arg := range args {
		fmt.Println(arg)
	}
	return nil
}

func (t *Demo) Two() error {
	fmt.Println("Two:", t.S)
	return nil
}

type DemoThree struct {
	d *Demo  // Demo flags are not repeated here, because D is not exported
	C string `name:"c" usage:"third level flag"`
}

func NewDemoThree(d *Demo) *DemoThree {
	var d3 DemoThree
	d3.d = d
	return &d3
}

func (d3 *DemoThree) Run(args []string) error {
	fmt.Println("Three", d3.d.S, d3.C)
	return nil
}

// command dependencies start here

/** Root command - defines global options */
type DemoCmd struct {
	command.Base
	D *Demo // this is the application objec
}

func (t *DemoCmd) Usage() *command.Usage {
	return &command.Usage{
		Short: "sub-commands, decoupling of commands from user structures",
		Long: `Demostrates:
- Decoupling of Options from command package, via pointers
- Sub-Commands
- Configure() of an intermediate command
- Multiple commands implemented by the same user object
`,
		Example: `
go run demo2.go -s x
demo2.go one
demo2.go one three -h
demo2.go -s sss two three -t 98 a b c
`}
}

func (t *DemoCmd) Init() error {
	return t.D.Init()
}

/** Sub-commands */
func (t *DemoCmd) Commands() map[string]command.Command {
	commands := make(map[string]command.Command)
	commands["one"] = &DemoOneCmd{D: t.D}
	commands["two"] = &DemoTwoCmd{D: t.D}
	return commands
}

type DemoOneCmd struct {
	command.Base
	D *Demo `name:""` // exclude D from flags, even though it is exported.  D flags have already been added from the parent command
}

func (t *DemoOneCmd) Usage() *command.Usage {
	return &command.Usage{
		Short: "A command"}
}

func (t *DemoOneCmd) Run(args []string) error {
	fmt.Println("DemoOne.Run()")
	return t.D.One(args)
}

type DemoTwoCmd struct {
	command.Base
	D *Demo `name:""` // exclude D from flags
}

func (t *DemoTwoCmd) Usage() *command.Usage {
	return &command.Usage{
		Short: "Another command, with sub-commands",
	}
}

func (t *DemoTwoCmd) Configured() error {
	fmt.Println("DemoTwo.Configured() s=" + t.D.S)
	return nil
}

func (t *DemoTwoCmd) Run(args []string) error {
	return t.D.Two()
}

func (t *DemoTwoCmd) Commands() map[string]command.Command {
	commands := make(map[string]command.Command)
	d3 := NewDemoThree(t.D)
	commands["three"] = &DemoThreeCmd{D3: d3}
	return commands
}

/** third-level command, with additional  options */
type DemoThreeCmd struct {
	command.Base
	D3 *DemoThree
}

func (t *DemoThreeCmd) Usage() *command.Usage {
	return &command.Usage{
		Short: "2nd-level command",
	}
}

func (cmd *DemoThreeCmd) Run(args []string) error {
	return cmd.D3.Run(args)
}

func main() {
	demo := &Demo{}
	command.Main(&DemoCmd{D: demo})
}
