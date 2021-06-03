package libs

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func NewCreateFortress(win *pixelgl.Window) {
	//Setup 3x4x4
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				imd := imdraw.New(nil)
				imd.Color = colornames.White
				imd.Push(pixel.V(float64(65+30+k*18+j*173), float64(80+i*18)))
				imd.Push(pixel.V(float64(83+30+k*18+j*173), float64(98+i*18)))
				imd.Rectangle(3)
				imd.Draw(win)
			}
		}
	}

}
