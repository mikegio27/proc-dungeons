package generator

import (
	"github.com/mikegio27/proc-dungeons/model"
)

// edgeStartingCell returns a random starting cell on the outer edge of the
// grid that is not inside any room.
func (g *Generator) edgeStartingCell(roomCells map[model.Cell]bool) model.Cell {
	plane := g.cfg.Grid
	width := plane.MaxX - plane.MinX + 1
	height := plane.MaxY - plane.MinY + 1
	perimeter := 2*(width+height) - 4
	if perimeter <= 0 {
		return model.Cell{X: plane.MinX, Y: plane.MinY}
	}

	for {
		// Pick a random position along the perimeter.
		pos := g.rng.Int31n(perimeter)
		var x, y int32

		switch {
		case pos < width:
			// bottom edge, left to right
			x = plane.MinX + pos
			y = plane.MinY
		case pos < width+height-1:
			// right edge, bottom to top (excluding bottom corner)
			x = plane.MaxX
			y = plane.MinY + (pos - width + 1)
		case pos < 2*width+height-2:
			// top edge, right to left (excluding right corner)
			x = plane.MaxX - (pos - (width + height - 1))
			y = plane.MaxY
		default:
			// left edge, top to bottom (excluding top and bottom corners)
			x = plane.MinX
			y = plane.MaxY - (pos - (2*width + height - 2) + 1)
		}

		c := model.Cell{X: x, Y: y}
		if !roomCells[c] {
			return c
		}
	}
}

// randomVisitedCell picks a random cell from the visited set that is not
// inside any room. It returns false if no such cell exists.
func (g *Generator) randomCorridorCell(corridors map[model.Cell]bool, roomCells map[model.Cell]bool) (model.Cell, bool) {
	var candidates []model.Cell
	for c := range corridors {
		if !roomCells[c] {
			candidates = append(candidates, c)
		}
	}
	if len(candidates) == 0 {
		return model.Cell{}, false
	}
	return candidates[g.rng.Intn(len(candidates))], true
}

// findPath uses a BFS search to find a shortest path from start to target,
// avoiding room walls and room interiors (except for the final target
// cell). It returns the sequence of cells to step through, excluding the
// starting cell.
func (g *Generator) findPath(start, target model.Cell, roomCells, roomEdges map[model.Cell]bool) ([]model.Cell, bool) {
	dirs := []model.Cell{{X: 1, Y: 0}, {X: -1, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: -1}}

	queue := []model.Cell{start}
	prev := make(map[model.Cell]model.Cell)
	seen := make(map[model.Cell]bool)
	seen[start] = true

	for len(queue) > 0 {
		c := queue[0]
		queue = queue[1:]

		for _, d := range dirs {
			nx := c.X + d.X
			ny := c.Y + d.Y

			if !g.cfg.Grid.InBounds(model.Cell{X: nx, Y: ny}) {
				continue
			}
			nc := model.Cell{X: nx, Y: ny}
			if seen[nc] {
				continue
			}

			// Always allow reaching the target, even if it's inside the room.
			if nc == target {
				prev[nc] = c
				// Reconstruct path from target back to start.
				var path []model.Cell
				cur := nc
				for cur != start {
					path = append(path, cur)
					cur = prev[cur]
				}
				// Reverse to get start->target order (excluding start).
				for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
					path[i], path[j] = path[j], path[i]
				}
				return path, true
			}

			// Do not walk on room walls or through room interiors.
			if roomEdges[nc] || roomCells[nc] {
				continue
			}

			seen[nc] = true
			prev[nc] = c
			queue = append(queue, nc)
		}
	}

	return nil, false
}

// GenPaths generates a set of paths such that each room is connected to a
// single network of corridors that starts on the grid edge, eventually
// enters each room interior, and never runs along room edge walls.
// It returns both the visited cells (corridors) and the starting cells
// (typically a single edge start) for display.
func (g *Generator) GenPaths(d *model.Dungeon, rooms []model.Room) []model.Cell {
	// Build maps of all room cells and room edge cells so paths can avoid
	// running along walls and through rooms. For each room we also choose a
	// "door" cell on its edge where a corridor is allowed to connect.
	roomCells := make(map[model.Cell]bool)
	roomEdges := make(map[model.Cell]bool)
	roomDoors := make([]model.Cell, len(rooms))
	roomHasDoor := make([]bool, len(rooms))

	for i, room := range rooms {
		local := make(map[model.Cell]bool)
		g.ForEachRoomCell(room, func(c model.Cell) { local[c] = true })

		// Determine edge cells for this room.
		var edgeCells []model.Cell
		for c := range local {
			neighbors := []model.Cell{
				{X: c.X + 1, Y: c.Y},
				{X: c.X - 1, Y: c.Y},
				{X: c.X, Y: c.Y + 1},
				{X: c.X, Y: c.Y - 1},
			}
			for _, n := range neighbors {
				if !local[n] {
					edgeCells = append(edgeCells, c)
					break
				}
			}
		}

		// Choose a random edge cell as a door, if any exist.
		var door model.Cell
		if len(edgeCells) > 0 {
			door = edgeCells[g.rng.Intn(len(edgeCells))]
			roomDoors[i] = door
			roomHasDoor[i] = true
		}

		// Copy into global room cell and edge maps, skipping the door so that
		// corridors can pass through that cell.
		for c := range local {
			if roomHasDoor[i] && c == door {
				continue
			}
			roomCells[c] = true
		}
		for _, c := range edgeCells {
			if roomHasDoor[i] && c == door {
				continue
			}
			roomEdges[c] = true
		}
	}

	corridors := make(map[model.Cell]bool)
	var starts []model.Cell

	for i := range rooms {
		// Use the preselected door as the connection point for this room.
		if !roomHasDoor[i] {
			continue
		}
		target := roomDoors[i]

		var start model.Cell
		if len(corridors) == 0 {
			// first connection: start from edge
			start = g.edgeStartingCell(roomCells)
			starts = append(starts, start)
		} else {
			// Subsequent rooms: start from an existing corridor cell to
			// ensure all rooms are connected into one network.
			if c, ok := g.randomCorridorCell(corridors, roomCells); ok {
				start = c
			} else {
				start = g.edgeStartingCell(roomCells)
				starts = append(starts, start)
			}
		}
		if d.At(start) != model.TileDoor {
			d.Set(start, model.TileCorridor)
		}
		corridors[start] = true

		path, ok := g.findPath(start, target, roomCells, roomEdges)
		if !ok {
			continue
		}
		for _, c := range path {
			// preserve doors
			if d.At(c) == model.TileDoor {
				continue
			}
			d.Set(c, model.TileCorridor)
			corridors[c] = true
		}
	}

	return starts
}
