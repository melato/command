package cli

import (
	"fmt"
)

type Add struct {
}

func (t *Add) Integers(a, b int) {
	fmt.Printf("%d\n", a+b)
}
