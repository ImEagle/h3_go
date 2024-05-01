package main

import (
	"github.com/ImEagle/h3_go/internal/game"
	"github.com/ImEagle/h3_go/pkg/h3m"
	"github.com/ImEagle/h3_go/pkg/lod"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	mapName := "/Users/eagle/projects/h3_go/pkg/h3m/generated.h3m"
	lodName := "/Users/eagle/Downloads/h3/h3sprite.lod"

	h3Map, err := h3m.Load(mapName)
	if err != nil {
		panic(err)
	}

	h3Sprites := lod.NewReader(lodName)
	h3Sprites.LoadMetadata()

	gameRenderer := game.NewRendered(h3Map, h3Sprites)
	objectsRenderer := game.NewObjectsRenderer(h3Map, h3Sprites)
	fullGame := NewGame(gameRenderer, objectsRenderer)

	ebiten.SetWindowSize(1500, 1000)
	ebiten.RunGame(fullGame)

	// NOP
	// NOP
	// NOP
	// NOP
	// NOP
	// NOP
}
