package generator

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/mikegio27/proc-dungeons/model"
)

type Config struct {
	Grid       model.Grid
	MaxRooms   int
	RoomShapes []model.RoomId
	RoomMinW   int32
	RoomMaxW   int32
	RoomMinH   int32
	RoomMaxH   int32
	CorridorW  int32
}

type Generator struct {
	cfg Config
	rng *rand.Rand
}

func New(cfg Config, seed int64) *Generator {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	return &Generator{
		cfg: cfg,
		rng: rand.New(rand.NewSource(seed)),
	}
}

func (g *Generator) Generate() {
	rooms := g.Rooms(g.cfg.MaxRooms)
	visited, starts := g.GenPaths(rooms)
	g.AddRoomEdges(visited, rooms)
	g.DrawGrid(visited, starts)
	fmt.Printf("Rooms: %v\n", rooms)
}
