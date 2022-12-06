package cli

import (
	"fmt"
	"time"
)

type FormatTime struct {
	Layout string `usage:"layout"`
}

func (t *FormatTime) Init() error {
	t.Layout = "2006-01-02 15:04:05"
	return nil
}

func (t *FormatTime) Format() error {
	fmt.Printf("%s\n", time.Now().Format(t.Layout))
	return nil
}
