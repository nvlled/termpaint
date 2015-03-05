package main

import (
	term "github.com/nsf/termbox-go"
)

func addAttribute(s string, fg, bg term.Attribute) []term.Cell {
	var cells []term.Cell
	for _, ch := range s {
		cells = append(cells, term.Cell{ch, fg, bg})
	}
	return cells
}
