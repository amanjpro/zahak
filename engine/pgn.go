package engine

import (
	"fmt"
	"strings"
)

func (p *Position) ParseMoves(moveStr []string) []Move {
	if len(moveStr) == 0 {
		return []Move{}
	}
	currentMove := moveStr[0]
	if len(strings.TrimSpace(currentMove)) == 0 {
		return p.ParseMoves(moveStr[1:])
	} else {
		var parsed Move
		validMoves := p.PseudoLegalMoves()
		for _, move := range validMoves {
			if move.ToString() == currentMove {
				parsed = move
				break
			}
		}
		if parsed == 0 {
			panic(fmt.Sprintf("Expected a valid move, %s is not valid", currentMove))
		}
		ep, tg, hc, _ := p.MakeMove(parsed)
		otherMoves := p.ParseMoves(moveStr[1:])
		p.UnMakeMove(parsed, tg, ep, hc)
		return append(append([]Move{}, parsed), otherMoves...)
	}
}

func (p *Position) MoveToPGN(move Move) string {
	if move.IsKingSideCastle() {
		return "O-O"
	}
	if move.IsQueenSideCastle() {
		return "O-O-O"
	}
	isCapture := move.IsCapture()
	movingPiece := move.MovingPiece()
	source := move.Source()
	dest := move.Destination()
	promoType := move.PromoType()

	// highly inefficient, but so what
	ambiguity := 0
	var alternativeMove Move
	if movingPiece.Type() != Pawn {
		validMoves := p.PseudoLegalMoves()
		for _, m := range validMoves {
			if m != move && m.MovingPiece() == movingPiece && m.Destination() == dest {
				alternativeMove = m
				ambiguity += 1
			}
		}
	}

	moveStr := ""
	if movingPiece.Type() != Pawn {
		moveStr = movingPiece.Type().Name()
	} else if movingPiece.Type() == Pawn && isCapture {
		moveStr = source.File().Name()
	}

	if ambiguity == 1 {
		s := alternativeMove.Source()
		if s.File() != source.File() {
			moveStr = fmt.Sprint(moveStr, source.File().Name())
		} else {
			moveStr = fmt.Sprint(moveStr, source.Rank().Name())
		}
	} else if ambiguity > 1 {
		moveStr = fmt.Sprint(moveStr, source.File().Name(), source.Rank().Name())
	}
	if isCapture {
		moveStr = fmt.Sprint(moveStr, "x")
	}
	moveStr = fmt.Sprint(moveStr, dest.Name())
	if promoType != NoType {
		moveStr = fmt.Sprint(moveStr, "=", promoType.Name())
	}
	// is Checkmate?
	p.partialMakeMove(move)
	if p.IsInCheck() {
		if len(p.PseudoLegalMoves()) == 0 {
			moveStr = fmt.Sprint(moveStr, "#")
		} else {
			moveStr = fmt.Sprint(moveStr, "+")
		}
	}
	p.partialUnMakeMove(move)

	return moveStr
}
