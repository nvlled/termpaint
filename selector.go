package main

import (
	term "github.com/nsf/termbox-go"
	"github.com/nvlled/wind"
	"github.com/nvlled/wind/size"
)

const (
	INDENT_SIZE = 2
	PUDDING     = 2
)

var highlightColor = term.ColorRed

type selector struct {
	index int
	title string
	list  []string
	width int
}

func NewSelector(title string, list ...string) *selector {
	width := len(title)
	for _, s := range list {
		if len(s) > width {
			width = len(s)
		}
	}
	return &selector{
		index: 0,
		title: title,
		list:  list,
		width: width + PUDDING,
	}
}

func (sel *selector) Width() size.T { return size.Const(sel.width) }

func (sel *selector) Height() size.T {
	// +1 including the title
	return size.Const(len(sel.list) + 1)
}

func (sel *selector) Render(canvas wind.Canvas) {
	canvas.Clear()
	for x, ch := range []rune(sel.title) {
		canvas.Draw(x, 0, ch, 0, 0)
	}
	y := 1
	for i, s := range sel.list {
		var bg uint16 = 0
		if i == sel.index {
			bg = uint16(highlightColor)
		}
		canvas.DrawText(INDENT_SIZE, y, s, 0, bg)
		y++
	}
}

func (sel *selector) Select(events chan term.Event) (int, string) {
	for e := range events {
		switch e.Key {
		case term.KeyEnter:
			goto there

		case term.KeyArrowUp:
			sel.index--
			if sel.index < 0 {
				sel.index = 0
			}

		case term.KeyArrowDown:
			sel.index++
			if sel.index >= len(sel.list) {
				sel.index = len(sel.list) - 1
			}

		}
		redraw()
	}
there:
	return sel.index, sel.list[sel.index]
}
