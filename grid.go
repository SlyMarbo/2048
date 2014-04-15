package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type Game struct {
	Size     int
	Grid     [][]*Tile
	Max      Number
	Won      bool
	Lost     bool
	Score    int
	Finished *sync.WaitGroup
}

func NewGame(size int) *Game {
	game := new(Game)
	game.Size = size
	game.Grid = make([][]*Tile, size)
	for i := range game.Grid {
		game.Grid[i] = make([]*Tile, size)
	}
	if Lucky {
		StartingTiles = 1
	}
	for i := 0; i < StartingTiles; i++ {
		game.AddRandomTile()
	}
	return game
}

func (g *Game) AddRandomTile() {
	pos := g.RandomAvailableCell()
	if pos == nil {
		return
	}

	if Lucky {
		pos = &Index{0, g.Size - 1}
		if g.Grid[pos.X][pos.Y] != nil {
			return
		}
	}

	var num Number                    // == 2
	if rand.Intn(100) < 10 || Lucky { // 10% chance
		num++ // == 4
	}

	g.Grid[pos.X][pos.Y] = &Tile{Number: num}
}

func (g *Game) AvailableCells() []Index {
	out := make([]Index, 0, g.Size*g.Size-2)
	for x, row := range g.Grid {
		for y, cell := range row {
			if cell == nil {
				out = append(out, Index{x, y})
			}
		}
	}
	return out
}

func (g *Game) String() string {
	buf := new(bytes.Buffer)
	for y := g.Size - 1; y >= 0; y-- {
		fmt.Fprintf(buf, " %s \n", strings.Repeat("---- ", g.Size))
		for x := 0; x < g.Size; x++ {
			var num, spaceL, spaceR string
			if val := g.Grid[x][y]; val != nil {
				num = val.String()
			}
			space := 4 - len(num)
			i := int(math.Floor(float64(space) / 2))
			if i > 0 {
				spaceL = strings.Repeat(" ", i)
			}
			if j := space - i; j > 0 {
				spaceR = strings.Repeat(" ", j)
			}
			fmt.Fprintf(buf, "|%s%s%s", spaceL, num, spaceR)
		}
		buf.WriteString("|\n")
	}
	fmt.Fprintf(buf, " %s", strings.Repeat("---- ", g.Size))
	return buf.String()
}

func (g *Game) RandomAvailableCell() *Index {
	available := g.AvailableCells()
	if len(available) == 0 {
		return nil
	}
	choice := available[rand.Intn(len(available))]
	return &choice
}

func (g *Game) Wait() {
	g.Finished = new(sync.WaitGroup)
	g.Finished.Add(1)
	g.Finished.Wait()
}
