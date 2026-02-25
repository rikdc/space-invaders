package main

import (
	"fmt"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	clearScreen = "\033[2J\033[H"
)

// Render draws the game to stdout.
func Render(g *Game) {
	// Build grid.
	grid := make([][]rune, Height)
	for i := range grid {
		grid[i] = make([]rune, Width)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	// Draw border.
	for x := 0; x < Width; x++ {
		grid[0][x] = '-'
		grid[Height-1][x] = '-'
	}
	for y := 0; y < Height; y++ {
		grid[y][0] = '|'
		grid[y][Width-1] = '|'
	}

	// Draw invaders.
	for r := 0; r < InvaderRows; r++ {
		for c := 0; c < InvaderCols; c++ {
			inv := g.Invaders[r][c]
			if inv.Active && inBounds(inv.Pos) {
				grid[inv.Pos.Y][inv.Pos.X] = 'W'
			}
		}
	}

	// Draw player.
	if inBounds(g.Player) {
		grid[g.Player.Y][g.Player.X] = 'A'
	}

	// Draw player bullet.
	if g.PlayerBullet.Active && inBounds(g.PlayerBullet.Pos) {
		grid[g.PlayerBullet.Pos.Y][g.PlayerBullet.Pos.X] = '|'
	}

	// Draw invader bullet.
	if g.InvaderBullet.Active && inBounds(g.InvaderBullet.Pos) {
		grid[g.InvaderBullet.Pos.Y][g.InvaderBullet.Pos.X] = 'v'
	}

	// Render to string with colors.
	var sb strings.Builder
	sb.WriteString(clearScreen)

	// HUD.
	sb.WriteString(fmt.Sprintf("%sScore: %d  Lives: %s%s\n",
		colorYellow, g.Score, livesStr(g.Lives), colorReset))

	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			ch := grid[y][x]
			switch ch {
			case 'W':
				sb.WriteString(colorGreen + string(ch) + colorReset)
			case 'A':
				sb.WriteString(colorCyan + string(ch) + colorReset)
			case '|':
				sb.WriteString(colorYellow + string(ch) + colorReset)
			case 'v':
				sb.WriteString(colorRed + string(ch) + colorReset)
			default:
				sb.WriteByte(byte(ch))
			}
		}
		sb.WriteByte('\n')
	}

	switch g.State {
	case StateWon:
		sb.WriteString(colorGreen + "\n*** YOU WIN! ***\n" + colorReset)
	case StateLost:
		sb.WriteString(colorRed + "\n*** GAME OVER ***\n" + colorReset)
	default:
		sb.WriteString("\nControls: A/D or arrow keys to move, SPACE to shoot, Q to quit\n")
	}

	fmt.Print(sb.String())
}

func inBounds(p Point) bool {
	return p.X > 0 && p.X < Width-1 && p.Y > 0 && p.Y < Height-1
}

func livesStr(lives int) string {
	s := ""
	for i := 0; i < lives; i++ {
		s += "A "
	}
	return s
}
