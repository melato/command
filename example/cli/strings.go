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
