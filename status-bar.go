package main

import (
	term "github.com/nsf/termbox-go"
	"github.com/nvlled/wind"
	"github.com/nvlled/wind/size"
)

type statusBar struct {
	savedContents []term.Cell
	contents      []term.Cell
	cursor        int
}

func NewStatusBar() *statusBar { return &statusBar{cursor: -1} }

func (sb *statusBar) Width() size.T          { return size.Free }
func (sb *statusBar) Height() size.T         { return size.Const(1) }
func (sb *statusBar) Elements() []wind.Layer { return nil }

func (sb *statusBar) Render(canvas wind.Canvas) {
	lastx := 0
	for x, cell := range sb.contents {
		canvas.Draw(x, 0, cell.Ch, uint16(cell.Fg), uint16(cell.Bg))
		lastx = x
	}
	for x := lastx + 1; x < canvas.Width(); x++ {
		canvas.Draw(x, 0, ' ', 0, 0)
	}
	if sb.cursor >= 0 {
		canvas.Draw(sb.cursor, 0, ' ', 0, uint16(term.ColorRed))
	}
}

func (sb *statusBar) SetText(text string) {
	var contents []term.Cell
	for _, ch := range text {
		contents = append(contents, term.Cell{ch, 0, 0})
	}
	sb.contents = contents
}

func (sb *statusBar) Input(prompt string, events chan term.Event) string {
	sb.savedContents = sb.contents
	prompt += ": "
	sb.SetText(prompt)
	sb.cursor = len(prompt)
	start := sb.cursor

	redraw()
	//TODO: redraw(sb)
	for e := range events {
		if e.Key == 0 {
			sb.contents = append(sb.contents, term.Cell{e.Ch, 0, 0})
			sb.cursor++
		} else {
			switch e.Key {
			case term.KeyEnter: // done
				goto done

			case term.KeyBackspace:
				fallthrough
			case term.KeyDelete:
				if sb.cursor-start > 0 {
					sb.contents = sb.contents[:len(sb.contents)-1]
					sb.cursor--
				}
			}
		}
		redraw()
	}
done:
	input := ""
	for i := start; i < len(sb.contents); i++ {
		input += string(sb.contents[i].Ch)
	}

	sb.contents = sb.savedContents
	sb.cursor = -1
	return input
}
