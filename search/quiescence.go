package search

import (
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

func (e *Engine) quiescence(position *Position, alpha int32, beta int32, ply int8, standPat int32) int32 {

	outcome := position.Status()
	if outcome == Checkmate {
		return -CHECKMATE_EVAL
	} else if outcome == Draw {
		return 0
	}

	withChecks := ply <= 4
	legalMoves := position.QuiesceneMoves(withChecks)
	orderedMoves := orderMoves(&ValidMoves{position, e.pv, legalMoves, 125})

	if standPat >= beta {
		return beta // fail hard
	}

	if alpha < standPat {
		alpha = standPat
	}

	if e.StopSearchFlag {
		return standPat
	}

	isInCheck := position.IsInCheck()

	w := WhitePawn
	deltaMargin := w.Weight() * 2 // 200 centipawns
	for _, move := range orderedMoves {
		if !isInCheck && move.HasTag(Capture) && !move.HasTag(EnPassant) {
			// SEE pruning
			board := position.Board
			movingPiece := board.PieceAt(move.Source)
			capturedPiece := board.PieceAt(move.Destination)
			gain := board.StaticExchangeEval(move.Destination, capturedPiece, move.Source, movingPiece)
			if gain < 0 {
				continue
			}
		}
		cp, ep, tg, hc := position.MakeMove(move)
		if !isInCheck && cp != NoPiece && standPat < alpha-deltaMargin { // is capture
			// Delta pruning meaningless captures
			position.UnMakeMove(move, tg, ep, cp, hc)
			continue
		}
		sp := Evaluate(position)
		score := -e.quiescence(position, -beta, -alpha, ply+1, sp)
		position.UnMakeMove(move, tg, ep, cp, hc)
		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha = score
		}
	}
	return alpha
}
