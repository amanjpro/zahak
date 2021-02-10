package engine

import (
	"fmt"
	"strings"
)

func (p *Position) ParseMoves(moveStr []string) []*Move {
	if len(moveStr) == 0 {
		return []*Move{}
	}
	currentMove := moveStr[0]
	if len(strings.TrimSpace(currentMove)) == 0 {
		return p.ParseMoves(moveStr[1:])
	} else {
		var parsed *Move
		for _, move := range p.LegalMoves() {
			if move.ToString() == currentMove {
				parsed = move
				break
			}
		}
		if parsed == nil {
			panic(fmt.Sprintf("Expectd a valid move, %s is not valid", currentMove))
		}
		cp, ep, tg := p.MakeMove(parsed)
		otherMoves := p.ParseMoves(moveStr[1:])
		p.UnMakeMove(parsed, tg, ep, cp)
		return append(append([]*Move{}, parsed), otherMoves...)
	}
}
