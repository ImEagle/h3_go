package game

import (
	"github.com/ImEagle/h3_go/pkg/def"
	"github.com/ImEagle/h3_go/pkg/lod"
)

type SpriteManager interface {
	Get(name string) ([]def.ImageDetails, error)
}

func NewInMemorySpriteManager(sprites *lod.Reader) *InMemorySpriteManager {
	return &InMemorySpriteManager{
		sprites: sprites,
		cache:   make(map[string][]def.ImageDetails),
	}
}

type InMemorySpriteManager struct {
	sprites *lod.Reader
	cache   map[string][]def.ImageDetails
}

func (i *InMemorySpriteManager) Get(name string) ([]def.ImageDetails, error) {
	if images, ok := i.cache[name]; ok {
		return images, nil
	}

	defPayload, err := i.sprites.GetFile(name)
	if err != nil {
		return nil, err
	}

	defReader := def.NewReader()
	images, err := defReader.Load(defPayload)
	if err != nil {
		return nil, err
	}

	i.cache[name] = images
	return images, nil
}
