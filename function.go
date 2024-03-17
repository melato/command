package command

import (
	"errors"
	"reflect"
	"strings"

	"melato.org/command/reflx"
)

func funcUsage(fn interface{}) string {
	fType := reflect.TypeOf(fn)
	if fType.Kind() != reflect.Func {
		return ""
	}
	n := fType.NumIn()
	types := make([]string, n)
	for i := 0; i < n; i++ {
		typeName := fType.In(i).String()
		types[i] = "<" + typeName + ">"
	}
	return strings.Join(types, " ")
}

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

func buildInputs(fn interface{}, args []string) ([]reflect.Value, error) {
	fType := reflect.TypeOf(fn)
	numIn := fType.NumIn()
	if numIn == 0 && len(args) > 0 {
		return nil, errors.New("function takes no arguments")
	}
	if !fType.IsVariadic() {
		if numIn == 1 && fType.In(0) == reflect.TypeOf(args) {
			// we make an exception for a function that takes a single []string argument
			return []reflect.Value{reflect.ValueOf(args)}, nil
		}
		if numIn != len(args) {
			return nil, errors.New("wrong number of arguments")
		}
	} else {
		if numIn > len(args)+1 {
			return nil, errors.New("not enough arguments")
		}
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
				parse, found = pm.Parser(aType)
				if !found {
					return nil, errors.New("no parser for " + aType.String())
				}
			}
			value := reflect.Indirect(reflect.New(aType))
			err := parse(value, arg)
			if err != nil {
				return nil, err
			}
			in[i] = value
		}
	}
	return in, nil

}

// wrap a function so it appears as a func that takes an array of string arguments and returns an error
// panic if this is not possible (to catch errors early, instead of waiting for the user to invoke this command).
func wrapFunc(fn interface{}) func([]string) error {
	if err := isFuncCompatible(fn); err != nil {
		panic(err)
	}
	return func(args []string) error {
		in, err := buildInputs(fn, args)
		if err != nil {
			return err
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
