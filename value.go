package command

import (
	"fmt"
	"reflect"

	"melato.org/command/reflx"
)

type fieldValue struct {
	Value      reflect.Value
	Parse      reflx.ParseFunc
	pType      reflect.Type
	DefaultUse string
}

func (t *fieldValue) IsBoolFlag() bool {
	return t.pType == reflect.TypeOf(false)
}

func (t *fieldValue) isString() bool {
	return t.pType == reflect.TypeOf("")
}

func (t *fieldValue) String() string {
	s := fmt.Sprint(t.Value)
	if s == "" {
		if t.DefaultUse != "" {
			return t.DefaultUse
		} else {
			return `""`
		}
	} else if t.isString() {
		return quote(s)
	} else {
		return s
	}
}

func (t *fieldValue) Set(s string) error {
	return t.Parse(t.Value, s)
}

type sliceValue struct {
	Value      reflect.Value
	Parse      reflx.ParseFunc
	pType      reflect.Type
	DefaultUse string

	eval  reflect.Value
	isSet bool
}

func (t *sliceValue) IsBoolFlag() bool {
	return t.pType == reflect.TypeOf(false)
}

func (t *sliceValue) String() string {
	s := fmt.Sprint(t.Value)
	if s == "[]" {
		// I couldn't figure out how to test that t.Value refers to an empty slice
		if t.DefaultUse != "" {
			return t.DefaultUse
		}
	}
	return s
}

func (t *sliceValue) Set(s string) error {
	err := t.Parse(t.eval, s)
	if err != nil {
		return err
	}
	if !t.isSet {
		t.Value.SetLen(0)
		t.isSet = true
	}
	newValue := reflect.Append(t.Value, t.eval)
	t.Value.Set(newValue)
	return nil
}

func newSliceValue(field *reflect.StructField, value reflect.Value, parse reflx.ParseFunc) *sliceValue {
	v := sliceValue{Value: value, Parse: parse}
	v.eval = reflect.New(field.Type.Elem()).Elem()
	return &v
}
