package search

import (
	. "github.com/amanjpro/zahak/engine"
)

func dynamicMargin(pos *Position) int16 {

	color := pos.Turn()
	delta := p

	if color == White {
		if pos.Board.GetBitboardOf(WhitePawn)&Rank7Fill != 0 {
			delta = q
		}
	} else {
		if pos.Board.GetBitboardOf(BlackPawn)&Rank2Fill != 0 {
			delta = q
		}
	}

	other := color.Other()
	if pos.Board.GetBitboardOf(GetPiece(Queen, other)) != 0 {
		return delta + q
	}

	if pos.Board.GetBitboardOf(GetPiece(Rook, other)) != 0 {
		return delta + r
	}

	if pos.Board.GetBitboardOf(GetPiece(Bishop, other)) != 0 || pos.Board.GetBitboardOf(GetPiece(Knight, other)) != 0 {
		return delta + b
	}

	return delta + p
}

func (e *Engine) quiescence(alpha int16, beta int16, searchHeight int8) int16 {

	e.VisitNode(searchHeight)

	position := e.Position
	e.tt.Prefetch(position.Hash())
	// pawnhash := e.Pawnhash

	currentMove := e.positionMoves[searchHeight]
	// Position is drawn
	if IsRepetition(position, e.pred, currentMove) || position.IsDraw() {
		return 0
	}

	hash := position.Hash()
	_, nEval, _, nType, ttHit := e.tt.Get(hash)
	// if ttHit {
	// 	// ttHit = position.IsPseudoLegal(nHashMove)
	// }
	// isNoisy := nHashMove.IsCapture() || nHashMove.PromoType() != NoType
	// if !ttHit || !isNoisy {
	// nHashMove = EmptyMove
	// }
	if ttHit {
		nEval = evalFromTT(nEval, searchHeight)
		if nEval >= beta && nType == LowerBound {
			e.CacheHit()
			return nEval
		}
		if nEval <= alpha && nType == UpperBound {
			e.CacheHit()
			return nEval
		}
		if nType == Exact {
			e.CacheHit()
			return nEval
		}
	}

	standPat := e.staticEvals[searchHeight]
	if standPat >= beta {
		return standPat // fail soft
	}

	if searchHeight >= MAX_DEPTH-1 {
		return standPat
	}

	if (e.isMainThread && e.TimeManager().ShouldStop(false, false)) || (!e.isMainThread && e.parent.Stop) {
		return 0
	}

	var isInCheck = e.Position.IsInCheck()
	bestscore := -CHECKMATE_EVAL + int16(searchHeight)
	if !isInCheck {
		bestscore = standPat
	}

	// Delta Pruning
	if standPat+dynamicMargin(position) < alpha {
		return alpha
	}

	if alpha < standPat {
		alpha = standPat
	}

	bestMove := EmptyMove
	originalAlpha := alpha

	movePicker := e.MovePickers[searchHeight]
	movePicker.RecycleWith(position, e, searchHeight, EmptyMove, true)

	noisyMoves := -1
	seeScores := movePicker.captureMoveList.Scores

	for true {
		move := movePicker.Next()
		if move == EmptyMove {
			break
		}

		// isCaptureMove := move.IsCapture()
		// if isCaptureMove || move.PromoType() != NoType {
		noisyMoves += 1
		// }

		if /* isCaptureMove && */ seeScores[noisyMoves] < 0 {
			// SEE pruning
			break
		}

		// if !IsPromoting(move) {
		// 	margin := p + move.CapturedPiece().Weight()
		// 	if standPat+margin <= alpha {
		// 		continue
		// 	}
		// }

		if ep, tg, hc, ok := position.MakeMove(move); ok {
			e.positionMoves[searchHeight+1] = move
			e.staticEvals[searchHeight+1] = position.Evaluate()

			e.pred.Push(position.Hash())
			score := -e.quiescence(-beta, -alpha, searchHeight+1)
			e.pred.Pop()
			position.UnMakeMove(move, tg, ep, hc)
			if score > bestscore {
				bestMove = move
				bestscore = score
				if score > alpha {
					alpha = score
					if score >= beta {
						break
					}
				}
			}
		}
	}

	if (e.isMainThread && !e.TimeManager().AbruptStop) || (!e.isMainThread && !e.parent.Stop) {
		flag := Exact
		if bestscore >= beta {
			flag = LowerBound
		} else if bestscore <= originalAlpha {
			flag = UpperBound
		}
		e.tt.Set(hash, bestMove, evalToTT(bestscore, searchHeight), 0, flag)
	}

	return bestscore
}

func (e *Engine) Quiescence(alpha int16, beta int16, searchHeight int8) int16 {
	return e.quiescence(alpha, beta, searchHeight)
}
