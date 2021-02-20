package search

import (
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

func quiescence(position *Position, alpha int16, beta int16, lastDepth int8, ply int8, standPat int16) int16 {

	outcome := position.Status()
	if outcome == Checkmate {
		return (-CHECKMATE_EVAL + int16(lastDepth+ply))
	} else if outcome == Draw {
		return 0
	}

	withChecks := ply <= 4
	legalMoves := position.QuiesceneMoves(withChecks)
	orderedMoves := orderMoves(&ValidMoves{position, legalMoves, 125})

	if standPat >= beta {
		return beta // fail hard
	}

	if alpha < standPat {
		alpha = standPat
	}

	if STOP_SEARCH_GLOBALLY {
		return standPat
	}

	w := WhitePawn
	deltaMargin := w.Weight() * 2 // 200 centipawns
	for _, move := range orderedMoves {
		if move.HasTag(Capture) && !move.HasTag(EnPassant) {
			// SEE pruning
			board := position.Board
			movingPiece := board.PieceAt(move.Source)
			capturedPiece := board.PieceAt(move.Destination)
			gain := position.Board.StaticExchangeEval(move.Destination, capturedPiece, move.Source, movingPiece)
			if gain <= 0 {
				continue
			}
		}
		cp, ep, tg, hc := position.MakeMove(move)
		sp := Evaluate(position)
		if cp != NoPiece && standPat < alpha-deltaMargin { // is capture
			// Delta pruning meaningless captures
			position.UnMakeMove(move, tg, ep, cp, hc)
			continue
		}
		score := -quiescence(position, -beta, -alpha, lastDepth, ply+1, sp)
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
