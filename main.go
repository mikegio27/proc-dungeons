package main

import (
	"fmt"

	"github.com/mikegio27/proc-dungeons/geography"
)

func main() {
	fmt.Println("Thinking how to generate procedural dungeons in go...")
	geography.InitGrid(50, 25)
	rooms := geography.Rooms()
	visited, starts := geography.GenPaths(rooms)
	geography.AddRoomEdges(visited, rooms)
	geography.DrawGrid(visited, starts)
	fmt.Printf("Rooms: %v\n", rooms)
}
