package main

import (
	"github.com/ImEagle/h3_go/internal/game"
	"github.com/ImEagle/h3_go/pkg/h3m"
	"github.com/ImEagle/h3_go/pkg/lod"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	mapName := "/Users/eagle/projects/h3_go/pkg/h3m/test_maps/test_1.h3m"
	lodName := "/Users/eagle/Downloads/h3/h3ab_bmp.lod"

	h3Map, err := h3m.Load(mapName)
	if err != nil {
		panic(err)
	}

	h3Sprites := lod.NewReader(lodName)
	h3Sprites.LoadMetadata()

	gameRenderer := game.NewRendered(h3Map, h3Sprites)
	fullGame := NewGame(gameRenderer)

	ebiten.SetWindowSize(800, 600)
	ebiten.RunGame(fullGame)
}
