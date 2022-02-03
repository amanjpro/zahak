package engine

import (
	"fmt"
	"strings"
)

func (p *Position) ParseSearchMoves(moveStr []string) []Move {
	parsed := make([]Move, 0, len(moveStr))
	validMoves := p.PseudoLegalMoves()
	for _, move := range validMoves {
		for _, cm := range moveStr {
			if move.ToString() == cm {
				parsed = append(parsed, move)
			}
		}
	}
	return parsed
}

func (g *Game) ParseGameMoves(moveStr []string) {
	p := g.Position()
	for _, currentMove := range moveStr {
		currentMove := strings.TrimSpace(currentMove)
		validMoves := p.PseudoLegalMoves()
		for _, move := range validMoves {
			if move.ToString() == currentMove {
				_, _, _, _ = p.GameMakeMove(move)
				if p.Turn() == White {
					g.numberOfMoves += 1
				}
				v, ok := p.Positions[p.Hash()]
				if ok {
					p.Positions[p.Hash()] = v + 1
				} else {
					p.Positions[p.Hash()] = 1
				}
				goto end
			}
		}
		panic(fmt.Sprintf("Illegal moves are in the path: %s", currentMove))
	end:
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
