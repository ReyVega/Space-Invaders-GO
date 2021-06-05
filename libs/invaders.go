package libs

import (
	"log"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Invader struct {
	pos        *pixel.Vec
	vel        float64
	sprite     *pixel.Sprite
	laser      *EnemyLaser
	laserDelay int

	lasers map[string]*EnemyLaser
}

const (
	rechargeTimeEnemy = 70
)

func NewCreateEnemies(win *pixelgl.Window, numAliens int, world *World) ([]Invader, error) {
	// Arreglo de invaders
	var invaders []Invader

	//estandar 50 num de aliens
	contador := 0
	for x := 0; x < 5; x++ {
		for i := 0; i < 10; i++ {
			pos := pixel.V(win.Bounds().Center().X+float64(i)*50-240.0, win.Bounds().Center().Y+float64(x)*40+100)
			if x == 0 {
				//Dibujar alien 3
				invader, err := NewInvader("assets/textures/spritealien1.png", pos, world)
				if err != nil {
					log.Fatal(err)
				}
				invaders = append(invaders, *invader)
			} else if x == 1 || x == 2 {
				//Dibujar alien 1
				invader, err := NewInvader("assets/textures/spritealien2.png", pos, world)
				if err != nil {
					log.Fatal(err)
				}
				invaders = append(invaders, *invader)
			} else {
				//Dibujar alien 2
				invader, err := NewInvader("assets/textures/spritealien3.png", pos, world)
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

func NewInvader(path string, pos pixel.Vec, world *World) (*Invader, error) {
	pic, err := NewloadPicture(path)
	if err != nil {
		return nil, err
	}
	spr := pixel.NewSprite(pic, pic.Bounds())
	initialPos := pos

	// Initialize the laser for the player
	l, err := NewEnemyLaser("assets/textures/mars.png", 200.0, world)
	if err != nil {
		return nil, err
	}

	invader := &Invader{
		pos:        &initialPos,
		vel:        100.00,
		sprite:     spr,
		laser:      l,
		laserDelay: 30,
		lasers:     make(map[string]*EnemyLaser),
	}
	return invader, nil
}

func (inv Invader) Draw(t pixel.Target) {
	mat := pixel.IM
	mat = mat.Moved(*inv.pos)
	mat = mat.Scaled(*inv.pos, 0.075)
	inv.sprite.Draw(t, mat)

	for _, l := range inv.lasers {
		l.Draw(t)
	}
}

func (inv *Invader) Update(movementX bool, movementY bool, dt float64) {
	go move(inv, dt, movementX, movementY)

	if rand.Intn((500-0)+0) > 480 {
		inv.shoot(dt)
	}

	for k, l := range inv.lasers {
		l.Update()

		// remove unused lasers
		if !l.isVisible {
			delete(inv.lasers, k)
		}
	}
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

func (inv *Invader) shoot(dt float64) {

	if inv.laserDelay >= 0 {
		inv.laserDelay--
	}

	if inv.laserDelay <= 0 {
		l := inv.laser.NewEnemLaser(*inv.pos)
		l.vel *= dt
		inv.lasers[NewULID()] = l
		inv.laserDelay = rechargeTimeEnemy
	}
}

func (inv *Invader) CheckFortressInvaders(coordenadasFortress [32]pixel.Vec, deadFortress [16]int) (coordenadasFortalezas [32]pixel.Vec, deadFortalezas [16]int) {
	for k, l := range inv.lasers {
		l.Update()
		for i := 0; i < 16; i++ {
			if l.pos.X >= coordenadasFortress[i*2].X && l.pos.X <= coordenadasFortress[i*2+1].X && l.pos.Y >= coordenadasFortress[i*2].Y && l.pos.Y <= coordenadasFortress[i*2+1].Y {

				delete(inv.lasers, k)
				deadFortress[i] = 1
				coordenadasFortress[i*2] = pixel.V(0, 0)
				coordenadasFortress[i*2+1] = pixel.V(0, 0)
			}
		}

	}
	return coordenadasFortress, deadFortress
}
