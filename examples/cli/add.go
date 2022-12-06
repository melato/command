package cli

import (
	"fmt"
)

type Add struct {
}

func (t *Add) Float(a, b float64) {
	fmt.Printf("%g\n", a+b)
}
