package reflx

import (
	"reflect"
	"strconv"
)

func ParseString(value reflect.Value, s string) error {
	value.SetString(s)
	return nil
}

func ParseInt(value reflect.Value, s string) error {
	d, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	value.SetInt(d)
	return nil
}

func ParseFloat(value reflect.Value, s string) error {
	d, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	value.SetFloat(d)
	return nil
}

func ParseBool(value reflect.Value, s string) error {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	value.SetBool(b)
	return nil
}
