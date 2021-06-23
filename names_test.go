package command

import (
	"testing"
)

func assertQuote(t *testing.T, expected, s string) {
	q := quote(s)
	if expected != q {
		t.Errorf("quote failed s=%s Quote(s)=%s", s, q)
	}
}

func verifyFlagName(t *testing.T, original, expected string) {
	name := createFlagName(original)
	if expected != name {
		t.Errorf("%s -> %s, expected %s", original, name, expected)
	}
}

func TestFlagName(t *testing.T) {
	verifyFlagName(t, "DryRun", "dry-run")
}

func TestQuote(t *testing.T) {
	assertQuote(t, `""`, "")
	assertQuote(t, `"a"`, "a")
	assertQuote(t, "\"\\\"\"", "\"")
	assertQuote(t, `"style=\"{{style%d}}\""`, `style="{{style%d}}"`)
}
