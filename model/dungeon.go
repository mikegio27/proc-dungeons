package model

type Dungeon struct {
	Rooms  []Room
	Tiles  []Tile
	Grid   Grid
	Starts []Cell
}

func NewDungeon(grid Grid) Dungeon {
	size := int(grid.Width() * grid.Height())
	tiles := make([]Tile, size)
	// zero value of Tile is TileEmpty, so no need to fill
	return Dungeon{
		Grid:  grid,
		Tiles: tiles,
	}
}

func (d Dungeon) InBounds(c Cell) bool { return d.Grid.InBounds(c) }

func (d Dungeon) At(c Cell) Tile {
	idx, ok := d.Grid.Index(c)
	if !ok {
		return TileEmpty
	}
	return d.Tiles[int(idx)]
}

func (d *Dungeon) Set(c Cell, t Tile) {
	idx, ok := d.Grid.Index(c)
	if !ok {
		return
	}
	d.Tiles[int(idx)] = t
}
