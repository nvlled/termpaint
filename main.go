package main

import (
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

var modified = true

//var mode Mode = normalMode
//func switchMode(m Mode) { mode = m }
func redraw() { modified = true }

type termPaint struct {
	bp       *brushPallete
	cp       *colorPallete
	dArea    *drawingArea
	modeSb   *statusBar
	secondSb *statusBar
}

func createPaintLayer(tpaint *termPaint) wind.Layer {
	hr := wind.Line('―')
	return wind.Vlayer(
		wind.Hlayer(tpaint.bp, wind.Text("| "), tpaint.modeSb),
		//hr,
		wind.Hlayer(tpaint.cp, wind.Text("| "), tpaint.secondSb),
		hr,
		tpaint.dArea,
		//wind.Zlayer(tpaint.dArea),
	)
}

func main() {
	term.Init()
	termw, termh := term.Size()

	brushPallete := NewBrushPallete(' ', 'け', '▪', '*', '~', '·', '▲', '★', '〄', '˚', 'Ξ', 'ϡ', 'њ', 'ر')
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
	dArea, err := loadDrawingArea("dArea")
	if err != nil {
		dArea = NewDrawingArea(termw, termh)
	}
	dArea.Flush()

	tpaint := &termPaint{
		bp:       brushPallete,
		cp:       colorPallete,
		dArea:    dArea,
		modeSb:   NewStatusBar(),
		secondSb: NewStatusBar(),
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
