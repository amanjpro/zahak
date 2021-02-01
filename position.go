package main

type Position struct {
	board     Bitboard
	enPassant Square
	tag       PositionTag
}

type PositionTag uint8

const (
	WhiteCanCastleKingSide PositionTag = 1 << iota
	WhiteCanCastleQueenSide
	BlackCanCastleKingSide
	BlackCanCastleQueenSide
	BlackToMove
	WhiteToMove
)

func (p *Position) SetTag(tag PositionTag)      { p.tag |= tag }
func (p *Position) ClearTag(tag PositionTag)    { p.tag &= ^tag }
func (p *Position) ToggleTag(tag PositionTag)   { p.tag ^= tag }
func (p *Position) HasTag(tag PositionTag) bool { return p.tag&tag != 0 }

func (p *Position) Turn() Color {
	if p.HasTag(WhiteToMove) {
		return White
	}
	return Black
}

func (p *Position) MakeMove(move Move) {
	p.board.Move(move.source, move.destination)
	movingPiece := p.board.PieceAt(move.source)

	if movingPiece == BlackKing {
		p.ClearTag(BlackCanCastleKingSide)
		p.ClearTag(BlackCanCastleQueenSide)
	} else if movingPiece == WhiteKing {
		p.ClearTag(WhiteCanCastleKingSide)
		p.ClearTag(WhiteCanCastleQueenSide)
	} else if movingPiece == BlackRook && *(move.source) == A8 {
		p.ClearTag(BlackCanCastleQueenSide)
	} else if movingPiece == BlackRook && *(move.source) == H8 {
		p.ClearTag(BlackCanCastleKingSide)
	} else if movingPiece == WhiteRook && *(move.source) == A1 {
		p.ClearTag(WhiteCanCastleQueenSide)
	} else if movingPiece == WhiteRook && *(move.source) == H1 {
		p.ClearTag(WhiteCanCastleKingSide)
	}

	p.ToggleTag(BlackToMove)
	p.ToggleTag(WhiteToMove)
}

func (p *Position) UnMakeMove(move Move, tag PositionTag, enPassant Square, capturedPiece Piece) {
	p.tag = tag
	p.enPassant = enPassant
	p.board.Move(move.destination, move.source)
	p.board.UpdateSquare(move.destination, capturedPiece)
	if move.HasTag(QueenSideCastle) {
		// white
		if *move.destination == C1 {
			p.UnMakeMove(Move{&A1, &D1, NoType, 0}, tag, enPassant, NoPiece)
		} else { // black
			p.UnMakeMove(Move{&A8, &D8, NoType, 0}, tag, enPassant, NoPiece)
		}
	} else if move.HasTag(KingSideCastle) {
		// white
		if *move.destination == G1 {
			p.UnMakeMove(Move{&H1, &F1, NoType, 0}, tag, enPassant, NoPiece)
		} else { // black
			p.UnMakeMove(Move{&H8, &F8, NoType, 0}, tag, enPassant, NoPiece)
		}
	}
}

func (p *Position) Hash() uint64 {
	return generateZobristHash(p)
}
