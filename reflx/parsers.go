package reflx

import (
	"reflect"
	"strconv"
)

func ParseString(value reflect.Value, s string) error {
	value.SetString(s)
	return nil
}

func ParseBool(value reflect.Value, s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	value.SetBool(v)
	return nil
}

func ParseInt(value reflect.Value, s string) error {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	value.SetInt(v)
	return nil
}

func ParseUint(value reflect.Value, s string) error {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}
	value.SetUint(v)
	return nil
}

func ParseFloat(value reflect.Value, s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	value.SetFloat(v)
	return nil
}

func ParseComplex(value reflect.Value, s string) error {
	v, err := strconv.ParseComplex(s, 64)
	if err != nil {
		return err
	}
	value.SetComplex(v)
	return nil
}
