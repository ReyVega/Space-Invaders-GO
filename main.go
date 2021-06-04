package main

import (
	"fmt"
	_ "image/png"
	"log"
	spacegame "spaceInvaders/libs"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
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
	var isRunning bool = true
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(colornames.Black)
		world.Draw(win)

		if isRunning {
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

			tvScore := text.New(pixel.V(20, 570), basicAtlas)
			tvLives := text.New(pixel.V(690, 570), basicAtlas)
			fmt.Fprintln(tvScore, "Score: ", 0)
			fmt.Fprintln(tvLives, "Lives: ", player.GetLife())
			tvScore.Draw(win, pixel.IM.Scaled(tvScore.Orig, 1.5))
			tvLives.Draw(win, pixel.IM.Scaled(tvLives.Orig, 1.5))
		} else {
			tvPause := text.New(pixel.V(windowWidth/2-70, windowHeight/2), basicAtlas)
			fmt.Fprintln(tvPause, "Paused")
			tvPause.Draw(win, pixel.IM.Scaled(tvPause.Orig, 4))
		}

		if win.JustPressed(pixelgl.KeyP) {
			isRunning = !isRunning
		}
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
