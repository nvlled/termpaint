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

	mode.updateStatus(tpaint.modeSb)
	redraw()

	e := <-events
	if e.Key != 0 {
		switch e.Key {
		case term.KeyCtrlS:
			filename := tpaint.filename
			if filename == "" {
				filename = tpaint.secondSb.Input("saving session, enter filename", events)
			}

			if filename == "" {
				tpaint.secondSb.SetText("no filename given, aborting...")
			} else {
				saveDrawingArea(filename, dArea)
				tpaint.secondSb.SetText("session saved on " + filename)
			}

		case term.KeyCtrlC:
			if tpaint.filename != "" {
				saveDrawingArea(tpaint.filename, dArea)
			}

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
			tpaint.secondSb.SetText("select brush then <enter>")
			_, ch := pallete.ChooseBrush(events)
			tpaint.secondSb.SetText("")
			mode.brush.Ch = ch
			mode.updateStatus(tpaint.modeSb)
		case 'c':
			pallete := tpaint.cp
			tpaint.secondSb.SetText("select color then <enter>")
			_, color := pallete.ChooseBrush(events)
			tpaint.secondSb.SetText("")
			mode.brush.Fg = term.Attribute(color)
			mode.updateStatus(tpaint.modeSb)
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
