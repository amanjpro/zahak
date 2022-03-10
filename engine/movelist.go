package engine

type MoveList struct {
	Moves    []Move
	Scores   []int32
	IsScored bool
	Size     int
	Next     int
}

func NewMoveList(capacity int) MoveList {
	return MoveList{
		make([]Move, capacity),
		make([]int32, capacity),
		false,
		0,
		0,
	}
}

func (ml *MoveList) Add(m Move) {
	ml.Moves[ml.Size] = m
	ml.Size += 1
}

func (ml *MoveList) AddFour(m1, m2, m3, m4 Move) {
	ml.Moves[ml.Size] = m1
	ml.Moves[ml.Size+1] = m2
	ml.Moves[ml.Size+2] = m3
	ml.Moves[ml.Size+3] = m4
	ml.Size += 4
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
