package generator

import (
	"math"

	"github.com/mikegio27/proc-dungeons/model"
)

// maxRoomAreaFraction controls the maximum fraction of the total grid area
// that any single room's bounding box is allowed to occupy.
const maxRoomAreaFraction = 0.15

// maxTotalRoomAreaFraction controls the maximum fraction of the grid area
// that all rooms combined are allowed to occupy.
const maxTotalRoomAreaFraction = 0.45

// minRoomGap is the minimum number of tiles that should separate the
// bounding boxes of any two rooms. This helps prevent rooms from being
// squished directly against each other.
const minRoomGap = 2

// roomDimensions returns random width and height for the bounding box of a
// given room shape. Dimensions are constrained so that the room stays
// within a reasonable size relative to the overall plane and preserves the
// basic proportions of each shape.
func (g *Generator) roomDimensions(shape model.RoomId) (width, height int32) {
	plane := g.cfg.Grid
	gridWidth := plane.MaxX - plane.MinX + 1
	gridHeight := plane.MaxY - plane.MinY + 1
	if gridWidth < 3 {
		gridWidth = 3
	}
	if gridHeight < 3 {
		gridHeight = 3
	}
	gridArea := gridWidth * gridHeight
	maxArea := int32(maxRoomAreaFraction * float64(gridArea))
	if maxArea < 9 { // ensure at least a small 3x3 room is possible
		maxArea = 9
	}

	minSize := int32(3)

	switch shape {
	case model.Rectangle, model.Triangle:
		// Allow independent width/height up to half the grid each, but
		// constrained by the maxArea.
		maxW := max(gridWidth/2, minSize)
		maxH := max(gridHeight/2, minSize)

		// Try random dimensions that satisfy the area constraint.
		for range 10 {
			w := g.rng.Int31n(maxW-minSize+1) + minSize
			h := g.rng.Int31n(maxH-minSize+1) + minSize
			if w*h <= maxArea {
				return w, h
			}
		}

		// Fallback: derive dimensions directly from the max area.
		w := min(max(int32(math.Sqrt(float64(maxArea))), minSize), maxW)
		h := min(max(maxArea/w, minSize), maxH)
		return w, h

	case model.Circle, model.Square:
		// Circles and squares use a square bounding box.
		maxSideByPlane := min(gridHeight, gridWidth)
		maxSideByArea := int32(math.Sqrt(float64(maxArea)))
		maxSide := max(min(maxSideByArea, maxSideByPlane), minSize)

		side := g.rng.Int31n(maxSide-minSize+1) + minSize
		return side, side

	default:
		// Reasonable default: small square room.
		return minSize, minSize
	}
}

// roomEdges chooses a random top-left position for a room of the given
// shape so that its bounding box fits entirely within the current grid.
func (g *Generator) roomEdges(shape model.RoomId) (topLeft, bottomRight model.Cell) {
	plane := g.cfg.Grid

	width, height := g.roomDimensions(shape)

	// Ensure the room fits in the grid; if not, clamp to grid size.
	gridWidth := plane.MaxX - plane.MinX + 1
	gridHeight := plane.MaxY - plane.MinY + 1
	if width > gridWidth {
		width = gridWidth
	}
	if height > gridHeight {
		height = gridHeight
	}

	// Compute valid ranges for the top-left corner.
	maxLeftX := plane.MaxX - width + 1
	maxTopY := plane.MaxY - height + 1

	xRange := maxLeftX - plane.MinX + 1
	yRange := maxTopY - plane.MinY + 1

	// Randomly choose a top-left within the valid ranges.
	x := g.rng.Int31n(xRange) + plane.MinX
	y := g.rng.Int31n(yRange) + plane.MinY

	topLeft = model.Cell{X: x, Y: y}
	bottomRight = model.Cell{X: x + width - 1, Y: y + height - 1}
	return
}

func (g *Generator) RandomRoom() model.Room {
	shapes := g.cfg.RoomShapes
	if len(shapes) == 0 {
		panic("no room shapes configured")
	}

	shape := shapes[g.rng.Intn(len(shapes))]
	// top left and bottom right positions must be within the constraints of the roomEdges
	topLeft, bottomRight := g.roomEdges(shape)
	return model.Room{
		Shape:       shape,
		TopLeft:     topLeft,
		BottomRight: bottomRight,
	}
}

// roomArea returns the area of the room's bounding box.
func roomArea(r model.Room) int32 {
	width := r.BottomRight.X - r.TopLeft.X + 1
	height := r.BottomRight.Y - r.TopLeft.Y + 1
	if width <= 0 || height <= 0 {
		return 0
	}
	return width * height
}

// roomsTooClose reports whether two rooms are closer than the specified
// gap, based on their bounding boxes.
func roomsTooClose(a, b model.Room, gap int32) bool {
	ax1, ax2 := a.TopLeft.X, a.BottomRight.X
	ay1, ay2 := a.TopLeft.Y, a.BottomRight.Y
	bx1, bx2 := b.TopLeft.X, b.BottomRight.X
	by1, by2 := b.TopLeft.Y, b.BottomRight.Y

	// Require at least `gap` tiles of separation along both axes.
	if ax1 > bx2+gap || bx1 > ax2+gap || ay1 > by2+gap || by1 > ay2+gap {
		return false
	}
	return true
}

