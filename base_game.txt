package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sync/atomic"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type spaceShip struct {
	kind  string
	name  string
	x     float64
	y     float64
	alive bool
	vx    float64
	vy    float64
	in    chan msg
}

type motherShip struct {
	register         map[string]spaceShip
	registerChannel  chan msg
	generalChannel   chan msg
	childrenChannels map[string](*chan msg)
}

type msg struct {
	cmd string
	val string
	p   spaceShip
}

var gos uint64

var window *pixelgl.Window

var ms = motherShip{register: make(map[string]spaceShip), childrenChannels: make(map[string](*chan msg)), generalChannel: make(chan msg, 20000), registerChannel: make(chan msg, 20000)}

var channelSpeed int
var shipSpeed int
var numAliens int
var score int
var live int
var shoot bool

func callGo(f func()) {
	atomic.AddUint64(&gos, 1)
	go f()
}

func (p *spaceShip) conductor() { // for spaceShip (player)
	callGo(p.actions)
	for range time.NewTicker(time.Duration(channelSpeed) * time.Millisecond).C {
		p.in <- msg{cmd: "Move"}
	}
}

func (m *motherShip) conductor() { // for motherShip (enemy)
	populate()
	callGo(ms.actions)
	var a = 0
	for range time.NewTicker(time.Duration(shipSpeed) * time.Millisecond).C {
		a = a + 1
		if a%2 == 0 {
			ms.generalChannel <- msg{cmd: "Display"}
		}
		if a%2 == 0 {
			ms.generalChannel <- msg{cmd: "CheckCollisions"}
		}
		if a%2 == 0 {
			ms.generalChannel <- msg{cmd: "CheckKeys"}
		}
	}
}

func (p *spaceShip) actions() { // for spaceShip (player)
	for {
		m := <-p.in
		if m.cmd == "Die" {
			p.alive = false
			ms.registerChannel <- msg{cmd: "Remove", p: *p}
			if p.kind == "Gun" {
				fmt.Println("Damn and blast")
				<-make(chan bool)
			}
		} else if m.cmd == "Left" {
			if p.x > 12 {
				p.vx = -80
			} else {
				p.vx = 0
			}
		} else if m.cmd == "Stop" {
			p.vx = 0
		} else if m.cmd == "Right" {
			if p.x < 500 {
				p.vx = 80
			} else {
				p.vx = 0
			}
		} else if m.cmd == "Shoot" {
			if p.x <= 12 {
				p.vx = 0
			}
			if p.x >= 500 {
				p.vx = 0
			}
			if shoot {
				ms.generalChannel <- msg{cmd: "Add", val: "Bullet", p: spaceShip{x: p.x, y: p.y, vx: 0, vy: 100}}
				shoot = false
			}
		}
		if m.cmd == "Move" {
			var xPixPerBeat = p.vx / 1000 * float64(shipSpeed)
			var yPixPerBeat = p.vy / 1000 * float64(shipSpeed)
			p.x = p.x + xPixPerBeat
			if p.alive && p.kind == "Alien" {
				if rand.Intn(950) < 1 {
					ms.generalChannel <- msg{cmd: "Add", val: "Bomb", p: spaceShip{x: p.x, y: p.y, vx: 0, vy: -100}}
				}
				if p.x > 500 || p.x < 10 {
					p.vx = -p.vx
					p.y = p.y - 10
				}
			} else {
				p.y = p.y + yPixPerBeat
			}
			if p.kind == "Bullet" && p.y > 600 {
				p.alive = false
				ms.registerChannel <- msg{cmd: "Remove", p: *p}
			} else if p.kind == "Bomb" && p.y < 0 {
				p.alive = false
				ms.registerChannel <- msg{cmd: "Remove", p: *p}
			}
			if p.alive {
				ms.registerChannel <- msg{cmd: "Set", p: *p}
			}
		}
	}
}

