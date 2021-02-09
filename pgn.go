package main

import (
	"fmt"
)

func (p *Position) ParseMoves(moveStr []string) []*Move {
	if len(moveStr) == 0 {
		return []*Move{}
	}
	var parsed *Move
	for _, move := range p.LegalMoves() {
		if move.ToString() == moveStr[0] {
			parsed = move
			break
		}
	}
	if parsed == nil {
		panic(fmt.Sprintf("Expectd a valid move, %s is not valid", moveStr))
	}
	tg := p.tag
	ep := p.enPassant
	cp := p.MakeMove(parsed)
	otherMoves := p.ParseMoves(moveStr[1:])
	p.UnMakeMove(parsed, tg, ep, cp)
	return append(append([]*Move{}, parsed), otherMoves...)
}
