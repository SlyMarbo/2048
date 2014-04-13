package main

import (
	"flag"
	"fmt"
	"os"
)

var term *Terminal
var game *Game
var shouldReset = false

func main() {
	var size = flag.Int("size", 4, "Width/height of the grid.")
	var goal = flag.Int("goal", 2048, "Target score.")

	flag.Parse()

	if *size < 2 {
		fmt.Fprintln(os.Stderr, "Error: size must be at least 2.")
		os.Exit(1)
	}
	if *goal < 4 {
		fmt.Fprintln(os.Stderr, "Error: goal must be at least 4.")
		os.Exit(1)
	}
	SetGoal(*goal)

	var err error
	term, err = NewTerminal(os.Stdin, os.Stdout, "")
	handle(err)
	defer term.Reset()
	term.Println("Controls: wasd or arrow keys to play, enter to exit.")

	game = NewGame(*size)
	update()

	term.KeyCallback = doKeypress

	go func() {
		term.ReadLine()
		game.Finished.Done()
	}()

	game.Finished.Wait()

	if game.Won {
		term.Println("\nCongratulations!")
	} else if game.Lost {
		term.Println("\nYou lose.")
	}
}

func doKeypress(key int) {
	var dir Direction
	switch key {
	case 'w', KeyUp:
		dir = Up
	case 'a', KeyLeft:
		dir = Left
	case 's', KeyDown:
		dir = Down
	case 'd', KeyRight:
		dir = Right
	default:
		return
	}

	// Try the move.
	if game.Play(dir) {
		// Move was successful.
		update()
	}
}

func update() {
	if shouldReset {
		handle(reset())
	}
	term.Println("Score:", game.Score)
	term.Print(game)
	shouldReset = true
}

func reset() error {
	_, err := term.Write([]byte{'\r'})
	if err != nil {
		return err
	}
	term.ClearLine()
	for i := 0; i < game.Size*2+1; i++ {
		_, err := term.Write([]byte{keyEscape, '[', 'A'})
		if err != nil {
			return err
		}
		term.ClearLine()
	}
	return nil
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}
