package cli

import (
	"fmt"
	"math"
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

func (t *FormatFloat) Format(values ...float64) error {
	if len(values) == 0 {
		t.print(t.Value)
	} else {
		for _, v := range values {
			t.print(v)
		}
	}
	return nil
}
