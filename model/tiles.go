package model

// Tile represents the type of a cell in the dungeon grid.
type Tile uint8

const (
	TileEmpty Tile = iota
	TileRoomFloor
	TileCorridor
	TileDoor
	TileWall
)

func (t Tile) String() string {
	switch t {
	case TileEmpty:
		return "Empty"
	case TileRoomFloor:
		return "RoomFloor"
	case TileCorridor:
		return "Corridor"
	case TileDoor:
		return "Door"
	default:
		return "Unknown"
	}
}

// Rune is an ASCII glyph for a simple renderer.
// TODO: move this out into a renderer package later if it makes sense.
func (t Tile) Rune() rune {
	switch t {
	case TileEmpty:
		return ' '
	case TileRoomFloor:
		return '.'
	case TileCorridor:
		return '#'
	case TileDoor:
		return '+'
	case TileWall:
		return '▒'
	default:
		return '?'
	}
}
