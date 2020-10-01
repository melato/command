package command

import (
	"fmt"
	"testing"
)

func assertQuote(t *testing.T, expected, s string) {
	q := Quote(s)
	if expected != q {
		fmt.Println("quote failed s=" + s + " Quote(s)=" + q)
		t.Fail()
	}
}

func verifyFlagName(t *testing.T, original, expected string) {
	name := CreateFlagName(original)
	if expected != name {
		t.Errorf("%s -> %s, expected %s", original, name, expected)
	}
}

func TestFlagName(t *testing.T) {
	verifyFlagName(t, "DryRun", "dry-run")
}

func TestQuote(t *testing.T) {
	//assertQuote(t, `""`, "")
	//assertQuote(t, `"a"`, "a")
	assertQuote(t, "\"\\\"\"", "\"")
}
