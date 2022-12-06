package cli

import (
	"fmt"
)

type Sprintf struct {
}

func (t *Sprintf) Sprintf(format string, args ...string) {
	vargs := make([]any, len(args))
	for i, arg := range args {
		vargs[i] = arg
	}
	s := fmt.Sprintf(format, vargs...)
	fmt.Printf("%s\n", s)
}
