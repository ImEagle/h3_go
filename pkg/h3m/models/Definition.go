package models

type MapObjectDefinition struct {
	*MapObjectData
	SpriteName string
}

type MapObjectData struct {
	PassableSquares [6]byte
	ActiveSquare    [6]byte
	Landscape       [2]byte
	LandscapeGroup  [2]byte
	Class           [4]byte
	Number          [4]byte
	Group           byte
	OverOrBelow     byte
	Unknown         [16]byte
}

type MapObjectPosition struct {
	X uint8
	Y uint8
	Z uint8
}
