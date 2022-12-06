package cli

import (
	"fmt"
	"strings"
)

type Strings struct {
	Separator string `name:"sep" usage:"separator"`
}

func (t *Strings) Init() error {
	t.Separator = ","
	return nil
}

func (t *Strings) Join(args []string) {
	fmt.Printf("%s\n", strings.Join(args, t.Separator))
}

func (t *Strings) Split(s string) {
	parts := strings.Split(s, t.Separator)
	for _, part := range parts {
		fmt.Printf("%s\n", part)
	}
}

func (t *Strings) Sprintf(format string, args ...string) {
	vargs := make([]any, len(args))
	for i, arg := range args {
		vargs[i] = arg
	}
	s := fmt.Sprintf(format, vargs...)
	fmt.Printf("%s\n", s)
}
