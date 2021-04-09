package search

import (
	. "github.com/amanjpro/zahak/engine"
)

type MovePicker struct {
	position       *Position
	engine         *Engine
	hashmove       Move
	quietMoves     []Move
	quietScores    []int32
	captureMoves   []Move
	captureScores  []int32
	moveOrder      int8
	nextQuiet      int
	nextCapture    int
	canUseHashMove bool
	isQuiescence   bool
}

func NewMovePicker(p *Position, e *Engine, moveOrder int8, hashmove Move, isQuiescence bool) *MovePicker {
	nextCapture := 0
	nextQuiet := 0
	if hashmove != EmptyMove {
		if hashmove.IsCapture() {
			nextCapture = 1
		} else {
			nextQuiet = 1
		}
	}
	mp := &MovePicker{
		position:       p,
		engine:         e,
		hashmove:       hashmove,
		quietMoves:     nil,
		quietScores:    nil,
		captureMoves:   nil,
		captureScores:  nil,
		moveOrder:      moveOrder,
		nextQuiet:      nextQuiet,
		nextCapture:    nextCapture,
		canUseHashMove: hashmove != EmptyMove,
		isQuiescence:   isQuiescence,
	}
	return mp
}

func (mp *MovePicker) Length() int {
	mp.generateQuietMoves()
	mp.generateCaptureMoves()
	return len(mp.captureMoves) + len(mp.quietMoves)
}

func (mp *MovePicker) generateQuietMoves() {
	if mp.isQuiescence || mp.quietMoves != nil {
		return
	}
	mp.quietMoves = mp.position.GetQuietMoves()
	mp.scoreQuietMoves()
}

func (mp *MovePicker) generateCaptureMoves() {
	if mp.captureMoves != nil || mp.quietMoves != nil {
		return
	}
	mp.captureMoves = mp.position.GetCaptureMoves()
	mp.scoreCaptureMoves()
}

func (mp *MovePicker) HasNoPVMove() bool {
	return mp.hashmove == EmptyMove
}

func (mp *MovePicker) UpgradeToPvMove(pvMove Move) {
	if pvMove == EmptyMove || mp.captureMoves != nil || mp.quietMoves != nil {
		return
	}
	mp.hashmove = pvMove
	mp.canUseHashMove = true
}

func (mp *MovePicker) scoreCaptureMoves() {
	position := mp.position
	board := position.Board
	mp.captureScores = make([]int32, len(mp.captureMoves))

	for i := 0; i < len(mp.captureMoves); i++ {
		move := mp.captureMoves[i]

		if move.EqualTo(mp.hashmove) {
			mp.captureScores[i] = 900_000_000
			mp.captureScores[0], mp.captureScores[i] = mp.captureScores[i], mp.captureScores[0]
			mp.captureMoves[0], mp.captureMoves[i] = mp.captureMoves[i], mp.captureMoves[0]
			mp.nextCapture = 1
			continue
		}

		source := move.Source()
		dest := move.Destination()
		piece := move.MovingPiece()
		//
		// capture ordering
		if move.IsCapture() {
			capPiece := move.CapturedPiece()
			cweight := capPiece.Weight()
			pweight := piece.Weight()
			if cweight <= pweight {
				// SEE for ordering
				gain := int32(board.StaticExchangeEval(dest, capPiece, source, piece))
				if gain < 0 {
					mp.captureScores[i] = -90_000_000 + gain
				} else if gain == 0 {
					mp.captureScores[i] = 100_000_000 + int32(cweight-pweight)
				} else {
					mp.captureScores[i] = 100_100_000 + gain
				}
			} else {
				mp.captureScores[i] = 100_100_000 + int32(cweight-pweight)
			}
			continue
		}

		promoType := move.PromoType()
		if promoType != NoType {
			p := GetPiece(promoType, White)
			mp.captureScores[i] = 50_000_000 + int32(p.Weight())
			continue
		}

	}
}