// Rooms generates up to maxRooms rooms, enforcing both a minimum spacing
// between rooms and a cap on the total area that all rooms may occupy.
func (g *Generator) Rooms(maxRooms int) []model.Room {
	plane := g.cfg.Grid
	gridWidth := plane.MaxX - plane.MinX + 1
	gridHeight := plane.MaxY - plane.MinY + 1
	gridArea := gridWidth * gridHeight
	maxTotalArea := int32(maxTotalRoomAreaFraction * float64(gridArea))
	if maxTotalArea <= 0 {
		maxTotalArea = gridArea
	}

	rooms := make([]model.Room, 0, maxRooms)
	var usedArea int32

	for len(rooms) < maxRooms {
		success := false
		// Try several times to place a room that satisfies constraints.
		for range 20 {
			candidate := g.RandomRoom()
			area := roomArea(candidate)
			if area == 0 || usedArea+area > maxTotalArea {
				continue
			}

			tooClose := false
			for _, existing := range rooms {
				if roomsTooClose(existing, candidate, minRoomGap) {
					tooClose = true
					break
				}
			}
			if tooClose {
				continue
			}

			rooms = append(rooms, candidate)
			usedArea += area
			success = true
			break
		}

		if !success {
			// Could not place any more rooms without breaking constraints.
			break
		}
	}

	return rooms
}

// fillRectRoom marks all cells inside the rectangular bounds of the room.
func (g *Generator) eachRect(room model.Room, fn func(model.Cell)) {
	for y := room.TopLeft.Y; y <= room.BottomRight.Y; y++ {
		for x := room.TopLeft.X; x <= room.BottomRight.X; x++ {
			fn(model.Cell{X: x, Y: y})
		}
	}
}

// fillCircleRoom approximates a circle inside the room's bounding box,
// smoothing the corners compared to a plain rectangle.
func (g *Generator) eachCircle(room model.Room, fn func(model.Cell)) {
	// Compute center of the bounding box.
	cx := float64(room.TopLeft.X+room.BottomRight.X) / 2.0
	cy := float64(room.TopLeft.Y+room.BottomRight.Y) / 2.0

	width := float64(room.BottomRight.X-room.TopLeft.X) + 1.0
	height := float64(room.BottomRight.Y-room.TopLeft.Y) + 1.0
	// Use the smaller half-dimension as radius.
	r := math.Min(width, height) / 2.0
	r2 := r * r

	for y := room.TopLeft.Y; y <= room.BottomRight.Y; y++ {
		for x := room.TopLeft.X; x <= room.BottomRight.X; x++ {
			dx := float64(x) - cx
			dy := float64(y) - cy
			if dx*dx+dy*dy <= r2+0.25 {
				fn(model.Cell{X: x, Y: y})
			}
		}
	}
}

// fillTriangleRoom fills an isosceles triangle within the room's
// bounding box. The triangle has its apex at the top and base at the
// bottom of the box.
func (g *Generator) eachTriangle(room model.Room, fn func(model.Cell)) {
	// vertical extent
	apexY := room.TopLeft.Y
	baseY := room.BottomRight.Y
	if baseY < apexY {
		// degenerate, just treat as rectangle
		g.eachRect(room, fn)
		return
	}

	height := float64(baseY - apexY)
	if height == 0 {
		// single row, again just a rectangle
		g.eachRect(room, fn)
		return
	}

	// horizontal center and maximum half-width
	cx := float64(room.TopLeft.X+room.BottomRight.X) / 2.0
	maxHalfWidth := float64(room.BottomRight.X-room.TopLeft.X) / 2.0

	for y := apexY; y <= baseY; y++ {
		// t goes from 0 at apex to 1 at base
		t := float64(y-apexY) / height
		halfWidth := maxHalfWidth * t
		minX := int32(math.Floor(cx - halfWidth))
		maxX := int32(math.Ceil(cx + halfWidth))
		for x := minX; x <= maxX; x++ {
			fn(model.Cell{X: x, Y: y})
		}
	}
}

// AddRoomEdges marks the cells of each room in the map
// according to its shape so that they appear in the drawn grid.
// avoids collisions with existing tiles.
func (g *Generator) AddRoomEdges(d *model.Dungeon, rooms []model.Room) {
	for _, room := range rooms {
		g.ForEachRoomCell(room, func(c model.Cell) {
			if d.At(c) == model.TileEmpty {
				d.Set(c, model.TileRoomFloor)
			}
		})
	}
}

func (g *Generator) ForEachRoomCell(room model.Room, fn func(model.Cell)) {
	switch room.Shape {
	case model.Rectangle, model.Square:
		g.eachRect(room, fn)
	case model.Circle:
		g.eachCircle(room, fn)
	case model.Triangle:
		g.eachTriangle(room, fn)
	default:
		g.eachRect(room, fn)
	}
}
