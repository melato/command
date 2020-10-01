package main

import (
	"melato.org/command"
)

func main() {
	rows := [][]string{
		[]string{"a", "b", "c"},
		[]string{"one", "", "three"},
	}
	command.PrintTableSpaces(rows, " ", " | ")
}
