package main

import (
	"testing"
)

func TestNewGame(t *testing.T) {
	g := NewGame()

	if g.Lives != PlayerLives {
		t.Errorf("expected %d lives, got %d", PlayerLives, g.Lives)
	}
	if g.Score != 0 {
		t.Errorf("expected score 0, got %d", g.Score)
	}
	if g.State != StatePlaying {
		t.Errorf("expected StatePlaying, got %v", g.State)
	}
	if g.InvaderDir != 1 {
		t.Errorf("expected initial invader direction 1, got %d", g.InvaderDir)
	}
	count := g.ActiveInvaderCount()
	expected := InvaderRows * InvaderCols
	if count != expected {
		t.Errorf("expected %d invaders, got %d", expected, count)
	}
}

func TestPlayerMovement(t *testing.T) {
	g := NewGame()
	startX := g.Player.X

	g.MoveLeft()
	if g.Player.X != startX-1 {
		t.Errorf("expected player at %d after move left, got %d", startX-1, g.Player.X)
	}

	g.MoveRight()
	if g.Player.X != startX {
		t.Errorf("expected player at %d after move right, got %d", startX, g.Player.X)
	}
}

func TestPlayerMovementBoundaries(t *testing.T) {
	g := NewGame()

	// Move all the way to the left boundary.
	for i := 0; i < Width; i++ {
		g.MoveLeft()
	}
	if g.Player.X < 1 {
		t.Errorf("player moved past left boundary: x=%d", g.Player.X)
	}

	// Move all the way to the right boundary.
	for i := 0; i < Width; i++ {
		g.MoveRight()
	}
	if g.Player.X >= Width-1 {
		t.Errorf("player moved past right boundary: x=%d", g.Player.X)
	}
}

func TestShoot(t *testing.T) {
	g := NewGame()

	if g.PlayerBullet.Active {
		t.Error("bullet should not be active initially")
	}

	g.Shoot()
	if !g.PlayerBullet.Active {
		t.Error("bullet should be active after shooting")
	}
	if g.PlayerBullet.Pos.X != g.Player.X {
		t.Errorf("bullet X should match player X: got %d, want %d", g.PlayerBullet.Pos.X, g.Player.X)
	}
	if g.PlayerBullet.Dir != -1 {
		t.Errorf("player bullet should travel up (dir=-1), got %d", g.PlayerBullet.Dir)
	}
}

func TestShootOnlyOneBulletAtATime(t *testing.T) {
	g := NewGame()
	g.Shoot()
	bulletY := g.PlayerBullet.Pos.Y

	// Advance bullet a bit.
	g.PlayerBullet.Pos.Y -= 3

	// Second shoot should not fire because one is already active.
	g.Shoot()
	if g.PlayerBullet.Pos.Y != bulletY-3 {
		t.Error("shooting while bullet active should not reset bullet position")
	}
}

func TestPlayerBulletMovesUp(t *testing.T) {
	g := NewGame()
	g.Shoot()
	startY := g.PlayerBullet.Pos.Y

	// Directly call the update enough times to move the bullet.
	for i := 0; i < InvaderSteps+1; i++ {
		g.Update()
	}

	if g.PlayerBullet.Active && g.PlayerBullet.Pos.Y >= startY {
		t.Errorf("player bullet should have moved up from y=%d", startY)
	}
}

func TestPlayerBulletDeactivatesAtTopBorder(t *testing.T) {
	g := NewGame()
	g.Shoot()
	g.PlayerBullet.Pos.Y = 2 // Near the top.

	// Run updates until bullet reaches top.
	for i := 0; i < 5; i++ {
		g.movePlayerBullet()
	}

	if g.PlayerBullet.Active {
		t.Error("bullet should be deactivated after hitting top border")
	}
}

func TestPlayerBulletHitsInvader(t *testing.T) {
	g := NewGame()

	// Place player bullet directly on an invader.
	invaderPos := g.Invaders[0][0].Pos
	g.PlayerBullet = Bullet{Pos: invaderPos, Active: true, Dir: -1}

	g.checkCollisions()

	if g.Invaders[0][0].Active {
		t.Error("invader should be destroyed after bullet hit")
	}
	if g.PlayerBullet.Active {
		t.Error("bullet should be deactivated after hitting invader")
	}
	if g.Score != 10 {
		t.Errorf("expected score 10 after killing invader, got %d", g.Score)
	}
}

