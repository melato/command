package command

import (
	"errors"
	"fmt"
	"reflect"
)

func isFuncCompatible(fn interface{}) error {
	fType := reflect.TypeOf(fn)
	if fType.Kind() != reflect.Func {
		return errors.New("not a function")
	}
	switch fType.NumOut() {
	case 0:
	case 1:
	// how can we check that the single return value is assignable to error?
	default:
		return errors.New("function has more than one output")
	}
	numIn := fType.NumIn()
	stringType := reflect.TypeOf("")
	for i := 0; i < numIn-1; i++ {
		if fType.In(i) != stringType {
			return errors.New(fmt.Sprintf("function input(%d) is not string", i))
		}
	}
	if numIn > 0 {
		lastType := fType.In(numIn - 1)
		if lastType != stringType && lastType != reflect.SliceOf(stringType) {
			return errors.New("last input should be string or []string")
		}
	}
	return nil
}

// wrap a function so it appears as a func that takes an array of string arguments and returns an error
// panic if this is not possible (to catch errors early, instead of waiting for the user to invoke this command).
func wrapFunc(fn interface{}) func([]string) error {
	if err := isFuncCompatible(fn); err != nil {
		panic(err)
	}
	return func(args []string) error {
		fType := reflect.TypeOf(fn)
		numIn := fType.NumIn()
		if numIn == 0 && len(args) > 0 {
			return errors.New("function takes no arguments")
		}
		if !fType.IsVariadic() && numIn != len(args) {
			return errors.New("too many arguments")
		}
		if numIn > len(args)+1 {
			return errors.New("not enough arguments")
		}
		in := make([]reflect.Value, len(args))
		for i, arg := range args {
			in[i] = reflect.ValueOf(arg)
		}
		result := reflect.ValueOf(fn).Call(in)
		switch len(result) {
		case 0:
			return nil
		case 1:
			r := result[0]
			if r.IsNil() {
				return nil
			}
			return (result[0].Interface()).(error)
		default:
			return errors.New("function has more than one output")
		}
	}
}
