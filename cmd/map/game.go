package main

import (
	"github.com/ImEagle/h3_go/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func NewGame(mapRenderer *game.Renderer, objRenderer *game.ObjectsRenderer) *Game {
	return &Game{
		mapRenderer: mapRenderer,
		objRenderer: objRenderer,
	}
}

type Game struct {
	mapRenderer *game.Renderer
	objRenderer *game.ObjectsRenderer
}

func (g Game) Update() error {
	return nil
}

func (g Game) Draw(screen *ebiten.Image) {
	g.mapRenderer.Draw(screen)
	g.objRenderer.Draw(screen)
}

func (g Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
