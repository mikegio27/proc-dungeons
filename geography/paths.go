package geography

import "fmt"

type Cell struct {
	X int32
	Y int32
}

// edgeStartingCell returns a random starting cell on the outer edge of the
// grid that is not inside any room.
func edgeStartingCell(roomCells map[Cell]bool) Cell {
	plane := GetPlane()
	width := plane.MaxX - plane.MinX + 1
	height := plane.MaxY - plane.MinY + 1
	perimeter := 2*(width+height) - 4
	if perimeter <= 0 {
		return Cell{X: plane.MinX, Y: plane.MinY}
	}

	for {
		// Pick a random position along the perimeter.
		pos := rng.Int31n(perimeter)
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

		c := Cell{X: x, Y: y}
		if !roomCells[c] {
			return c
		}
	}
}

// randomVisitedCell picks a random cell from the visited set that is not
// inside any room. It returns false if no such cell exists.
func randomVisitedCell(visited map[Cell]bool, roomCells map[Cell]bool) (Cell, bool) {
	var candidates []Cell
	for c := range visited {
		if !roomCells[c] {
			candidates = append(candidates, c)
		}
	}
	if len(candidates) == 0 {
		return Cell{}, false
	}
	return candidates[rng.Intn(len(candidates))], true
}

// findPath uses a BFS search to find a shortest path from start to target,
// avoiding room walls and room interiors (except for the final target
// cell). It returns the sequence of cells to step through, excluding the
// starting cell.
func findPath(start, target Cell, roomCells, roomEdges map[Cell]bool) ([]Cell, bool) {
	plane := GetPlane()
	dirs := [][2]int32{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

	queue := []Cell{start}
	prev := make(map[Cell]Cell)
	seen := make(map[Cell]bool)
	seen[start] = true

	for len(queue) > 0 {
		c := queue[0]
		queue = queue[1:]

		for _, d := range dirs {
			nx := c.X + d[0]
			ny := c.Y + d[1]

			if nx < plane.MinX || nx > plane.MaxX || ny < plane.MinY || ny > plane.MaxY {
				continue
			}
			nc := Cell{X: nx, Y: ny}
			if seen[nc] {
				continue
			}

			// Always allow reaching the target, even if it's inside the room.
			if nc == target {
				prev[nc] = c
				// Reconstruct path from target back to start.
				var path []Cell
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

func DrawGrid(visited map[Cell]bool, starts map[Cell]bool) {
	plane := GetPlane()
	for y := plane.MaxY; y >= plane.MinY; y-- {
		for x := plane.MinX; x <= plane.MaxX; x++ {
			ch := '.'
			pt := Cell{X: x, Y: y}
			if starts[pt] {
				ch = '*'
			} else if visited[pt] {
				ch = '#'
			}
			if x == 0 && y == 0 {
				ch = '+'
			}
			fmt.Printf("%c ", ch)
		}
		fmt.Println()
	}
}

// GenPaths generates a set of paths such that each room is connected to a
// single network of corridors that starts on the grid edge, eventually
// enters each room interior, and never runs along room edge walls.
// It returns both the visited cells (corridors) and the starting cells
// (typically a single edge start) for display.
func GenPaths(rooms []Room) (map[Cell]bool, map[Cell]bool) {
	// Build maps of all room cells and room edge cells so paths can avoid
	// running along walls and through rooms. For each room we also choose a
	// "door" cell on its edge where a corridor is allowed to connect.
	roomCells := make(map[Cell]bool)
	roomEdges := make(map[Cell]bool)
	roomDoors := make([]Cell, len(rooms))
	roomHasDoor := make([]bool, len(rooms))

	for i, room := range rooms {
		local := make(map[Cell]bool)

		switch room.Shape {
		case Rectangle, Square:
			fillRectRoom(local, room)
		case Circle:
			fillCircleRoom(local, room)
		case Triangle:
			fillTriangleRoom(local, room)
		default:
			fillRectRoom(local, room)
		}

		// Determine edge cells for this room.
		var edgeCells []Cell
		for c := range local {
			neighbors := []Cell{{c.X + 1, c.Y}, {c.X - 1, c.Y}, {c.X, c.Y + 1}, {c.X, c.Y - 1}}
			for _, n := range neighbors {
				if !local[n] {
					edgeCells = append(edgeCells, c)
					break
				}
			}
		}

		// Choose a random edge cell as a door, if any exist.
		var door Cell
		if len(edgeCells) > 0 {
			door = edgeCells[rng.Intn(len(edgeCells))]
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

	visited := make(map[Cell]bool)
	starts := make(map[Cell]bool)

	for i := range rooms {
		// Use the preselected door as the connection point for this room.
		if !roomHasDoor[i] {
			continue
		}
		target := roomDoors[i]

		var start Cell
		if i == 0 {
			// First room: start from the outer edge.
			start = edgeStartingCell(roomCells)
			starts[start] = true
		} else {
			// Subsequent rooms: start from an existing corridor cell to
			// ensure all rooms are connected into one network.
			if c, ok := randomVisitedCell(visited, roomCells); ok {
				start = c
			} else {
				start = edgeStartingCell(roomCells)
				starts[start] = true
			}
		}
		visited[start] = true

		path, ok := findPath(start, target, roomCells, roomEdges)
		if !ok {
			continue
		}
		for _, c := range path {
			visited[c] = true
		}
	}

	return visited, starts
}
