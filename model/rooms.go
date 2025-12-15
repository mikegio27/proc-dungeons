package model

type RoomId int

type Room struct {
	TopLeft     Cell
	BottomRight Cell
	Shape       RoomId
}

const (
	Rectangle RoomId = iota
	Circle
	Square
	Triangle
)

var shapeName = map[RoomId]string{
	Rectangle: "Rectangle",
	Circle:    "Circle",
	Square:    "Square",
	Triangle:  "Triangle",
}

// String implements fmt.Stringer for RoomId, returning the human-readable
// name of the shape.
func (id RoomId) String() string {
	if name, ok := shapeName[id]; ok {
		return name
	}
	return "Unknown"
}
