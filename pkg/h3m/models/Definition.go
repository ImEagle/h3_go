package models

type MapObjectDefinition struct {
	*MapObjectData
	SpriteName string
	*MapObjectPosition
}

type MapObjectData struct {
	PassableSquares [6]byte
	ActiveSquare    [6]byte
	Landscape       [2]byte
	LandscapeGroup  [2]byte
	Class           uint32
	Number          uint32
	Group           byte
	OverOrBelow     byte
	Unknown         [16]byte
}

type MapObjectPosition struct {
	X              uint8
	Y              uint8
	Z              uint8
	ObjectDefIndex uint32
}
