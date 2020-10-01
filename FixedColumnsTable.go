package command

import (
	"io"
	"os"
	"strings"
)

type FixedColumnsTable struct {
	Prefix    string
	Separator string
}

func repeatString(s string, n int) string {
	var buf strings.Builder
	for i := 0; i < n; i++ {
		buf.WriteString(s)
	}
	return buf.String()
}

func (t *FixedColumnsTable) Fprint(writer io.StringWriter, rows [][]string) {
	cols := 0
	for _, row := range rows {
		n := len(row)
		if n > cols {
			cols = n
		}
	}
	colSizes := make([]int, cols, cols)
	for _, row := range rows {
		for i, s := range row {
			n := len(s)
			if n > colSizes[i] {
				colSizes[i] = n
			}
		}
	}
	maxSize := 0
	for _, n := range colSizes {
		if n > maxSize {
			maxSize = n
		}
	}
	spaces := repeatString(" ", maxSize)
	for _, row := range rows {
		var buf strings.Builder
		buf.WriteString(t.Prefix)
		ns := 0
		for i, s := range row {
			if i > 0 {
				buf.WriteString(spaces[0:ns])
				buf.WriteString(t.Separator)
			}
			buf.WriteString(s)
			ns = colSizes[i] - len(s)
		}
		buf.WriteString("\n")
		writer.WriteString(buf.String())
	}
}

func (t *FixedColumnsTable) Print(rows [][]string) {
	t.Fprint(os.Stdout, rows)
}
