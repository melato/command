package cli

import (
	"fmt"
	"math"
	"strconv"
)

// FormatFloat demonstrates flags and variable arguments
type FormatFloat struct {
	Value float64 `name:"v" usage:"default value to format"`
	Fmt   string  `name:"f" usage:"fmt format string"`
}

func (t *FormatFloat) Init() error {
	t.Fmt = "%f"
	t.Value = math.Pi
	return nil
}

func (t *FormatFloat) print(v float64) {
	s := fmt.Sprintf(t.Fmt, v)
	fmt.Printf("%s\n", s)
}

func (t *FormatFloat) Format(args []string) error {
	if len(args) == 0 {
		t.print(t.Value)
	} else {
		for _, arg := range args {
			v, err := strconv.ParseFloat(arg, 64)
			if err != nil {
				return err
			}
			t.print(v)
		}
	}
	return nil
}
