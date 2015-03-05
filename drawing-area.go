package main

import (
	"encoding/gob"
	term "github.com/nsf/termbox-go"
	"github.com/nvlled/wind"
	"github.com/nvlled/wind/size"
	"os"
)

type invalidFileError string

func (err invalidFileError) Error() string {
	return string(err) + ": invalid file format"
}

type drawingArea struct {
	Buffer  [][]term.Cell
	cursor  *pos
	flush   bool
	baseX   int
	baseY   int
	actualW int
	actualH int
	Shit    int
	//redrawList  []pos
}

func init() {
	var buf [][]term.Cell
	gob.Register(buf)
	gob.Register(&drawingArea{})
}

func (dArea *drawingArea) Width() size.T {
	return size.Const(len(dArea.Buffer[0]))
}

func (dArea *drawingArea) Height() size.T {
	return size.Const(len(dArea.Buffer))
}

func (dArea *drawingArea) Elements() []wind.Layer { return nil }

func (dArea *drawingArea) Render(canvas wind.Canvas) {
	// possibly need synchronization from other threads
	dArea.baseX, dArea.baseY = canvas.Abs(0, 0)
	dArea.actualW, dArea.actualH = canvas.Dimension()
	// I could avoid the need to refer to canvas' properties by
	//  * avoiding the use of termbox.SetCursor
	//  * defer drawing operations

	// redraw.add(x, y)

	curs := dArea.cursor
	if curs.x >= 0 && curs.y >= 0 {
		x, y := canvas.Abs(curs.x, curs.y)
		term.SetCursor(x, y)
	}
	if dArea.flush {
		for y, row := range dArea.Buffer {
			for x, cell := range row {
				canvas.Draw(x, y,
					cell.Ch, uint16(cell.Fg), uint16(cell.Bg))
			}
		}
		dArea.flush = false
	} /*TODO: else {
		for _, pos := range dArea.redraw {
			//draw as done above
		}
	}*/
}

func (dArea *drawingArea) Draw(x, y int, ch rune, fg, bg term.Attribute) {
	dArea.Buffer[y][x] = term.Cell{ch, fg, bg}
	// redraw(x, y)
	term.SetCell(dArea.baseX+x, dArea.baseY+y, ch, fg, bg)
}

func (dArea *drawingArea) DrawCell(x, y int, cell term.Cell) {
	dArea.Draw(x, y, cell.Ch, cell.Fg, cell.Bg)
}

func (dArea *drawingArea) Flush() {
	dArea.flush = true
}

func (dArea *drawingArea) ActualSize() (int, int) {
	return dArea.actualW, dArea.actualH
}

func saveDrawingArea(filename string, dArea *drawingArea) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(file)
	return enc.Encode(dArea)
}

func loadDrawingArea(filename string) (*drawingArea, error) {
	var dArea drawingArea

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	dec := gob.NewDecoder(file)
	err = dec.Decode(&dArea)

	if err != nil {
		return nil, invalidFileError(filename)
	}
	dArea.cursor = newpos(0, 0)
	return &dArea, nil
}

func NewDrawingArea(width, height int) *drawingArea {
	// width and height must be greater than zero
	buffer := createBuffer(width, height)
	return &drawingArea{
		Buffer: buffer,
		cursor: newpos(0, 0),
		flush:  false,
	}
}

func createBuffer(width, height int) [][]term.Cell {
	buffer := make([][]term.Cell, height)
	for y := range buffer {
		buffer[y] = make([]term.Cell, width)
	}
	return buffer
}
