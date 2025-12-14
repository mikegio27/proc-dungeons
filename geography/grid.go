package geography

type Grid struct {
	MaxX int32
	MaxY int32
	MinY int32
	MinX int32
}

var current Grid

func InitGrid(x, y int32) {
	current = Grid{MaxX: x, MaxY: y, MinX: -x, MinY: -y}
}

func GetPlane() Grid {
	return current
}
