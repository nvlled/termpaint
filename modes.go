package main

import (
	term "github.com/nsf/termbox-go"
)

type Mode interface {
	Handle(*termPaint, chan term.Event)
}

type normalMode struct {
	brush *term.Cell
	auto  bool
}

func NewNormalMode() *normalMode {
	color := term.ColorDefault
	return &normalMode{
		brush: &term.Cell{' ', color, color},
		auto:  false,
	}
}

func (mode *normalMode) SetBrush(ch rune, color uint16) {
	mode.brush = &term.Cell{ch, term.Attribute(color), 0}
}

func (mode *normalMode) Handle(tpaint *termPaint, events chan term.Event) {
	dArea := tpaint.dArea
	cursor := dArea.cursor

	mode.updateStatus(tpaint.sb)
	redraw()

	e := <-events
	if e.Key != 0 {
		switch e.Key {
		case term.KeyCtrlC:
			saveDrawingArea("dArea", dArea)
			//notify on status bar

		case term.KeyArrowLeft:
			mode.moveCursor(tpaint, -1, 0)
		case term.KeyArrowRight:
			mode.moveCursor(tpaint, 1, 0)
		case term.KeyArrowDown:
			mode.moveCursor(tpaint, 0, 1)
		case term.KeyArrowUp:
			mode.moveCursor(tpaint, 0, -1)

		case term.KeySpace:
			if mode.auto {
				tpaint.dArea.DrawCell(cursor.x, cursor.y, *mode.brush)
			}
			mode.auto = !mode.auto
		}

		w, h := dArea.ActualSize()
		cursor.clamp(0, 0, w, h)

	} else {
		switch e.Ch {
		// insert char
		case 'r':
			tpaint.dArea.DrawCell(cursor.x, cursor.y, *mode.brush)

		// select brush
		case 's':
			pallete := tpaint.bp
			//name := tpaint.sb.Input("enter name", events)
			//tpaint.sb.SetText("You name is dumb " + name)

			_, ch := pallete.ChooseBrush(events)
			mode.brush.Ch = ch
			mode.updateStatus(tpaint.sb)
		case 'c':
			pallete := tpaint.cp
			tpaint.sb.SetText("select color then <enter>")
			_, color := pallete.ChooseBrush(events)
			mode.brush.Fg = term.Attribute(color)
			mode.updateStatus(tpaint.sb)
		}
	}
	redraw()
}

func (mode *normalMode) moveCursor(tpaint *termPaint, dx, dy int) {
	cursor := tpaint.dArea.cursor
	if mode.auto {
		tpaint.dArea.DrawCell(cursor.x, cursor.y, *mode.brush)
	}
	cursor.x += dx
	cursor.y += dy
}

func (mode *normalMode) updateStatus(sb *statusBar) {
	var cells []term.Cell
	text := "normal "
	defColor := term.ColorDefault
	for _, ch := range text {
		cells = append(cells, term.Cell{ch, defColor, defColor})
	}
	cells = append(cells, *mode.brush, term.Cell{' ', 0, 0})

	if mode.auto {
		text = "auto"
		cells = append(cells, addAttribute(text, term.ColorBlack, term.ColorRed)...)
	} else {
		text = "    "
		for _, ch := range text {
			cells = append(cells, term.Cell{ch, defColor, defColor})
		}
	}

	sb.contents = cells
}
