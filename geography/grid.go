package geography

type Grid struct {
	MaxX int64
	MaxY int64
	MinY int64
	MinX int64
}

var current Grid

func InitGrid(x, y int64) {
	current = Grid{MaxX: x, MaxY: y, MinX: -x, MinY: -y}
}

func GetPlane() Grid {
	return current
}
