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
