package cli

import (
	"fmt"
	"time"
)

type FormatTime struct {
	NoNewline bool   `name:"n" usage:"do not append newline"`
	Layout    string `usage:"format layout"`
}

func (t *FormatTime) Init() error {
	t.Layout = "2006-01-02 15:04:05"
	return nil
}

func (t *FormatTime) Format() error {
	s := time.Now().Format(t.Layout)
	if t.NoNewline {
		fmt.Printf("%s", s)
	} else {
		fmt.Printf("%s\n", s)
	}
	return nil
}
