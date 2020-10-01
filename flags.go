package command

import (
	"flag"
	"fmt"
	"reflect"
	"strings"

	"melato.org/command/reflx"
)

type commandFlag struct {
	Names  []string
	Usage  string
	Prefix *flagPrefix
	Value  flag.Value
}

func (t *commandFlag) PrimaryNameIndex() int {
	for i, name := range t.Names {
		if len(name) > 1 {
			return i
		}
	}
	return 0
}

type flagPrefix struct {
	Name  string
	Usage string
}

func (t *flagPrefix) ComposeName(name string) string {
	if t == nil || t.Name == "" {
		return name
	}
	return t.Name + "." + name
}

func (t *flagPrefix) ComposeUsage(usage string) string {
	if t == nil || t.Usage == "" {
		return usage
	} else if usage == "" {
		return t.Usage
	} else {
		return t.Usage + " " + usage
	}
}

func (t *flagPrefix) Append(p *flagPrefix) *flagPrefix {
	if t == nil {
		return p
	}
	if p.Name == "" {
		return t
	}
	var result flagPrefix
	result.Name = t.ComposeName(p.Name)
	result.Usage = t.ComposeUsage(p.Usage)
	return &result
}

func extractFlags(cmd interface{}, prefix *flagPrefix) []*commandFlag {
	var cmdType reflect.Type = reflect.TypeOf(cmd)
	//fmt.Println("extractFlags cmdType", cmdType)
	var t reflect.Type = cmdType.Elem()
	var value reflect.Value = reflect.ValueOf(cmd).Elem()
	return extractFlagsV(value, t, prefix)
}

func extractFlagsV(value reflect.Value, t reflect.Type, prefix *flagPrefix) []*commandFlag {
	var flags []*commandFlag
	n := t.NumField()
	for i := 0; i < n; i++ {
		var field reflect.StructField = t.Field(i)
		if !isExported(field.Name) {
			continue
		}
		pm := reflx.NewParserManager()
		pType := field.Type
		fValue := value.Field(i)
		//kind := field.Type.Kind()
		kind := fValue.Type().Kind()
		if kind == reflect.Slice {
			pType = field.Type.Elem()
		}

		var nameStr string = field.Tag.Get("name")
		if nameStr == "-" || nameStr == "" {
			if _, exists := field.Tag.Lookup("name"); exists {
				// tag "name" is explicitly set to exclude this field from flags
				continue
			}
		}
		var names []string
		if nameStr != "" {
			names = strings.Split(nameStr, ",")
		}
		fPrefix := &flagPrefix{Name: "", Usage: field.Tag.Get("usage")}
		if len(names) == 0 {
			names = append(names, CreateFlagName(field.Name))
		} else {
			fPrefix.Name = names[0]
		}

		//fmt.Println(i, field.Name, pType, kind, fValue)

		if kind == reflect.Struct {
			//fmt.Println("struct: "+field.Name, fValue)
			sFlags := extractFlagsV(fValue, fValue.Type(), prefix.Append(fPrefix))
			flags = append(flags, sFlags...)
			continue
		}

		if kind == reflect.Ptr {
			//fmt.Println("pointer: "+field.Name, fValue)
			if !fValue.IsNil() {
				var ptrValue interface{} = fValue.Interface()
				ptrFlags := extractFlags(ptrValue, prefix.Append(fPrefix))
				flags = append(flags, ptrFlags...)
			}
			continue
		}

		if kind == reflect.Interface {
			if !fValue.IsNil() {
				var ptrValue interface{} = fValue.Interface()
				ptrFlags := extractFlags(ptrValue, prefix.Append(fPrefix))
				flags = append(flags, ptrFlags...)
			}
			continue
		}

		var cf commandFlag
		cf.Usage = field.Tag.Get("usage")
		cf.Prefix = prefix
		parse, found := pm.GetParserT(pType)
		if found {
			if kind == reflect.Slice {
				sv := newSliceValue(&field, value.Field(i), parse)
				sv.pType = pType
				sv.DefaultUse = field.Tag.Get("default")
				cf.Value = sv
			} else {
				fv := &fieldValue{Value: value.Field(i), Parse: parse}
				fv.pType = pType
				fv.DefaultUse = field.Tag.Get("default")
				cf.Value = fv
			}
			cf.Names = names
			flags = append(flags, &cf)
		} else if len(names) > 0 || cf.Usage != "" {
			fmt.Println("no parser for " + field.Type.String())
		}
	}
	return flags
}
