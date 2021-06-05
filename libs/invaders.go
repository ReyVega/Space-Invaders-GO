package libs

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func NewCreateEnemies(window *pixelgl.Window, numAliens int) {
	pic1, err := NewloadPicture("assets/textures/spritealien1.png")
	if err != nil {
		panic(err)
	}
	pic2, err := NewloadPicture("assets/textures/spritealien2.png")
	if err != nil {
		panic(err)
	}
	pic3, err := NewloadPicture("assets/textures/spritealien3.png")
	if err != nil {
		panic(err)
	}

	// pic4, err := loadPicture("/assets/textures/spritealien4.png")
	// if err != nil {
	// 	panic(err)
	// }
	spriteAlien1 := pixel.NewSprite(pic1, pic1.Bounds())
	spriteAlien2 := pixel.NewSprite(pic2, pic2.Bounds())
	spriteAlien3 := pixel.NewSprite(pic3, pic3.Bounds())
	//spriteAlien4 := pixel.NewSprite(pic4, pic4.Bounds())
	//estandar 50 num de aliens
	contador := 0
	for x := 0; x < 5; x++ {
		for i := 0; i < 10; i++ {
			mat := pixel.IM
			mat = mat.Moved(pixel.V(window.Bounds().Center().X+float64(i)*50-240.0, window.Bounds().Center().Y+float64(x)*40+100))
			mat = mat.Scaled(pixel.V(window.Bounds().Center().X+float64(i)*50-240.0, window.Bounds().Center().Y+float64(x)*40+100), 0.075)
			if x == 0 {
				//Dibujar alien 3
				spriteAlien3.Draw(window, mat)
			} else if x == 1 || x == 2 {
				//Dibujar alien 1
				spriteAlien1.Draw(window, mat)
			} else {
				//Dibujar alien 2
				spriteAlien2.Draw(window, mat)
			}
			contador++
			if contador >= numAliens {
				break
			}
		}
	}
	//spriteAlien4.Draw(window, mat)
}
