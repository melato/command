package reflx

import (
	"errors"
	"reflect"
)

// ParserFunc - Parses a string and sets the result to a Value
type ParseFunc func(reflect.Value, string) error

type ParserManager interface {
	// GetParser - Get the parser for a value that has the given name and type
	// check for a named parser first, then for a typed parser
	GetParser(name string, t reflect.Type) (ParseFunc, error)
	GetParserT(t reflect.Type) (ParseFunc, bool)
	Parse(name string, value reflect.Value, s string) error

	// SetParser - define a parser for a type
	SetParser(t reflect.Type, f ParseFunc)
	// SetNamedParser - define a parser for a name, such as the field name of a struct
	SetNamedParser(name string, f ParseFunc)
}

type parserManager struct {
	typeParsers  map[reflect.Type]ParseFunc
	namedParsers map[string]ParseFunc
}

func NewParserManager() ParserManager {
	var mgr parserManager
	mgr.typeParsers = make(map[reflect.Type]ParseFunc)
	mgr.typeParsers[reflect.TypeOf("")] = ParseString
	mgr.typeParsers[reflect.TypeOf(int64(0))] = ParseInt
	mgr.typeParsers[reflect.TypeOf(int32(0))] = ParseInt
	mgr.typeParsers[reflect.TypeOf(int(0))] = ParseInt
	mgr.typeParsers[reflect.TypeOf(float64(0))] = ParseFloat
	mgr.typeParsers[reflect.TypeOf(float32(0))] = ParseFloat
	mgr.typeParsers[reflect.TypeOf(false)] = ParseBool
	mgr.namedParsers = make(map[string]ParseFunc)
	return &mgr
}

func (mgr *parserManager) SetParser(t reflect.Type, f ParseFunc) {
	mgr.typeParsers[t] = f
}

func (mgr *parserManager) SetNamedParser(name string, f ParseFunc) {
	mgr.namedParsers[name] = f
}

func (mgr *parserManager) GetParserT(t reflect.Type) (ParseFunc, bool) {
	parse, found := mgr.typeParsers[t]
	return parse, found
}

func (mgr *parserManager) GetParser(name string, t reflect.Type) (ParseFunc, error) {
	parse, found := mgr.namedParsers[name]
	if !found {
		parse, found = mgr.typeParsers[t]
	}
	if !found {
		return nil, errors.New("no parser for " + name + " (" + t.Name() + ")")
	}
	return parse, nil
}

func (mgr *parserManager) Parse(name string, value reflect.Value, s string) error {
	t := value.Type()
	parse, err := mgr.GetParser(name, t)
	if err != nil {
		return nil
	}
	return parse(value, s)
}
