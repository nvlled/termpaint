package main

type pos struct{ x, y int }

func newpos(x, y int) *pos { return &pos{x, y} }

func (p *pos) clamp(minX, minY, maxX, maxY int) {
	if p.x < minX {
		p.x = minX
	}
	if p.x >= maxX {
		p.x = maxX - 1
	}

	if p.y < minY {
		p.y = minY
	}
	if p.y >= maxY {
		p.y = maxY - 1
	}
}
