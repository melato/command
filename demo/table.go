package main

import (
	"melato.org/export/command"
)

func main() {
	rows := [][]string{
		[]string{"a", "b", "c"},
		[]string{"one", "", "three"},
	}
	command.PrintTableSpaces(rows, " ", " | ")
}