func (mp *MovePicker) scoreQuietMoves() {
	engine := mp.engine
	moveOrder := mp.moveOrder
	mp.quietScores = make([]int32, len(mp.quietMoves))

	for i := 0; i < len(mp.quietMoves); i++ {
		move := mp.quietMoves[i]

		if move.EqualTo(mp.hashmove) {
			mp.quietScores[i] = 900_000_000
			mp.quietScores[0], mp.quietScores[i] = mp.quietScores[i], mp.quietScores[0]
			mp.quietMoves[0], mp.quietMoves[i] = mp.quietMoves[i], mp.quietMoves[0]
			mp.nextQuiet = 1
			continue
		}

		dest := move.Destination()
		piece := move.MovingPiece()

		killer := engine.KillerMoveScore(move, moveOrder)
		if killer != 0 {
			mp.quietScores[i] = killer
			continue
		}

		history := engine.MoveHistoryScore(piece, dest, moveOrder)
		if history != 0 {
			mp.quietScores[i] = history
			continue
		}

		// prefer checks
		// if move.IsCheck() {
		// 	mp.quietScores[i] = 10_000
		// 	continue
		// }

		// King safety (castling)
		isCastling := move.IsKingSideCastle() || move.IsQueenSideCastle()
		if isCastling {
			mp.quietScores[i] = 3_000
			continue
		}

		// Prefer smaller pieces
		if piece.Type() == King {
			mp.quietScores[i] = 0
			continue
		}

		mp.quietScores[i] = 1100 - int32(piece.Weight())
	}
}

func (mp *MovePicker) Reset() {
	mp.canUseHashMove = mp.hashmove != EmptyMove
	mp.nextQuiet = 0
	mp.nextCapture = 0
	if mp.canUseHashMove {
		if mp.hashmove.IsCapture() {
			mp.nextCapture = 1
		} else {
			mp.nextQuiet = 1
		}
	}
}

func (mp *MovePicker) Next() Move {
	if mp.hashmove != EmptyMove && mp.canUseHashMove {
		mp.canUseHashMove = false
		return mp.hashmove
	}

	move := mp.getNextCapture()
	if move == EmptyMove {
		return mp.getNextQuiet()
	}
	return move
}

func (mp *MovePicker) getNextCapture() Move {
	if mp.captureMoves == nil {
		mp.generateCaptureMoves()
	}

	if mp.nextCapture >= len(mp.captureMoves) {
		return EmptyMove
	}

	bestIndex := mp.nextCapture
	for i := mp.nextCapture + 1; i < len(mp.captureMoves); i++ {
		if mp.captureScores[i] > mp.captureScores[bestIndex] {
			bestIndex = i
		}
	}
	if mp.captureScores[bestIndex] < 0 {
		alt := mp.getNextQuiet()
		if alt != EmptyMove {
			return alt
		}
	}
	best := mp.captureMoves[bestIndex]
	mp.captureMoves[mp.nextCapture], mp.captureMoves[bestIndex] = mp.captureMoves[bestIndex], mp.captureMoves[mp.nextCapture]
	mp.captureScores[mp.nextCapture], mp.captureScores[bestIndex] = mp.captureScores[bestIndex], mp.captureScores[mp.nextCapture]
	mp.nextCapture += 1
	return best
}

func (mp *MovePicker) getNextQuiet() Move {
	if mp.quietMoves == nil {
		mp.generateQuietMoves()
	}

	if mp.nextQuiet >= len(mp.quietMoves) {
		return EmptyMove
	}

	bestIndex := mp.nextQuiet
	for i := mp.nextQuiet + 1; i < len(mp.quietMoves); i++ {
		if mp.quietScores[i] > mp.quietScores[bestIndex] {
			bestIndex = i
		}
	}
	best := mp.quietMoves[bestIndex]
	mp.quietMoves[mp.nextQuiet], mp.quietMoves[bestIndex] = mp.quietMoves[bestIndex], mp.quietMoves[mp.nextQuiet]
	mp.quietScores[mp.nextQuiet], mp.quietScores[bestIndex] = mp.quietScores[bestIndex], mp.quietScores[mp.nextQuiet]
	mp.nextQuiet += 1
	return best
}
