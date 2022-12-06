package cli

import (
	"fmt"
	"strconv"
)

type Sprintf struct {
	Convert ConvertArgs // flags from a struct field
	Trace   bool        `usage:"print debug info"`
}

var TraceConvert bool

type ConvertArgs struct {
	Convert string `name:"c" usage:"conversion flags"`
}

func (t *Sprintf) Configured() error {
	if t.Trace {
		TraceConvert = true
	}
	return nil
}

func (t *Sprintf) Sprintf(format string, args ...string) error {
	vargs, err := t.Convert.Args(args)
	if err != nil {
		return err
	}
	if t.Trace {
		fmt.Printf("fmt.Sprintf %s %v\n", format, vargs)
		for i, arg := range vargs {
			fmt.Printf("arg[%d]: %v (%T)\n", i, arg, arg)
		}
	}
	s := fmt.Sprintf(format, vargs...)
	fmt.Printf("%s\n", s)
	return nil
}

func (t *ConvertArgs) Args(args []string) ([]any, error) {
	vargs := make([]any, len(args))
	for i, arg := range args {
		if i >= len(t.Convert) {
			vargs[i] = arg
		} else {
			flag := t.Convert[i]
			switch flag {
			case 'i':
				if TraceConvert {
					fmt.Printf("convert %s to int64\n", arg)
				}
				d, err := strconv.ParseInt(arg, 10, 64)
				if err != nil {
					return nil, err
				}
				vargs[i] = d
			case 'f':
				if TraceConvert {
					fmt.Printf("convert %s to float64\n", arg)
				}
				f, err := strconv.ParseFloat(arg, 64)
				if err != nil {
					return nil, err
				}
				vargs[i] = f
			case 's':
				vargs[i] = arg
			default:
				return nil, fmt.Errorf("invalid convertion flag: %s.  Valid flags are: i => int64), f => float64, s => string ")
			}
		}
	}
	return vargs, nil
}
