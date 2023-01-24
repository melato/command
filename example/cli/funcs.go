package cli

import "fmt"

type Funcs struct{}

func (t *Funcs) DoError() error {
	return fmt.Errorf("here's an error")
}

func (t *Funcs) Mixed(a string, n int) {
	fmt.Printf("a=%s n=%d\n", a, n)
}

func (t *Funcs) Strings(args ...string) {
	fmt.Println(args)
}

func (t *Funcs) Floats(args ...float64) {
	fmt.Println(args)
}
