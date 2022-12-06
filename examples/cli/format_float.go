package cli

import (
	"fmt"
)

type FormatFloat struct {
	Fmt string `name:"f" usage:"fmt format string"`
}

func (t *FormatFloat) Init() error {
	t.Fmt = "%f"
	return nil
}

func (t *FormatFloat) Format(f float64) error {
	fmt.Printf(t.Fmt+"\n", f)
	return nil
}
