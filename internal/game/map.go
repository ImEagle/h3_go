package game

import (
	"github.com/ImEagle/h3_go/pkg/h3m"
	"github.com/ImEagle/h3_go/pkg/lod"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	_ "image/png"
)

func NewRendered(mapData *h3m.H3m, spriteManager *lod.Reader) *Renderer {

	memSprites := NewInMemorySpriteManager(spriteManager)

	return &Renderer{
		mapData:       mapData,
		spriteManager: memSprites,
	}
}

type Renderer struct {
	mapData       *h3m.H3m
	spriteManager SpriteManager
}

func (r *Renderer) Draw(screen *ebiten.Image) {
	tileSize := 32

	var mapSize, mapY, mapX uint32
	mapSize = r.mapData.MapSize

	for mapY = 0; mapY < mapSize; mapY++ {
		for mapX = 0; mapX < mapSize; mapX++ {
			dio := &ebiten.DrawImageOptions{}

			tileIdx := (mapY * mapSize) + mapX
			tileInfo := r.mapData.LandMap[tileIdx] // #TODO: Fix for rendering land and underground

			x := 0
			y := 0
			w := x + tileSize
			h := y + tileSize

			tileImg := r.getLandImage(tileInfo.Surface, tileInfo.SurfacePicture)

			if tileImg == nil {
				continue
			}

			subImage := tileImg.SubImage(image.Rect(x, y, w, h)).(*ebiten.Image)

			dioX := float64(int(mapX) * tileSize)
			dioY := float64(int(mapY) * tileSize)
			dio.GeoM.Translate(dioX, dioY)

			screen.DrawImage(subImage, dio)

		}
	}

}

func (r *Renderer) getLandImage(landType byte, pictureIndex byte) *ebiten.Image {
	spriteMapper := map[byte]string{
		0: "dirttl.def",
		1: "sandtl.def",
		2: "grass",
		3: "snow",
		4: "swamp",
		5: "rough",
		6: "subterranean",
		7: "lava",
		8: "water",
		9: "rock",
	}

	spriteName, ok := spriteMapper[landType]
	if !ok {
		return nil
	}

	// TODO: Add cache for the files
	images, err := r.spriteManager.Get(spriteName)
	if err != nil {
		return nil
	}

	return ebiten.NewImageFromImage(images[pictureIndex].Image)
}
