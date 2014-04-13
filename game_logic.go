package main

import (
	"math"
)

var StartingTiles int = 2
var goal Number = 10 // 2048

func SetGoal(n int) {
	val := math.Log2(float64(n))
	if val < 0 || math.Floor(val) != val {
		panic("Goal must be a positive power of 2.")
	}
	goal = Number(val - 1)
}

// Returns success.
func (g *Game) Play(d Direction) bool {
	if g.Move(d) {
		g.AddRandomTile()
		return true
	}
	switch d {
	case Left:
		if !g.Try(Right) && !g.Try(Up) && !g.Try(Down) {
			g.Lost = true
			g.Finished.Done()
		}
	case Right:
		if !g.Try(Left) && !g.Try(Up) && !g.Try(Down) {
			g.Lost = true
			g.Finished.Done()
		}
	case Up:
		if !g.Try(Left) && !g.Try(Right) && !g.Try(Down) {
			g.Lost = true
			g.Finished.Done()
		}
	case Down:
		if !g.Try(Left) && !g.Try(Right) && !g.Try(Up) {
			g.Lost = true
			g.Finished.Done()
		}
	default:
		panic("Unknown direction")
	}

	return false
}

// Returns success.
func (g *Game) Move(d Direction) bool {
	moved := false

	// Reset state.
	for _, row := range g.Grid {
		for _, tile := range row {
			if tile != nil {
				tile.Merged = false
			}
		}
	}

	var (
		rowStart = 0
		rowEnd   = g.Size
		rowDelta = 1
		colStart = 0
		colEnd   = g.Size
		colDelta = 1
	)

	switch d {
	case Left:
		rowStart = 0
		rowEnd = g.Size
		rowDelta = 1
	case Right:
		rowStart = g.Size - 1
		rowEnd = -1
		rowDelta = -1
	case Up:
		colStart = g.Size - 1
		colEnd = -1
		colDelta = -1
	case Down:
		colStart = 0
		colEnd = g.Size
		colDelta = 1
	}

	for x := rowStart; x != rowEnd; x += rowDelta {
		for y := colStart; y != colEnd; y += colDelta {
			//fmt.Println(x, y, g.Size)
			tile := g.Grid[x][y]
			if tile == nil {
				continue
			}

			// Move if possible.
			pos := &Index{x, y}
			for pos.Increment(d, g.Size) {
				// Look for a space ahead.
				if g.Grid[pos.X][pos.Y] != nil {
					pos.Decrement(d, g.Size)
					break
				}
			}

			if pos.X != x || pos.Y != y {
				moved = true
			}

			// Look for a merge ahead.
			next := pos.Next(d, g.Size)
			if next != nil {
				nextTile := g.Grid[next.X][next.Y]
				if nextTile != nil && nextTile.Number == tile.Number && !nextTile.Merged {
					// Successful merge.
					nextTile.Merged = true
					nextTile.Number++
					if nextTile.Number > g.Max {
						g.Max++
					}
					if nextTile.Number == goal {
						g.Won = true
						g.Finished.Done()
					}
					g.Score += tile.Number.Int()
					tile = nil
					moved = true
				}
			}

			// Remove old position.
			g.Grid[x][y] = nil

			// Add new position if we still exist.
			// This may be the position we just cleared.
			if tile != nil {
				g.Grid[pos.X][pos.Y] = tile
			}
		}
	}

	return moved
}

// Return success.
func (g *Game) Try(d Direction) bool {
	// Reset state.
	for _, row := range g.Grid {
		for _, tile := range row {
			if tile != nil {
				tile.Merged = false
			}
		}
	}

	var (
		rowStart = 0
		rowEnd   = g.Size
		rowDelta = 1
		colStart = 0
		colEnd   = g.Size
		colDelta = 1
	)

	switch d {
	case Left:
		rowStart = 0
		rowEnd = g.Size
		rowDelta = 1
	case Right:
		rowStart = g.Size - 1
		rowEnd = -1
		rowDelta = -1
	case Up:
		colStart = g.Size - 1
		colEnd = -1
		colDelta = -1
	case Down:
		colStart = 0
		colEnd = g.Size
		colDelta = 1
	}

	for x := rowStart; x != rowEnd; x += rowDelta {
		for y := colStart; y != colEnd; y += colDelta {
			//fmt.Println(x, y, g.Size)
			tile := g.Grid[x][y]
			if tile == nil {
				continue
			}

			// Move if possible.
			pos := &Index{x, y}
			for pos.Increment(d, g.Size) {
				// Look for a space ahead.
				if g.Grid[pos.X][pos.Y] != nil {
					pos.Decrement(d, g.Size)
					break
				}
			}

			if pos.X != x || pos.Y != y {
				return true
			}

			// Look for a merge ahead.
			next := pos.Next(d, g.Size)
			if next != nil {
				nextTile := g.Grid[next.X][next.Y]
				if nextTile != nil && nextTile.Number == tile.Number && !nextTile.Merged {
					// Successful merge.
					return true
				}
			}
		}
	}

	return false
}
