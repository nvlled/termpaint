package main

import (
	term "github.com/nsf/termbox-go"
	"github.com/nvlled/wind"
	"github.com/nvlled/wind/size"
)

// Note the similarities to brushPallete

const (
	SLOT        = '▩'
	SLOT_SELECT = '█'
)

type colorPallete struct {
	showCursor bool
	cursor     int
	colors     []uint16
}

func NewColorPallete(colors ...uint16) *colorPallete {
	return &colorPallete{
		showCursor: false,
		cursor:     0,
		colors:     colors,
	}
}

func (cp *colorPallete) Width() size.T {
	return size.Const(len(cp.colors) * 2)
}

func (cp *colorPallete) Height() size.T {
	return size.Const(1)
}

func (cp *colorPallete) Color(index int) uint16 {
	if index >= 0 && index < len(cp.colors) {
		return cp.colors[index]
	}
	return ' '
}

func (cp *colorPallete) Render(canvas wind.Canvas) {
	for x, color := range cp.colors {
		canvas.Draw(x*2, 0, SLOT, color, uint16(term.ColorDefault))
	}
	if cp.showCursor {
		color := cp.colors[cp.cursor]
		canvas.Draw(cp.cursor*2, 0, SLOT_SELECT,
			color, color)
	}
}

func (cp *colorPallete) ChooseBrush(events chan term.Event) (int, uint16) {
	end := len(cp.colors)
	cp.showCursor = true
	origCursor := cp.cursor

	redraw()
	for e := range events {
		switch e.Key {
		case term.KeyEsc:
			cp.cursor = origCursor
			goto done
		case term.KeyArrowLeft:
			cp.cursor--
			if cp.cursor < 0 {
				cp.cursor = 0
			}
		case term.KeyArrowRight:
			cp.cursor++
			if cp.cursor >= end {
				cp.cursor = end - 1
			}
		case term.KeyEnter:
			goto done
		}
		redraw()
	}
	/**/ done:

	index := cp.cursor
	cp.showCursor = false
	return index, cp.Color(index)
}