func (m *motherShip) actions() { // for motherShip (enemy o barrier)
	for !window.Closed() {
		select {
		case message := <-m.registerChannel:
			if message.cmd == "Set" {
				m.register[message.p.name] = spaceShip{name: message.p.name, kind: message.p.kind, x: message.p.x, y: message.p.y, alive: message.p.alive}
			} else if message.cmd == "Remove" {
				delete(m.register, message.p.name)
			}
		case message := <-m.generalChannel:
			if message.cmd == "Add" {
				name := message.val
				if name != "Gun" {
					name = fmt.Sprintf(message.val, time.Now())
				}
				p := spaceShip{name: name, kind: message.val, x: message.p.x, y: message.p.y, alive: true, vx: message.p.vx, vy: message.p.vy, in: make(chan msg)}
				m.childrenChannels[p.name] = &p.in
				callGo(p.conductor)
			}
			if message.cmd == "CheckCollisions" {
				for _, s1 := range m.register {
					if s1.alive && s1.kind == "Bullet" {
						for _, s2 := range m.register {
							if s2.alive && (s2.kind == "Alien" || s2.kind == "Fortress") {
								if (s2.x-s1.x)*(s2.x-s1.x)+(s2.y-s1.y)*(s2.y-s1.y) < 40 {
									*m.childrenChannels[s2.name] <- msg{cmd: "Die"}
									score += 10
									*m.childrenChannels[s1.name] <- msg{cmd: "Die"}
								}
							}
						}
					}
					if s1.alive && s1.kind == "Bomb" {
						for _, s2 := range m.register {
							if s2.alive && s2.kind == "Gun" {
								if (s2.x-s1.x)*(s2.x-s1.x)+(s2.y-s1.y)*(s2.y-s1.y) < 40 {
									*m.childrenChannels[s2.name] <- msg{cmd: "One live less"}
									live--
									if live == 0 {
										*m.childrenChannels[s2.name] <- msg{cmd: "Die"}
									}
									*m.childrenChannels[s1.name] <- msg{cmd: "Die"}
								}
							}
							if s2.alive && s2.kind == "Fortress" {
								if (s2.x-s1.x)*(s2.x-s1.x)+(s2.y-s1.y)*(s2.y-s1.y) < 40 {
									*m.childrenChannels[s2.name] <- msg{cmd: "Die"}
									*m.childrenChannels[s1.name] <- msg{cmd: "Die"}
								}
							}
						}
					}
					if s1.alive && s1.kind == "Alien" {
						for _, s2 := range m.register {
							if s2.alive && s2.kind == "Gun" {
								if (s2.x-s1.x)*(s2.x-s1.x)+(s2.y-s1.y)*(s2.y-s1.y) < 40 {
									*m.childrenChannels[s2.name] <- msg{cmd: "One live less"}
									live--
									if live == 0 {
										*m.childrenChannels[s2.name] <- msg{cmd: "Die"}
									}
									*m.childrenChannels[s1.name] <- msg{cmd: "Die"}
								}
							}
							if s2.alive && s2.kind == "Fortress" {
								if (s2.x-s1.x)*(s2.x-s1.x)+(s2.y-s1.y)*(s2.y-s1.y) < 40 {
									*m.childrenChannels[s2.name] <- msg{cmd: "Die"}
									*m.childrenChannels[s1.name] <- msg{cmd: "Die"}
								}
							}
						}
					}
				}
			}
			if message.cmd == "CheckKeys" {
				if window != nil {
					if m.childrenChannels["Gun"] != nil {
						if window.Pressed(pixelgl.KeyLeft) {
							*m.childrenChannels["Gun"] <- msg{cmd: "Left"}
						} else if window.Pressed(pixelgl.KeyRight) {
							*m.childrenChannels["Gun"] <- msg{cmd: "Right"}
						} else if window.Pressed(pixelgl.KeySpace) {
							*m.childrenChannels["Gun"] <- msg{cmd: "Shoot"}
						} else {
							*m.childrenChannels["Gun"] <- msg{cmd: "Stop"}
						}
					}
				}
			}
			if message.cmd == "Display" {
				if window != nil {
					var imd = imdraw.New(nil)
					imd.Clear()
					for _, s := range m.register {
						if s.alive && s.kind == "Gun" {
							imd.Color = colornames.Green
							imd.Push(pixel.V(float64(s.x), float64(s.y)))
							imd.Circle(6, 3)
						} else if s.alive && s.kind == "Alien" {
							imd.Color = colornames.Red
							imd.Push(pixel.V(float64(s.x), float64(s.y)))
							imd.Circle(4, 2)
						} else if s.alive && s.kind == "Fortress" {
							imd.Color = colornames.White
							imd.Push(pixel.V(float64(s.x), float64(s.y)))
							imd.Circle(5, 4)
						} else if s.alive && s.kind == "Bullet" {
							imd.Color = colornames.Orange
							imd.Push(pixel.V(float64(s.x), float64(s.y)))
							imd.Circle(2, 2)
						} else if s.alive && s.kind == "Bomb" {
							imd.Color = colornames.Blue
							imd.Push(pixel.V(float64(s.x), float64(s.y)))
							imd.Circle(3, 2)
						}
					}
					window.Clear(colornames.Black)
					imd.Draw(window)

					basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
					tvScore := text.New(pixel.V(15, 490), basicAtlas)
					tvLives := text.New(pixel.V(410, 490), basicAtlas)

					fmt.Fprintln(tvScore, "Score: ", score)
					fmt.Fprintln(tvLives, "Lives: ", live)
					tvScore.Draw(window, pixel.IM.Scaled(tvScore.Orig, 1.25))
					tvLives.Draw(window, pixel.IM.Scaled(tvLives.Orig, 1.25))
					window.Update()
				}
			}
		}
	}
	window.Destroy()
	os.Exit(1)
}

func populate() {
	var spacingPx = 50
	var con = 0
	for r := 0; r < 300/spacingPx; r++ {
		for i := 0; i < 400/spacingPx; i++ {
			if con < numAliens {
				ms.generalChannel <- msg{cmd: "Add", val: "Alien", p: spaceShip{x: float64(100 + spacingPx*i), y: float64(200 + spacingPx*r), vx: float64(-100 + rand.Intn(1)), vy: 0}}
			}
			con++
		}
	}
	for f := 0; f < 3; f++ {
		for c := 0; c < 3; c++ {
			for r := 0; r < 3; r++ {
				ms.generalChannel <- msg{cmd: "Add", val: "Fortress", p: spaceShip{x: float64(100 + f*100 + c*10), y: float64(60 + r*10), vx: 0, vy: 0}}
			}
		}
	}
	ms.generalChannel <- msg{cmd: "Add", val: "Gun", p: spaceShip{x: 100, y: 10, vx: 0, vy: 0}}
}

func start() {
	cfg := pixelgl.WindowConfig{
		Title:  "Space Invaders",
		Bounds: pixel.R(0, 0, float64(512), float64(512)),
		VSync:  true,
	}

	window, _ = pixelgl.NewWindow(cfg)
	window.SetPos(window.GetPos().Add(pixel.V(0, 1)))

	go func() {
		ticker := time.NewTicker(time.Second / 2)
		for _ = range ticker.C {
			shoot = true
		}
	}()

	callGo(ms.conductor)
	<-make(chan bool)
}

func main() {
	channelSpeed = 8
	shipSpeed = 17
	live = 10
	flag.IntVar(&numAliens, "aliens", 10, "Number of aliens")

	flag.Parse()
	pixelgl.Run(start)
}
