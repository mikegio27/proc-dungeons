package main

import (
	"fmt"
	"time"

	"github.com/mikegio27/proc-dungeons/generator"
	"github.com/mikegio27/proc-dungeons/model"
	"github.com/mikegio27/proc-dungeons/render"
)

func main() {
	fmt.Println("Procedurally generating dungeon...")
	// Use current time as seed for randomness
	// TODO: Allow user to specify seed via command-line argument
	seed := time.Now().UnixNano()
	fmt.Printf("Using seed: %d\n", seed)
	gridX := int32(50)
	gridY := int32(20)
	maxRooms := 20
	g := generator.New(generator.Config{
		Grid:         model.Grid{MaxX: gridX, MaxY: gridY, MinX: -gridX, MinY: -gridY},
		MaxRooms:     maxRooms,
		CorridorW:    2,
		CorridorBuff: 1,
		RoomShapes: []model.RoomId{
			model.Rectangle,
			model.Circle,
			model.Square,
			model.Triangle,
		},
	}, seed)
	d := g.Generate()
	render.DrawDungeon(&d)
	fmt.Printf("Rooms: %v\n", d.Rooms)
}
