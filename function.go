package command

import (
	"errors"
	"reflect"

	"melato.org/command/reflx"
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
		pm := reflx.NewParserManager()
		stringType := reflect.TypeOf("")
		var aType reflect.Type
		var parse reflx.ParseFunc
		for i, arg := range args {
			if i < numIn {
				aType = fType.In(i)
			}
			if i == numIn-1 && fType.IsVariadic() {
				aType = aType.Elem()
			}
			if aType == stringType {
				in[i] = reflect.ValueOf(arg)
			} else {
				if i < numIn {
					var found bool
					parse, found = pm.GetParserT(aType)
					if !found {
						return errors.New("no parser for " + aType.String())
					}
				}
				value := reflect.Indirect(reflect.New(aType))
				parse(value, arg)
				in[i] = value
			}
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