func TestInvaderBulletHitsPlayer(t *testing.T) {
	g := NewGame()
	initialLives := g.Lives

	// Place invader bullet directly on player.
	g.InvaderBullet = Bullet{Pos: g.Player, Active: true, Dir: 1}

	g.checkCollisions()

	if g.InvaderBullet.Active {
		t.Error("invader bullet should be deactivated after hitting player")
	}
	if g.Lives != initialLives-1 {
		t.Errorf("expected %d lives after hit, got %d", initialLives-1, g.Lives)
	}
}

func TestWinCondition(t *testing.T) {
	g := NewGame()

	// Destroy all invaders.
	for r := 0; r < InvaderRows; r++ {
		for c := 0; c < InvaderCols; c++ {
			g.Invaders[r][c].Active = false
		}
	}

	g.checkWinLoss()

	if g.State != StateWon {
		t.Errorf("expected StateWon after all invaders destroyed, got %v", g.State)
	}
}

func TestLossConditionNoLives(t *testing.T) {
	g := NewGame()
	g.Lives = 0

	g.checkWinLoss()

	if g.State != StateLost {
		t.Errorf("expected StateLost with 0 lives, got %v", g.State)
	}
}

func TestLossConditionInvaderReachesPlayer(t *testing.T) {
	g := NewGame()

	// Move an invader to the player's row.
	g.Invaders[0][0].Pos.Y = g.Player.Y

	g.checkWinLoss()

	if g.State != StateLost {
		t.Errorf("expected StateLost when invader reaches player row, got %v", g.State)
	}
}

func TestInvadersMoveSideways(t *testing.T) {
	g := NewGame()
	startX := g.Invaders[0][0].Pos.X

	g.moveInvaders()

	newX := g.Invaders[0][0].Pos.X
	if newX != startX+g.InvaderDir && newX != startX-g.InvaderDir {
		t.Errorf("invader should have moved: was %d, now %d", startX, newX)
	}
}

func TestInvadersReverseAndDropOnBoundary(t *testing.T) {
	g := NewGame()

	// Force invaders to the right boundary.
	for r := 0; r < InvaderRows; r++ {
		for c := 0; c < InvaderCols; c++ {
			g.Invaders[r][c].Pos.X = Width - 2
		}
	}
	g.InvaderDir = 1
	startY := g.Invaders[0][0].Pos.Y

	g.moveInvaders()

	if g.InvaderDir != -1 {
		t.Errorf("invader direction should have reversed, got %d", g.InvaderDir)
	}
	if g.Invaders[0][0].Pos.Y != startY+1 {
		t.Errorf("invaders should have dropped one row: was %d, now %d", startY, g.Invaders[0][0].Pos.Y)
	}
}

func TestActiveInvaderCount(t *testing.T) {
	g := NewGame()

	total := InvaderRows * InvaderCols
	if g.ActiveInvaderCount() != total {
		t.Errorf("expected %d active invaders, got %d", total, g.ActiveInvaderCount())
	}

	g.Invaders[0][0].Active = false
	if g.ActiveInvaderCount() != total-1 {
		t.Errorf("expected %d active invaders after kill, got %d", total-1, g.ActiveInvaderCount())
	}
}

func TestUpdateDoesNothingWhenNotPlaying(t *testing.T) {
	g := NewGame()
	g.State = StateLost
	tick := g.Tick

	g.Update()

	if g.Tick != tick {
		t.Error("Update should not advance tick when game is not playing")
	}
}

func TestInvaderBulletMoves(t *testing.T) {
	g := NewGame()
	g.InvaderBullet = Bullet{Pos: Point{10, 5}, Active: true, Dir: 1}

	g.moveInvaderBullet()

	if g.InvaderBullet.Pos.Y != 6 {
		t.Errorf("invader bullet should have moved down to y=6, got y=%d", g.InvaderBullet.Pos.Y)
	}
}

func TestInvaderBulletDeactivatesAtBottom(t *testing.T) {
	g := NewGame()
	g.InvaderBullet = Bullet{Pos: Point{10, Height - 2}, Active: true, Dir: 1}

	g.moveInvaderBullet()

	if g.InvaderBullet.Active {
		t.Error("invader bullet should deactivate when reaching bottom")
	}
}
