package cli

import (
	"errors"
	"fmt"
)

// Flags of various types
type Flags struct {
	S       string `name:"s,s-flag" usage:"string flag with two names"`
	B       bool   // bool flag with default name, no usage
	IntFlag int    `name:"i" usage:"int flag"`
	Float   `name:"f" usage:"aliased float32"`
	Sub1    EmbeddedType  // struct
	Sub2    *EmbeddedType `name:"sub2" usage:"sub-2: "` // pointer to struct, with prefix name
	Sub3    EmbeddedType  `name:"-"`                    // no flags
}

type Float float32

type EmbeddedType struct {
	X string `usage:"X"`
	Y string `usage:"Y"`
}

type AdditionalFlags struct {
	Types  *Flags `name:"-"` // flags are specified by the parent command
	Prefix string
}

// Initialize some flags (optional).
func (t *Flags) Init() error {
	t.S = "s-default"
	t.Sub2.Y = "y-default"
	return nil
}

// Check if the flags defined by the user are ok, before running the command.
func (t *Flags) Configured() error {
	if t.S == "" {
		return errors.New("missing -s")
	}
	return nil
}

func (t *Flags) PrintFlags() error {
	fmt.Printf("s: %s\n", t.S)
	fmt.Printf("b: %v\n", t.B)
	fmt.Printf("i: %d\n", t.IntFlag)
	fmt.Printf("f: %f\n", t.Float)
	fmt.Printf("sub.x: %s\n", t.Sub1.X)
	return nil
}

func (t *AdditionalFlags) Init() error {
	t.Prefix = "hello"
	return nil
}

func (t *AdditionalFlags) Run(args []string) error {
	for _, arg := range args {
		fmt.Printf("%s: %s\n", t.Prefix, arg)
	}
	return nil
}
