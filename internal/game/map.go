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

			// Note: Debugging invalid images
			tileImg := r.getLandImage(tileInfo.Surface, tileInfo.SurfacePicture)

			if tileImg == nil {
				continue
			}

			subImage := tileImg.SubImage(image.Rect(x, y, w, h)).(*ebiten.Image)

			dioX := float64(int(mapX) * tileSize)
			dioY := float64(int(mapY) * tileSize)

			if tileInfo.TerrainFlipX() && tileInfo.TerrainFlipY() {
				dio.GeoM.Scale(-1, -1)
				dio.GeoM.Translate(float64(tileSize), float64(tileSize))
			} else if tileInfo.TerrainFlipX() && !tileInfo.TerrainFlipY() {
				dio.GeoM.Scale(-1, 1)
				dio.GeoM.Translate(float64(tileSize), 0)
			} else if !tileInfo.TerrainFlipX() && tileInfo.TerrainFlipY() {
				dio.GeoM.Scale(1, -1)
				dio.GeoM.Translate(0, float64(tileSize))
			} else if !tileInfo.TerrainFlipX() && !tileInfo.TerrainFlipY() {
				dio.GeoM.Scale(1, 1)
			}

			dio.GeoM.Translate(dioX, dioY)

			screen.DrawImage(subImage, dio)

		}
	}

}

func (r *Renderer) getLandImage(landType byte, pictureIndex byte) *ebiten.Image {
	spriteMapper := map[byte]string{
		0: "dirttl.def",
		1: "sandtl.def",
		2: "grastl.def",
		3: "snowtl.def",
		4: "swmptl.def",
		5: "rougtl.def",
		6: "subbtl.def",
		7: "lavatl.def",
		8: "watrtl.def",
		9: "rocktl.dev",
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

	if pictureIndex >= byte(len(images)) {
		// TODO: Fix this
		lastIdx := len(images) - 1
		return ebiten.NewImageFromImage(images[lastIdx].Image)
	}

	return ebiten.NewImageFromImage(images[pictureIndex].Image)
}
