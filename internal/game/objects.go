package game

import (
	"github.com/ImEagle/h3_go/pkg/h3m"
	"github.com/ImEagle/h3_go/pkg/lod"
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

func (r *ObjectsRenderer) Draw() {
	// NOP
}