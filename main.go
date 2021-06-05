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
	var firstTime bool = true
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(colornames.Black)
		world.Draw(win)

		if firstTime {

			logoImg, err := spacegame.NewloadPicture("assets/textures/logo.png")
			logo := pixel.NewSprite(logoImg, logoImg.Bounds())
			mat := pixel.IM
			mat = mat.Moved(pixel.V(win.Bounds().Center().X, win.Bounds().Center().Y+175))
			mat = mat.Scaled(pixel.V(win.Bounds().Center().X, win.Bounds().Center().Y), 0.5)
			logo.Draw(win, mat)
			if err != nil {
				panic(err)
			}

			tvMenu := text.New(pixel.V(windowWidth/2-150, windowHeight/2-175), basicAtlas)
			fmt.Fprintln(tvMenu, "Press ENTER to start")
			tvMenu.Draw(win, pixel.IM.Scaled(tvMenu.Orig, 2))
			if win.JustPressed(pixelgl.KeyEnter) {
				firstTime = false
			}
		}

		if isRunning && !firstTime {
			if win.Pressed(pixelgl.KeyLeft) {
				direction = spacegame.LeftDirection
			}

			if win.Pressed(pixelgl.KeyRight) {
				direction = spacegame.RightDirection
			}
			if win.Pressed(pixelgl.KeySpace) {
				action = spacegame.ShootAction
			}
			//50 is the standard for alien number
			spacegame.NewCreateEnemies(win, 50)
			coordenadasFortalezas := spacegame.NewCreateFortress(win)
			fmt.Println(coordenadasFortalezas[1])

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
		} else if !firstTime {
			tvPause := text.New(pixel.V(windowWidth/2-70, windowHeight/2), basicAtlas)
			fmt.Fprintln(tvPause, "Paused")
			tvPause.Draw(win, pixel.IM.Scaled(tvPause.Orig, 4))
		}

		if win.JustPressed(pixelgl.KeyP) && !firstTime {
			isRunning = !isRunning
		}
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
