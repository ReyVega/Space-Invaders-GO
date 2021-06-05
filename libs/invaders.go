package libs

import (
	"log"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Invader struct {
	pos    *pixel.Vec
	vel    float64
	sprite *pixel.Sprite
}

func NewCreateEnemies(win *pixelgl.Window, numAliens int) ([]Invader, error) {
	// Arreglo de invaders
	var invaders []Invader

	//estandar 50 num de aliens
	contador := 0
	for x := 0; x < 5; x++ {
		for i := 0; i < 10; i++ {
			pos := pixel.V(win.Bounds().Center().X+float64(i)*50-240.0, win.Bounds().Center().Y+float64(x)*40+100)
			if x == 0 {
				//Dibujar alien 3
				invader, err := NewInvader("assets/textures/spritealien1.png", pos)
				if err != nil {
					log.Fatal(err)
				}
				invaders = append(invaders, *invader)
			} else if x == 1 || x == 2 {
				//Dibujar alien 1
				invader, err := NewInvader("assets/textures/spritealien2.png", pos)
				if err != nil {
					log.Fatal(err)
				}
				invaders = append(invaders, *invader)
			} else {
				//Dibujar alien 2
				invader, err := NewInvader("assets/textures/spritealien3.png", pos)
				if err != nil {
					log.Fatal(err)
				}
				invaders = append(invaders, *invader)
			}
			contador++
			if contador >= numAliens {
				break
			}
		}
	}
	return invaders, nil
}

func NewInvader(path string, pos pixel.Vec) (*Invader, error) {
	pic, err := NewloadPicture(path)
	if err != nil {
		return nil, err
	}
	spr := pixel.NewSprite(pic, pic.Bounds())
	initialPos := pos

	invader := &Invader{
		pos:    &initialPos,
		vel:    50.00,
		sprite: spr,
	}
	return invader, nil
}

func (inv Invader) Draw(t pixel.Target) {
	mat := pixel.IM
	mat = mat.Moved(*inv.pos)
	mat = mat.Scaled(*inv.pos, 0.075)
	inv.sprite.Draw(t, mat)
}

func (inv *Invader) Update(movementX bool, movementY bool, dt float64) {
	go move(inv, dt, movementX, movementY)
}

func move(invader *Invader, dt float64, movementX bool, movementY bool) {
	if movementX {
		invader.pos.X = invader.pos.X - (invader.vel * dt)
	} else {
		invader.pos.X = invader.pos.X + (invader.vel * dt)
	}

	if movementY {
		invader.pos.Y = invader.pos.Y - 10
	}
}
