package main

import (
	"os"
)

const esc = 27

var size = 4

var term *Terminal
var game *Game
var shouldReset = false

func main() {
	var err error
	term, err = NewTerminal(os.Stdin, os.Stdout, "")
	handle(err)
	defer term.Reset()
	term.Println("Controls: wasd or arrow keys.")

	game = NewGame(size)
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
	for i := 0; i < size*2+1; i++ {
		_, err := term.Write([]byte{esc, '[', 'A'})
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
