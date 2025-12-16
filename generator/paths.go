package generator

import "github.com/mikegio27/proc-dungeons/model"

// GenPaths connects every room to a single corridor network.
// - One door per room (edge cell)
// - First corridor starts at perimeter
// - Subsequent rooms connect from existing corridor cell
// - Corridors keep distance from rooms via CorridorBuff (except at doors)
// - CorridorW controls thickness
func (g *Generator) GenPaths(d *model.Dungeon, rooms []model.Room) []model.Cell {
	// ---- 1) Room footprints + doors ----

	roomCells := make(map[model.Cell]bool) // interior (excluding door)
	roomEdges := make(map[model.Cell]bool) // edge ring (excluding door)

	roomDoors := make([]model.Cell, len(rooms))
	roomHasDoor := make([]bool, len(rooms))

	for i, room := range rooms {
		local := make(map[model.Cell]bool)
		g.ForEachRoomCell(room, func(c model.Cell) { local[c] = true })

		// edge cells: any cell with a neighbor not in local
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

		// choose a door
		if len(edgeCells) > 0 {
			// ensure door is not on the edge of the dungeon grid
			var validEdgeCells []model.Cell
			for _, c := range edgeCells {
				if c.X > g.cfg.Grid.MinX && c.X < g.cfg.Grid.MaxX &&
					c.Y > g.cfg.Grid.MinY && c.Y < g.cfg.Grid.MaxY {
					validEdgeCells = append(validEdgeCells, c)
				}
			}
			if len(validEdgeCells) == 0 {
				// fallback to any edge cell
				validEdgeCells = edgeCells
			}
			door := validEdgeCells[g.rng.Intn(len(validEdgeCells))]
			roomDoors[i] = door
			roomHasDoor[i] = true
			d.Set(door, model.TileDoor)

			// add to global maps, skipping door
			for c := range local {
				if c == door {
					continue
				}
				roomCells[c] = true
			}
			for _, c := range edgeCells {
				if c == door {
					continue
				}
				roomEdges[c] = true
			}
		} else {
			// degenerate: treat whole thing as roomCells
			for c := range local {
				roomCells[c] = true
			}
		}
	}

	roomSolid := mergeBoolMaps(roomCells, roomEdges)

	// ---- 2) Build ONE blocked map (rooms + edge ring + buffer) ----

	// Base blocked: everything within CorridorBuff of rooms/edges.
	blockedBase := g.expand(roomSolid, g.cfg.CorridorBuff)
	for c := range roomSolid {
		blockedBase[c] = true
	}

	// Allow doors + a *small* approach area so BFS can actually attach.
	// If you clear the full buff radius, you basically undo the whole idea.
	for i := range rooms {
		if !roomHasDoor[i] {
			continue
		}
		door := roomDoors[i]
		// Door cell must be allowed
		delete(blockedBase, door)

		// Also allow a 1-tile halo outside the door so corridors can “plug in”
		g.clearRadius(blockedBase, door, 1)
	}

	// ---- 3) Connect rooms into a single corridor network ----

	corridors := make(map[model.Cell]bool)
	var starts []model.Cell

	for i := range rooms {
		if !roomHasDoor[i] {
			continue
		}
		target := roomDoors[i]

		var start model.Cell
		if len(corridors) == 0 {
			start = g.edgeStartingCell(roomSolid)
			starts = append(starts, start)
		} else {
			if c, ok := g.randomCorridorCell(corridors); ok {
				start = c
			} else {
				start = g.edgeStartingCell(roomSolid)
				starts = append(starts, start)
			}
		}

		// carve start
		if d.At(start) != model.TileDoor {
			g.carveCorridor(d, start, blockedBase)
		}
		corridors[start] = true

		path, ok := g.findPath(start, target, blockedBase)
		if !ok {
			continue
		}

		for _, c := range path {
			if d.At(c) == model.TileDoor {
				continue
			}
			g.carveCorridor(d, c, blockedBase)
			corridors[c] = true
		}
	}

	return starts
}

