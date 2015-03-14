package main

import (
	"fmt"
	term "github.com/nsf/termbox-go"
	"github.com/nvlled/wind"
	"os"
	"time"
)

// TODO:
//  - save/load
//  - undo/redo
//  - other modes
//  - fork/modify termbox to support flushing internal buffers to a file
//  - log errors to a file
//  - edit palletes
//  - customizable keybindings
//  - handle terminal resize
//  - drawing area size must be different from terminal size
//  - insert string mode

var modified = true

//var mode Mode = normalMode
//func switchMode(m Mode) { mode = m }
func redraw() { modified = true }

var popupLayer wind.Layer

func getPopupLayer() wind.Layer { return popupLayer }

func setPopupLayer(layer wind.Layer) {
	popupLayer = wind.Border('-', '|', layer)
	redraw()
}

func hidePopupLayer() {
	popupLayer = nil
	redraw()
}

type termPaint struct {
	filename       string
	bp             *brushPallete
	cp             *colorPallete
	dArea          *drawingArea
	modeSb         *statusBar
	secondSb       *statusBar
	sessionBrowser *selector
	editor         *brushPalleteEditor
}

func createPaintLayer(tpaint *termPaint) wind.Layer {
	//hr := wind.Line('―')
	return wind.Vlayer(
		wind.Hlayer(
			wind.Defer(func() wind.Layer { return tpaint.bp }),
			wind.Text("| "),
			tpaint.modeSb,
		),
		//hr,
		wind.Hlayer(tpaint.cp, wind.Text("| "), tpaint.secondSb),
		//hr,
		wind.Zlayer(
			wind.Border('-', '.', tpaint.dArea),
			// RenderLayer returns Free size, not the popupLayer size
			wind.SyncSize(wind.Defer(getPopupLayer), wind.RenderLayer(
				func(canvas wind.Canvas) {
					if popupLayer != nil {
						canvas.Clear()
						popupLayer.Render(canvas)
					}
				},
			)),
		),
	)
}

func main() {
	term.Init()

	brushPalletes := [][]rune{
		{'¶', '»', 'º', '±', 'ß', '÷', 'Ħ'},
		{'ł', 'Œ', 'Ŧ', 'Ʒ', 'ǥ', 'Γ', 'Σ'},
		{'.', 'Ж', 'љ', 'ק', 'گ', '‰', '※'},
	}

	brushPallete := NewBrushPallete(brushPalletes[0]...)
	colorPallete := NewColorPallete(
		uint16(term.ColorDefault),
		uint16(term.ColorCyan),
		uint16(term.ColorRed),
		uint16(term.ColorBlack),
		uint16(term.ColorGreen),
		uint16(term.ColorYellow),
		uint16(term.ColorBlue),
		uint16(term.ColorMagenta),
		uint16(term.ColorWhite),
	)

	var dArea *drawingArea
	var filename string
	if len(os.Args) > 1 {
		var err error
		filename = os.Args[1]
		dArea, err = loadDrawingArea(filename)
		if err != nil {
			term.Close()
			fmt.Printf("file not loaded, %v\n", err)
			os.Exit(1)
		}
		dArea.Flush()
	} else {
		dArea = NewDrawingArea(70, 20)
	}

	tpaint := &termPaint{
		filename:       filename,
		bp:             brushPallete,
		cp:             colorPallete,
		dArea:          dArea,
		modeSb:         NewStatusBar(),
		secondSb:       NewStatusBar(),
		sessionBrowser: NewSelector("Recent sessions", "/home/test/sample", "/tmp/testing", "/var/aaaa", "/home/user/file1"),
		editor:         NewBrushPalleteEditor(brushPalletes),
	}

	paintLayer := createPaintLayer(tpaint)
	canvas := wind.NewTermCanvas()
	events := make(chan term.Event, 1)

	go func() {
		for range time.Tick(33 * time.Millisecond) {
			if modified {
				paintLayer.Render(canvas)
				term.Flush()
				modified = false
			}
		}
	}()

	go func() {
		for {
			e := term.PollEvent()
			events <- e
			if e.Key == term.KeyCtrlC {
				term.Close()
				os.Exit(1)
			}
		}
	}()

	mode := NewNormalMode()
	mode.SetBrush(brushPallete.Brush(0), colorPallete.Color(0))
	for {
		mode.Handle(tpaint, events)
	}
}
