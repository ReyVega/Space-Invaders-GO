package main

import (
	_ "image/png"
	"log"
	spacegame "spaceInvaders/libs"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	windowWidth  = 800
	windowHeight = 600
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Space Invaders",
		Bounds: pixel.R(0, 0, windowWidth, windowHeight),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		log.Fatal(err)
	}

	world := spacegame.NewWorld(windowWidth, windowHeight)
	if err := world.AddBackground("assets/textures/background.png"); err != nil {
		log.Fatal(err)
	}

	player, err := spacegame.NewPlayer("assets/textures/ship.png", 5, world)
	if err != nil {
		log.Fatal(err)
	}

	direction := spacegame.Idle
	last := time.Now()
	action := spacegame.NoneAction

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(colornames.Black)
		world.Draw(win)

		if win.Pressed(pixelgl.KeyLeft) {
			direction = spacegame.LeftDirection
		}

		if win.Pressed(pixelgl.KeyRight) {
			direction = spacegame.RightDirection
		}
		if win.Pressed(pixelgl.KeySpace) {
			action = spacegame.ShootAction
		}

		spacegame.NewCreateEnemies(win)
		spacegame.NewCreateFortress(win)

		player.Update(direction, action, dt)
		player.Draw(win)
		direction = spacegame.Idle
		action = spacegame.NoneAction

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
