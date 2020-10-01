package reflx

import (
	"reflect"
	"testing"
)

type B struct {
	X string
}

func TestParser(t *testing.T) {
	mgr := NewParserManager()
	var b B
	v := reflect.ValueOf(&b).Elem().Field(0)
	mgr.Parse("", v, "x1")
	if "x1" != b.X {
		t.Errorf("expected: %s actual: %s", "x1", b.X)
	}
}
