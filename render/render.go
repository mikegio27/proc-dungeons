package render

import (
	"fmt"

	"github.com/mikegio27/proc-dungeons/model"
)

func DrawDungeon(d *model.Dungeon) {
	DrawWalls(d)

	g := d.Grid

	starts := make(map[model.Cell]bool, len(d.Starts))
	for _, s := range d.Starts {
		starts[s] = true
	}

	for y := g.MaxY + 1; y >= g.MinY-1; y-- {
		for x := g.MinX - 1; x <= g.MaxX+1; x++ {
			c := model.Cell{X: x, Y: y}
			fmt.Printf("%c ", glyphAt(d, c, starts))
		}
		fmt.Println()
	}
}

func glyphAt(d *model.Dungeon, c model.Cell, starts map[model.Cell]bool) rune {
	if starts[c] {
		return '*'
	}

	if !d.InBounds(c) {
		if adjacentToStart(c, starts) {
			return '*'
		}
		return model.TileWall.Rune()
	}

	return d.At(c).Rune()
}

func adjacentToStart(c model.Cell, starts map[model.Cell]bool) bool {
	dirs := []model.Cell{
		{X: 1, Y: 0}, {X: -1, Y: 0},
		{X: 0, Y: 1}, {X: 0, Y: -1},
	}
	for _, di := range dirs {
		if starts[model.Cell{X: c.X + di.X, Y: c.Y + di.Y}] {
			return true
		}
	}
	return false
}

func DrawWalls(d *model.Dungeon) {
	dirs := []model.Cell{
		{X: 1, Y: 0}, {X: -1, Y: 0},
		{X: 0, Y: 1}, {X: 0, Y: -1},
	}

	for y := d.Grid.MinY; y <= d.Grid.MaxY; y++ {
		for x := d.Grid.MinX; x <= d.Grid.MaxX; x++ {
			c := model.Cell{X: x, Y: y}
			if d.At(c) != model.TileEmpty {
				continue
			}

			for _, di := range dirs {
				n := model.Cell{X: c.X + di.X, Y: c.Y + di.Y}
				switch d.At(n) {
				case model.TileRoomFloor, model.TileDoor:
					d.Set(c, model.TileWall)
					goto next
				}
			}
		next:
		}
	}
}
