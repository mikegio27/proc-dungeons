package dungeon

import (
	"fmt"
	"math/rand"

	"github.com/mikegio27/proc-dungeons/grid"
)

func fiftyFifty() bool {
	randVal := rand.Intn(10000)
	return randVal%2 == 0
}

func randSign() int64 {
	if fiftyFifty() {
		return 1
	}
	return -1
}

func genStartingPoint() (int64, int64) {
	plane := grid.GetPlane()

	switch rand.Intn(4) {
	case 0:
		return plane.MinX, rand.Int63n(plane.MaxY+1) * randSign()
	case 1:
		return plane.MaxX, rand.Int63n(plane.MaxY+1) * randSign()
	case 2:
		return rand.Int63n(plane.MaxX+1) * randSign(), plane.MinY
	default:
		return rand.Int63n(plane.MaxX+1) * randSign(), plane.MaxY
	}
}

func nextPoint(p grid.Path, visited map[grid.Cell]bool) (int64, int64, bool) {
	plane := grid.GetPlane()
	dxdy := [][2]int64{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

	rand.Shuffle(len(dxdy), func(i int, j int) { dxdy[i], dxdy[j] = dxdy[j], dxdy[i] })
	for _, d := range dxdy {
		nx := p.Start.X + d[0]
		ny := p.Start.Y + d[1]

		if nx < plane.MinX || nx > plane.MaxX || ny < plane.MinY || ny > plane.MaxY {
			continue
		}
		pt := grid.Cell{X: nx, Y: ny}
		if visited[pt] {
			continue
		}
		return nx, ny, true
	}

	// no available moves
	return p.Start.X, p.Start.Y, false
}

func GenPaths() map[grid.Cell]bool {
	startX, startY := genStartingPoint()
	p := grid.Path{Start: grid.Cell{X: startX, Y: startY}}

	visited := make(map[grid.Cell]bool)
	visited[grid.Cell{X: startX, Y: startY}] = true

	fmt.Printf("Starting point: %d, %d\n", p.Start.X, p.Start.Y)
	for i := 0; i < 50; i++ { // longer path
		nextX, nextY, ok := nextPoint(p, visited)
		if !ok {
			fmt.Println("No more moves, stopping.")
			break
		}
		p.Cells = append(p.Cells, grid.Cell{X: nextX, Y: nextY})
		p.Start.X = nextX
		p.Start.Y = nextY
		visited[grid.Cell{X: nextX, Y: nextY}] = true
		fmt.Printf("Next point: %d, %d\n", nextX, nextY)
	}

	fmt.Println("Generated path:", p)
	return visited
}

func DrawGrid(visited map[grid.Cell]bool) {
	plane := grid.GetPlane()
	for y := plane.MaxY; y >= plane.MinY; y-- {
		for x := plane.MinX; x <= plane.MaxX; x++ {
			ch := '.'
			pt := grid.Cell{X: x, Y: y}
			if visited[pt] {
				ch = '#'
			}
			if x == 0 && y == 0 {
				ch = '+'
			}
			fmt.Printf("%c ", ch)
		}
		fmt.Println()
	}
}
