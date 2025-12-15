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

func (g *Generator) Generate() model.Dungeon {
	d := model.NewDungeon(g.cfg.Grid)
	rooms := g.Rooms(g.cfg.MaxRooms)
	d.Rooms = rooms
	starts := g.GenPaths(&d, rooms)
	d.Starts = starts
	g.AddRoomEdges(&d, rooms)

	return d
}

func (g *Generator) DrawDungeon(d model.Dungeon) {
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
			if x == 0 && y == 0 {
				ch = '+'
			}
			fmt.Printf("%c ", ch)
		}
		fmt.Println()
	}
}
