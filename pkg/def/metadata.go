package def

import "math"

type Frame struct {
	X uint32
	Y uint32
	W uint32
	H uint32
}

var battleCreature = map[uint32]string{
	0:  "moving",
	1:  "mouseon",
	2:  "holding",
	3:  "hitted",
	4:  "defence",
	5:  "death",
	6:  "death_ranged",
	7:  "turn_l",
	8:  "turn_r",
	9:  "turn_l2",
	10: "turn_r2",
	11: "attack_up",
	12: "attack_front",
	13: "attack_down",
	14: "shoot_up",
	15: "shoot_front",
	16: "shoot_down",
	17: "cast_up",
	18: "cast_front",
	19: "cast_down",
	20: "move_start",
	21: "move_end",
	22: "dead",
	23: "dead_ranged",
}

type Metadata struct {
	Frames map[string][]Frame
}

func NewMetadata(r *Reader) *Metadata {
	frameWidth := r.header.Width
	frameHeight := r.header.Height
	countInRow := 8 // TODO: Is it important? Can be changed to X?
	md := &Metadata{
		Frames: make(map[string][]Frame),
	}

	index := 0

	for i := uint32(0); i < r.header.BlocksCount; i++ {
		block := r.blocks[i]
		animationType := battleCreature[block.Id]

		for j := uint32(0); j < block.Count; j++ {
			x := uint32(index % countInRow)
			y := uint32(math.Floor(float64(index / countInRow)))

			f := Frame{
				X: x * frameWidth,
				Y: y * frameHeight,
				W: frameWidth,
				H: frameHeight,
			}

			md.Frames[animationType] = append(md.Frames[animationType], f)
			index++
		}
	}

	return md
}