// edgeStartingCell returns a random cell on the perimeter that is not inside roomSolid.
func (g *Generator) edgeStartingCell(roomSolid map[model.Cell]bool) model.Cell {
	plane := g.cfg.Grid
	width := plane.MaxX - plane.MinX + 1
	height := plane.MaxY - plane.MinY + 1
	perimeter := 2*(width+height) - 4
	if perimeter <= 0 {
		return model.Cell{X: plane.MinX, Y: plane.MinY}
	}

	for {
		pos := g.rng.Int31n(perimeter)
		var x, y int32

		switch {
		case pos < width:
			x = plane.MinX + pos
			y = plane.MinY
		case pos < width+height-1:
			x = plane.MaxX
			y = plane.MinY + (pos - width + 1)
		case pos < 2*width+height-2:
			x = plane.MaxX - (pos - (width + height - 1))
			y = plane.MaxY
		default:
			x = plane.MinX
			y = plane.MaxY - (pos - (2*width + height - 2) + 1)
		}

		c := model.Cell{X: x, Y: y}
		if !roomSolid[c] {
			return c
		}
	}
}

// randomCorridorCell picks a random existing corridor cell.
func (g *Generator) randomCorridorCell(corridors map[model.Cell]bool) (model.Cell, bool) {
	if len(corridors) == 0 {
		return model.Cell{}, false
	}
	cands := make([]model.Cell, 0, len(corridors))
	for c := range corridors {
		cands = append(cands, c)
	}
	return cands[g.rng.Intn(len(cands))], true
}

// carveCorridor writes corridor tiles around center cell according to CorridorW,
// refuses blocked cells, and never overwrites non-empty tiles.
func (g *Generator) carveCorridor(d *model.Dungeon, center model.Cell, blocked map[model.Cell]bool) {
	w := g.cfg.CorridorW
	if w < 1 {
		w = 1
	}
	r := w / 2

	for y := center.Y - r; y <= center.Y+r; y++ {
		for x := center.X - r; x <= center.X+r; x++ {
			c := model.Cell{X: x, Y: y}
			if !d.InBounds(c) {
				continue
			}
			if blocked[c] {
				continue
			}
			if d.At(c) != model.TileEmpty {
				continue
			}
			d.Set(c, model.TileCorridor)
		}
	}
}

// findPath BFS from start to target, avoiding blocked cells.
// Returns path excluding start (includes target).
func (g *Generator) findPath(start, target model.Cell, blocked map[model.Cell]bool) ([]model.Cell, bool) {
	dirs := []model.Cell{{X: 1, Y: 0}, {X: -1, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: -1}}

	queue := []model.Cell{start}
	prev := make(map[model.Cell]model.Cell)
	seen := make(map[model.Cell]bool)
	seen[start] = true

	for len(queue) > 0 {
		c := queue[0]
		queue = queue[1:]

		for _, d := range dirs {
			nc := model.Cell{X: c.X + d.X, Y: c.Y + d.Y}
			if !g.cfg.Grid.InBounds(nc) {
				continue
			}
			if seen[nc] {
				continue
			}

			// allow reaching target even if it's blocked
			if nc != target && blocked[nc] {
				continue
			}

			seen[nc] = true
			prev[nc] = c

			if nc == target {
				var path []model.Cell
				for cur := nc; cur != start; cur = prev[cur] {
					path = append(path, cur)
				}
				for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
					path[i], path[j] = path[j], path[i]
				}
				return path, true
			}

			queue = append(queue, nc)
		}
	}

	return nil, false
}

// expand includes all cells within Chebyshev radius r of any input cell.
func (g *Generator) expand(cells map[model.Cell]bool, r int32) map[model.Cell]bool {
	out := make(map[model.Cell]bool, len(cells))
	if r <= 0 {
		for c := range cells {
			out[c] = true
		}
		return out
	}
	for c := range cells {
		for dy := -r; dy <= r; dy++ {
			for dx := -r; dx <= r; dx++ {
				n := model.Cell{X: c.X + dx, Y: c.Y + dy}
				if g.cfg.Grid.InBounds(n) {
					out[n] = true
				}
			}
		}
	}
	return out
}

func (g *Generator) clearRadius(m map[model.Cell]bool, center model.Cell, r int32) {
	if r < 0 {
		return
	}
	for dy := -r; dy <= r; dy++ {
		for dx := -r; dx <= r; dx++ {
			delete(m, model.Cell{X: center.X + dx, Y: center.Y + dy})
		}
	}
}

func mergeBoolMaps(a, b map[model.Cell]bool) map[model.Cell]bool {
	out := make(map[model.Cell]bool, len(a)+len(b))
	for k := range a {
		out[k] = true
	}
	for k := range b {
		out[k] = true
	}
	return out
}
