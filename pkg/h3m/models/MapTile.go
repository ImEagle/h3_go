package models

type MapTile struct {
	Surface         byte
	SurfacePicture  byte
	RiverType       byte
	RiverProperties byte
	RoadType        byte
	RoadProperties  byte
	Mirroring       byte
}

func (mt *MapTile) TerrainFlipX() bool {
	return mt.Mirroring&0x01 != 0
}

func (mt *MapTile) TerrainFlipY() bool {
	return mt.Mirroring&0x02 != 0
}
