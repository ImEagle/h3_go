package game

import (
	"github.com/ImEagle/h3_go/pkg/h3m"
	"github.com/ImEagle/h3_go/pkg/h3m/models"
	"github.com/ImEagle/h3_go/pkg/lod"
	"github.com/hajimehoshi/ebiten/v2"
)

type ObjectsRenderer struct {
	mapData       *h3m.H3m
	spriteManager SpriteManager
}

func NewObjectsRenderer(mapData *h3m.H3m, spriteManager *lod.Reader) *ObjectsRenderer {
	memSprites := NewInMemorySpriteManager(spriteManager)

	return &ObjectsRenderer{
		mapData:       mapData,
		spriteManager: memSprites,
	}
}

func (r *ObjectsRenderer) Draw(screen *ebiten.Image) {
	tileSize := 32

	for _, objDef := range r.mapData.ObjectsDefinition {
		if !isRenderable(objDef.MapObjectPosition) {
			continue
		}

		dio := &ebiten.DrawImageOptions{}

		dioX := float64(int(objDef.MapObjectPosition.X) * tileSize)
		dioY := float64(int(objDef.MapObjectPosition.Y) * tileSize)

		dio.GeoM.Translate(dioX, dioY)

		images, err := r.spriteManager.Get(objDef.SpriteName)
		if err != nil {
			continue
		}

		img := images[0]
		eImg := ebiten.NewImageFromImage(img.Image)

		screen.DrawImage(eImg, dio)

	}
}

func isRenderable(position *models.MapObjectPosition) bool {
	if position == nil {
		return false
	}

	if position.X < 0 || position.Y < 0 {
		return false
	}

	return true
}
