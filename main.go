package main

import (
	"fmt"
	"time"

	"github.com/mikegio27/proc-dungeons/geography"
)

func main() {
	fmt.Println("Procedurally generating dungeon...")
	// Use current time as seed for randomness
	// TODO: Allow user to specify seed via command-line argument
	seed := time.Now().UnixNano()
	geography.SetRandSeed(seed)
	fmt.Printf("Using seed: %d\n", seed)
	gridX := int32(20)
	gridY := int32(20)
	maxRooms := 10
	geography.InitGrid(gridX, gridY)
	rooms := geography.Rooms(maxRooms)
	visited, starts := geography.GenPaths(rooms)
	geography.AddRoomEdges(visited, rooms)
	geography.DrawGrid(visited, starts)
	fmt.Printf("Rooms: %v\n", rooms)
}
