package main

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/term"
)

const tickDuration = 80 * time.Millisecond

// key constants for mapped input.
type key byte

const (
	keyLeft  key = 1
	keyRight key = 2
	keyShoot key = 3
	keyQuit  key = 4
)

func main() {
	// Put terminal into raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to set raw mode: %v\n", err)
		os.Exit(1)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	game := NewGame()

	inputCh := make(chan key, 10)
	go readInput(inputCh)

	ticker := time.NewTicker(tickDuration)
	defer ticker.Stop()

	Render(game)

	for {
		select {
		case k := <-inputCh:
			handleInput(game, k)
			if game.State != StatePlaying {
				Render(game)
				time.Sleep(2 * time.Second)
				return
			}
		case <-ticker.C:
			game.Update()
			Render(game)
			if game.State != StatePlaying {
				time.Sleep(2 * time.Second)
				return
			}
		}
	}
}

// readInput reads from stdin and maps to key constants.
func readInput(ch chan<- key) {
	buf := make([]byte, 4)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			return
		}
		b := buf[0]
		switch {
		case b == 'q' || b == 'Q' || b == 3: // 3 = Ctrl+C
			ch <- keyQuit
		case b == 'a' || b == 'A':
			ch <- keyLeft
		case b == 'd' || b == 'D':
			ch <- keyRight
		case b == ' ':
			ch <- keyShoot
		case b == 27 && n >= 3 && buf[1] == '[':
			// ANSI escape sequence for arrow keys.
			switch buf[2] {
			case 'C': // right arrow
				ch <- keyRight
			case 'D': // left arrow
				ch <- keyLeft
			}
		}
	}
}

// handleInput processes a mapped key press.
func handleInput(g *Game, k key) {
	switch k {
	case keyQuit:
		os.Exit(0)
	case keyLeft:
		g.MoveLeft()
	case keyRight:
		g.MoveRight()
	case keyShoot:
		g.Shoot()
	}
}
