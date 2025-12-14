package main

import (
	"fmt"

	"github.com/mikegio27/proc-dungeons/dungeon"
	"github.com/mikegio27/proc-dungeons/grid"
)

func main() {
	fmt.Println("Thinking how to generate procedural dungeons in go...")
	grid.InitGrid(10, 10)
	dungeon.DrawGrid(dungeon.GenPaths())
}
