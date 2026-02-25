package main

import "math/rand"

const (
	Width        = 40
	Height       = 20
	InvaderRows  = 3
	InvaderCols  = 8
	PlayerLives  = 3
	BulletSpeed  = 1
	InvaderSteps = 15 // ticks between invader moves
)

type Point struct {
	X, Y int
}

type Bullet struct {
	Pos    Point
	Active bool
	// Direction: -1 = up (player), +1 = down (invader)
	Dir int
}

type Invader struct {
	Pos    Point
	Active bool
}

type GameState int

const (
	StatePlaying GameState = iota
	StateWon
	StateLost
)

type Game struct {
	Player        Point
	PlayerBullet  Bullet
	Invaders      [InvaderRows][InvaderCols]Invader
	InvaderBullet Bullet
	Score         int
	Lives         int
	Tick          int
	InvaderDir    int // 1 = right, -1 = left
	State         GameState
}

func NewGame() *Game {
	g := &Game{
		Player:     Point{Width / 2, Height - 2},
		Lives:      PlayerLives,
		InvaderDir: 1,
	}
	for r := 0; r < InvaderRows; r++ {
		for c := 0; c < InvaderCols; c++ {
			g.Invaders[r][c] = Invader{
				Pos:    Point{3 + c*4, 2 + r*2},
				Active: true,
			}
		}
	}
	return g
}

// MoveLeft moves the player left.
func (g *Game) MoveLeft() {
	if g.Player.X > 1 {
		g.Player.X--
	}
}

// MoveRight moves the player right.
func (g *Game) MoveRight() {
	if g.Player.X < Width-2 {
		g.Player.X++
	}
}

// Shoot fires a player bullet if none is active.
func (g *Game) Shoot() {
	if !g.PlayerBullet.Active {
		g.PlayerBullet = Bullet{
			Pos:    Point{g.Player.X, g.Player.Y - 1},
			Active: true,
			Dir:    -1,
		}
	}
}

// Update advances the game by one tick.
func (g *Game) Update() {
	if g.State != StatePlaying {
		return
	}
	g.Tick++

	g.movePlayerBullet()
	g.moveInvaderBullet()

	if g.Tick%InvaderSteps == 0 {
		g.moveInvaders()
		g.maybeInvaderShoot()
	}

	g.checkCollisions()
	g.checkWinLoss()
}

func (g *Game) movePlayerBullet() {
	if !g.PlayerBullet.Active {
		return
	}
	g.PlayerBullet.Pos.Y += g.PlayerBullet.Dir
	if g.PlayerBullet.Pos.Y < 1 {
		g.PlayerBullet.Active = false
	}
}

func (g *Game) moveInvaderBullet() {
	if !g.InvaderBullet.Active {
		return
	}
	g.InvaderBullet.Pos.Y += g.InvaderBullet.Dir
	if g.InvaderBullet.Pos.Y >= Height-1 {
		g.InvaderBullet.Active = false
	}
}

func (g *Game) moveInvaders() {
	// Check if any invader would go out of bounds.
	needDrop := false
	for r := 0; r < InvaderRows; r++ {
		for c := 0; c < InvaderCols; c++ {
			inv := &g.Invaders[r][c]
			if !inv.Active {
				continue
			}
			nx := inv.Pos.X + g.InvaderDir
			if nx <= 0 || nx >= Width-1 {
				needDrop = true
			}
		}
	}

	if needDrop {
		g.InvaderDir = -g.InvaderDir
		for r := 0; r < InvaderRows; r++ {
			for c := 0; c < InvaderCols; c++ {
				if g.Invaders[r][c].Active {
					g.Invaders[r][c].Pos.Y++
				}
			}
		}
	} else {
		for r := 0; r < InvaderRows; r++ {
			for c := 0; c < InvaderCols; c++ {
				if g.Invaders[r][c].Active {
					g.Invaders[r][c].Pos.X += g.InvaderDir
				}
			}
		}
	}
}

func (g *Game) maybeInvaderShoot() {
	if g.InvaderBullet.Active {
		return
	}
	// Collect active invaders.
	var active []Point
	for r := 0; r < InvaderRows; r++ {
		for c := 0; c < InvaderCols; c++ {
			if g.Invaders[r][c].Active {
				active = append(active, g.Invaders[r][c].Pos)
			}
		}
	}
	if len(active) == 0 {
		return
	}
	shooter := active[rand.Intn(len(active))]
	g.InvaderBullet = Bullet{
		Pos:    Point{shooter.X, shooter.Y + 1},
		Active: true,
		Dir:    1,
	}
}

func (g *Game) checkCollisions() {
	// Player bullet hits invader.
	if g.PlayerBullet.Active {
		for r := 0; r < InvaderRows; r++ {
			for c := 0; c < InvaderCols; c++ {
				inv := &g.Invaders[r][c]
				if inv.Active && inv.Pos == g.PlayerBullet.Pos {
					inv.Active = false
					g.PlayerBullet.Active = false
					g.Score += 10
				}
			}
		}
	}

	// Invader bullet hits player.
	if g.InvaderBullet.Active && g.InvaderBullet.Pos == g.Player {
		g.InvaderBullet.Active = false
		g.Lives--
	}
}

func (g *Game) checkWinLoss() {
	if g.Lives <= 0 {
		g.State = StateLost
		return
	}

	// Check if any invader reached the player row.
	for r := 0; r < InvaderRows; r++ {
		for c := 0; c < InvaderCols; c++ {
			if g.Invaders[r][c].Active && g.Invaders[r][c].Pos.Y >= g.Player.Y {
				g.State = StateLost
				return
			}
		}
	}

	// Check if all invaders destroyed.
	for r := 0; r < InvaderRows; r++ {
		for c := 0; c < InvaderCols; c++ {
			if g.Invaders[r][c].Active {
				return
			}
		}
	}
	g.State = StateWon
}

// ActiveInvaderCount returns the number of active invaders.
func (g *Game) ActiveInvaderCount() int {
	n := 0
	for r := 0; r < InvaderRows; r++ {
		for c := 0; c < InvaderCols; c++ {
			if g.Invaders[r][c].Active {
				n++
			}
		}
	}
	return n
}
