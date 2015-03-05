package main

import (
	term "github.com/nsf/termbox-go"
	"github.com/nvlled/wind"
	"github.com/nvlled/wind/size"
)

type brushPallete struct {
	showCursor bool
	cursor     int
	brushes    []rune
}

func NewBrushPallete(brushes ...rune) *brushPallete {
	return &brushPallete{
		showCursor: false,
		cursor:     0,
		brushes:    brushes,
	}
}

func (bp *brushPallete) Width() size.T {
	return size.Const(len(bp.brushes) * 2)
}

func (bp *brushPallete) Height() size.T {
	return size.Const(1)
}

func (bp *brushPallete) Brush(index int) rune {
	if index >= 0 && index < len(bp.brushes) {
		return bp.brushes[index]
	}
	return ' '
}

func (_ *brushPallete) Elements() []wind.Layer { return nil }

func (bp *brushPallete) Render(canvas wind.Canvas) {
	for x, brush := range bp.brushes {
		canvas.Draw(x*2, 0, brush, 0, 0)
	}
	if bp.showCursor {
		canvas.Draw(bp.cursor*2, 0, bp.brushes[bp.cursor],
			uint16(term.ColorDefault), uint16(term.ColorRed))
	}
}

func (bp *brushPallete) ChooseBrush(events chan term.Event) (int, rune) {
	end := len(bp.brushes)
	bp.showCursor = true
	origCursor := bp.cursor

	redraw()
	for e := range events {
		switch e.Key {
		case term.KeyEsc:
			bp.cursor = origCursor
			goto done
		case term.KeyArrowLeft:
			bp.cursor--
			if bp.cursor < 0 {
				bp.cursor = 0
			}
		case term.KeyArrowRight:
			bp.cursor++
			if bp.cursor >= end {
				bp.cursor = end - 1
			}
		case term.KeyEnter:
			goto done
		}
		redraw()
	}
	/**/ done:

	index := bp.cursor
	bp.showCursor = false
	return index, bp.Brush(index)
}
