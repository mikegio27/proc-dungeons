package main

import (
	"fmt"

	"github.com/mikegio27/proc-dungeons/geography"
)

func main() {
	fmt.Println("Procedurally generating dungeon...")
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
