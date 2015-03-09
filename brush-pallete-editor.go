package main

import (
	term "github.com/nsf/termbox-go"
	"github.com/nvlled/wind"
	"github.com/nvlled/wind/size"
)

const (
	PROMPT = "enter char"
)

type brushPalleteEditor struct {
	brushPalletes []*brushPallete
	vlayer        wind.Layer
	prompt        *statusBar
	cursor        int
}

func (bp *brushPalleteEditor) PalleteAtCursor() *brushPallete {
	i := bp.cursor
	if i >= 0 && i < len(bp.brushPalletes) {
		return bp.brushPalletes[i]
	}
	return nil
}

func NewBrushPalleteEditor(groups [][]rune) *brushPalleteEditor {
	editor := &brushPalleteEditor{}

	var palletes []*brushPallete
	var layers []wind.Layer
	for _, brushes := range groups {
		p := NewBrushPallete(brushes...)
		tapped := wind.TapRender(p, func(pallete wind.Layer, canvas wind.Canvas) {
			// note equality might not work as expected
			if pallete == editor.PalleteAtCursor() {
				canvas = wind.ChangeDefaultColor(0, uint16(term.ColorRed), canvas)
			}
			pallete.Render(canvas)
		})
		palletes = append(palletes, p)
		layers = append(layers, tapped)
	}

	prompt := NewStatusBar()
	if len(palletes) > 0 {
		layers = append(layers, prompt)
	}

	editor.prompt = prompt
	editor.brushPalletes = palletes
	editor.vlayer = wind.Vlayer(layers...)
	return editor
}

func (bp *brushPalleteEditor) Width() size.T {
	return bp.vlayer.Width()
}

func (bp *brushPalleteEditor) Height() size.T {
	return bp.vlayer.Height()
}

func (bp *brushPalleteEditor) Render(canvas wind.Canvas) {
	bp.vlayer.Render(canvas)
}

func (bp *brushPalleteEditor) Edit(events chan term.Event) *brushPallete {
	var pallete *brushPallete
	for e := range events {
		switch e.Key {
		case term.KeyEsc:
			goto end
		case term.KeyEnter:
			pallete = bp.PalleteAtCursor()
			goto end
		case term.KeyCtrlE:
			pallete := bp.PalleteAtCursor()
			if pallete != nil {
				curs := bp.cursor
				bp.cursor = -1

				events_ := make(chan term.Event, 1)
				editing := true

				go func() {
					for e := range events {
						switch e.Key {
						case term.KeyEsc:
							editing = false
							close(events_)
							return
						default:
							events_ <- e
						}
					}
				}()
				for editing {
					i, _ := pallete.ChooseBrush(events_)
					s := bp.prompt.Input(events_, PROMPT)
					if len(s) > 0 {
						pallete.SetBrush(i, rune(s[0]))
					}
				}

				bp.cursor = curs
			}

		case term.KeyArrowDown:
			bp.cursor++
			n := len(bp.brushPalletes) - 1
			if bp.cursor > n {
				bp.cursor = n
			}

		case term.KeyArrowUp:
			bp.cursor--
			if bp.cursor < 0 {
				bp.cursor = 0
			}

		}
		redraw()
	}
end:
	return pallete
}
