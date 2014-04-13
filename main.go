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
		complain("Error: size must be at least 2.")
	}
	if *goal < 4 {
		complain("Error: goal must be at least 4.")
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

	game.Wait()

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
		update() // Move was successful.
	}
}

func update() {
	if shouldReset {
		_, err := term.Write([]byte{'\r'})
		handle(err)
		term.ClearLine()
		for i := 0; i < game.Size*2+1; i++ {
			_, err = term.Write([]byte{keyEscape, '[', 'A'})
			handle(err)
			term.ClearLine()
		}
	}
	term.Println("Score:", game.Score)
	term.Print(game)
	shouldReset = true
}

func complain(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}
