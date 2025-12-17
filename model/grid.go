package model

type Grid struct {
	MaxX int32
	MaxY int32
	MinY int32
	MinX int32
}

type Cell struct {
	X int32
	Y int32
}

func (g Grid) Width() int32 {
	return g.MaxX - g.MinX + 1
}

func (g Grid) Height() int32 {
	return g.MaxY - g.MinY + 1
}

func (g Grid) InBounds(c Cell) bool {
	return c.X >= g.MinX && c.X <= g.MaxX && c.Y >= g.MinY && c.Y <= g.MaxY
}

func (g Grid) RoomInBoundsWithPadding(r Room, pad int32) bool {
	return r.TopLeft.X >= g.MinX+pad &&
		r.BottomRight.X <= g.MaxX-pad &&
		r.TopLeft.Y >= g.MinY+pad &&
		r.BottomRight.Y <= g.MaxY-pad
}

func (g Grid) OnGridBoundary(c Cell) bool {
	return c.X == g.MinX || c.X == g.MaxX || c.Y == g.MinY || c.Y == g.MaxY
}

func (g Grid) OnYBoundary(c Cell) bool {
	return c.Y == g.MinY || c.Y == g.MaxY
}

func (g Grid) OnXBoundary(c Cell) bool {
	return c.X == g.MinX || c.X == g.MaxX
}

func (g Grid) Index(c Cell) (int32, bool) {
	if !g.InBounds(c) {
		return -1, false
	}
	width := g.Width()
	xOffset := c.X - g.MinX
	yOffset := c.Y - g.MinY
	index := yOffset*width + xOffset
	return index, true
}
