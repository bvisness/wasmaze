package main

import (
	"math/rand"
	"syscall/js"
)

type Coords struct {
	X, Y int
}

func (c *Coords) ToJS() js.Value {
	return js.ValueOf([]interface{}{c.X, c.Y})
}

type Cell struct {
	ConnectsTo []Coords
}

func (c *Cell) ToJS() js.Value {
	var jsCt []interface{}
	for _, ct := range c.ConnectsTo {
		jsCt = append(jsCt, ct.ToJS())
	}

	return js.ValueOf(map[string]interface{}{
		"connectsTo": jsCt,
	})
}

type Maze struct {
	Width, Height int
	Cells         [][]Cell
}

func (m *Maze) ToJS() js.Value {
	cells := make([]interface{}, m.Width)
	for x := 0; x < m.Width; x++ {
		col := make([]interface{}, m.Height)

		for y := 0; y < m.Height; y++ {
			col[y] = m.Cells[x][y].ToJS()
		}

		cells[x] = col
	}

	result := js.ValueOf(map[string]interface{}{
		"width":  m.Width,
		"height": m.Height,
		"cells":  cells,
	})

	return result
}

func genMazeGo(args []js.Value) {
	width := args[0].Int()
	height := args[1].Int()
	callback := args[2]

	cells := make([][]Cell, width)
	for x := 0; x < width; x++ {
		cells[x] = make([]Cell, height)

		for y := 0; y < height; y++ {
			cells[x][y] = Cell{}
		}
	}

	visited := make([][]bool, width)
	for x := 0; x < width; x++ {
		visited[x] = make([]bool, height)
	}

	start := Coords{rand.Intn(width), rand.Intn(height)}
	stack := []Coords{start}

	for len(stack) > 0 {
		c := stack[len(stack)-1]

		visited[c.X][c.Y] = true

		var unvisitedNeighbors []Coords
		for _, d := range []Coords{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			newX := c.X + d.X
			newY := c.Y + d.Y

			if newX < 0 || newX >= width || newY < 0 || newY >= height || visited[newX][newY] {
				continue
			} else {
				unvisitedNeighbors = append(unvisitedNeighbors, Coords{newX, newY})
			}
		}

		if len(unvisitedNeighbors) > 0 {
			n := unvisitedNeighbors[rand.Intn(len(unvisitedNeighbors))]

			currentCell := &cells[c.X][c.Y]
			currentCell.ConnectsTo = append(currentCell.ConnectsTo, n)

			neighborCell := &cells[n.X][n.Y]
			neighborCell.ConnectsTo = append(neighborCell.ConnectsTo, c)

			stack = append(stack, n)
		} else {
			stack = stack[:len(stack)-1]
		}
	}

	maze := Maze{
		Width:  width,
		Height: height,
		Cells:  cells,
	}

	callback.Invoke(maze.ToJS())
}

func registerCallbacks() {
	js.Global().Set("genMazeGo", js.NewCallback(genMazeGo))
}

func main() {
	c := make(chan struct{}, 0)

	println("WASM Go Initialized")
	// register functions
	registerCallbacks()
	<-c
}
