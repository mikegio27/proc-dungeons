package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/mikegio27/proc-dungeons/geography"
)

func main() {
	fmt.Println("Procedurally generating dungeon...")
	// Use a time-based seed for randomness. TODO: Allow user to specify seed.
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	fmt.Printf("Using seed: %d\n", seed.Int63())
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
