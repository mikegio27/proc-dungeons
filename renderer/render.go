package renderer

import (
	"fmt"

	"github.com/mikegio27/proc-dungeons/model"
)

func DrawDungeon(d model.Dungeon) {
	DrawWalls(d)
	plane := d.Grid
	starts := make(map[model.Cell]bool, len(d.Starts))
	for _, s := range d.Starts {
		starts[s] = true
	}

	for y := plane.MaxY; y >= plane.MinY; y-- {
		for x := plane.MinX; x <= plane.MaxX; x++ {
			pt := model.Cell{X: x, Y: y}

			ch := d.At(pt).Rune()
			if starts[pt] {
				ch = '*'
			}
			fmt.Printf("%c ", ch)
		}
		fmt.Println()
	}

}

func DrawWalls(d model.Dungeon) {
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
				case model.TileRoomFloor, model.TileCorridor, model.TileDoor:
					d.Set(c, model.TileWall)
					goto next
				}
			}
		next:
		}
	}
}
