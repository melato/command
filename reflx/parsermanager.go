package reflx

import (
	"reflect"
)

// ParserFunc - Parses a string and sets the result to a Value
type ParseFunc func(reflect.Value, string) error

type ParserManager interface {
	// Parser - Get the parser for the given type
	Parser(t reflect.Type) (ParseFunc, bool)

	// SetParser - define a parser for a type
	SetParser(t reflect.Type, f ParseFunc)
}

type parserManager struct {
	typeParsers map[reflect.Type]ParseFunc
	kindParsers map[reflect.Kind]ParseFunc
}

func NewParserManager() ParserManager {
	var mgr parserManager
	mgr.typeParsers = make(map[reflect.Type]ParseFunc)
	mgr.kindParsers = make(map[reflect.Kind]ParseFunc)
	mgr.kindParsers[reflect.String] = ParseString
	mgr.kindParsers[reflect.Int64] = ParseInt
	mgr.kindParsers[reflect.Int32] = ParseInt
	mgr.kindParsers[reflect.Int16] = ParseInt
	mgr.kindParsers[reflect.Int8] = ParseInt
	mgr.kindParsers[reflect.Int] = ParseInt
	mgr.kindParsers[reflect.Uint64] = ParseUint
	mgr.kindParsers[reflect.Uint32] = ParseUint
	mgr.kindParsers[reflect.Uint16] = ParseUint
	mgr.kindParsers[reflect.Uint8] = ParseUint
	mgr.kindParsers[reflect.Uint] = ParseUint
	mgr.kindParsers[reflect.Float64] = ParseFloat
	mgr.kindParsers[reflect.Float32] = ParseFloat
	mgr.kindParsers[reflect.Complex64] = ParseComplex
	mgr.kindParsers[reflect.Complex128] = ParseComplex
	mgr.kindParsers[reflect.Bool] = ParseBool
	return &mgr
}

func (mgr *parserManager) SetParser(t reflect.Type, f ParseFunc) {
	mgr.typeParsers[t] = f
}

func (mgr *parserManager) Parser(t reflect.Type) (ParseFunc, bool) {
	parse, found := mgr.typeParsers[t]
	if found {
		return parse, true
	}
	parse, found = mgr.kindParsers[t.Kind()]
	return parse, found
}
