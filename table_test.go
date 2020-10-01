package command

import (
	"testing"

	"strings"
)

func FixedTable(t *testing.T) {
	rows := [][]string{
		[]string{"a", "b", "c"},
		[]string{"one", "", "three"},
	}
	var buf strings.Builder
	(&FixedColumnsTable{Prefix: "  ", Separator: " "}).Fprint(&buf, rows)
}
