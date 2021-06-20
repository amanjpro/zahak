package engine

type MoveList struct {
	Moves  []Move
	Scores []int32
	Size   int
	Next   int
}

func NewMoveList(capacity int) *MoveList {
	return &MoveList{
		make([]Move, capacity),
		make([]int32, capacity),
		0,
		0,
	}
}

// func (ml *MoveList) Add(m Move) {
// 	ml.Moves[ml.Size] = m
// 	ml.Size += 1
// }

func (ml *MoveList) Add(ms ...Move) {
	for _, m := range ms {
		ml.Moves[ml.Size] = m
		ml.Size += 1
	}
}

func (ml *MoveList) IsEmpty() bool {
	return ml.Size == 0
}

func (ml *MoveList) SwapWith(best int) {
	ml.Swap(0, best)
}

func (ml *MoveList) Swap(first int, second int) {
	ml.Moves[first], ml.Moves[second] = ml.Moves[second], ml.Moves[first]
	ml.Scores[first], ml.Scores[second] = ml.Scores[second], ml.Scores[first]
}

func (ml *MoveList) IncNext() {
	ml.Next += 1
}
