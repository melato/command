package cli

import (
	"fmt"
	"regexp"
)

type Regexp struct {
	Pattern string `name:"r" usage:"regular expression"`
	re      *regexp.Regexp
}

func (t *Regexp) Configured() error {
	re, err := regexp.Compile(t.Pattern)
	if err != nil {
		return err
	}
	t.re = re
	return nil
}

func (t *Regexp) FindStringSubmatch(s string) {
	parts := t.re.FindStringSubmatch(s)
	for _, part := range parts {
		fmt.Printf("%s\n", part)
	}
}

func (t *Regexp) FindAllString(s string, n int) {
	matches := t.re.FindAllString(s, n)
	for _, match := range matches {
		fmt.Printf("%s\n", match)
	}
}

func (t *Regexp) Split(s string, n int) {
	parts := t.re.Split(s, n)
	for _, part := range parts {
		fmt.Printf("%s\n", part)
	}
}
