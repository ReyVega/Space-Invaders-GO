package main

import (
<<<<<<< Updated upstream
	"log"

=======
>>>>>>> Stashed changes
	"image"
	_ "image/png"
	"log"
	"os"
	spacegame "spaceInvaders/libs"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	windowWidth  = 800
	windowHeight = 600
)

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func createEnemies(window *pixelgl.Window) {
	pic1, err := loadPicture("assets/textures/spritealien1.png")
	if err != nil {
		panic(err)
	}
	pic2, err := loadPicture("assets/textures/spritealien2.png")
	if err != nil {
		panic(err)
	}
	pic3, err := loadPicture("assets/textures/spritealien3.png")
	if err != nil {
		panic(err)
	}

	// pic4, err := loadPicture("spritealien4.png")
	// if err != nil {
	// 	panic(err)
	// }

	spriteAlien1 := pixel.NewSprite(pic1, pic1.Bounds())
	spriteAlien2 := pixel.NewSprite(pic2, pic2.Bounds())
	spriteAlien3 := pixel.NewSprite(pic3, pic3.Bounds())
	//spriteAlien4 := pixel.NewSprite(pic4, pic4.Bounds())

	mat := pixel.IM
	mat = mat.Moved(window.Bounds().Center())
	mat = mat.ScaledXY(window.Bounds().Center(), pixel.V(0.05, 0.05))
	for x := 0; x < 5; x++ {
		for i := 0; i < 10; i++ {
			if x == 0 {
				//Dibujar alien 3

			} else if x == 1 || x == 2 {
				//Dibujar alien 1

			} else {
				//Dibujar alien 2

			}
		}
	}

	spriteAlien1.Draw(window, mat)
	spriteAlien2.Draw(window, mat)
	spriteAlien3.Draw(window, mat)

	//spriteAlien4.Draw(window, mat)
}

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

		createEnemies(win)

		player.Update(direction, action, dt)
		player.Draw(win)
		direction = spacegame.Idle
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
