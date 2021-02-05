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

func (p *Position) ToggleTurn() {
	p.ToggleTag(BlackToMove)
	p.ToggleTag(WhiteToMove)
}

func (p *Position) MakeMove(move Move) {
	p.board.Move(move.source, move.destination)
	movingPiece := p.board.PieceAt(move.source)

	p.enPassant = findEnPassantSquare(move, movingPiece)

	// Do promotion
	if move.promoType != NoType {
		p.board.UpdateSquare(move.destination, getPiece(move.promoType, p.Turn()))
	}

	if movingPiece == BlackKing {
		p.ClearTag(BlackCanCastleKingSide)
		p.ClearTag(BlackCanCastleQueenSide)
	} else if movingPiece == WhiteKing {
		p.ClearTag(WhiteCanCastleKingSide)
		p.ClearTag(WhiteCanCastleQueenSide)
	} else if movingPiece == BlackRook && move.source == A8 {
		p.ClearTag(BlackCanCastleQueenSide)
	} else if movingPiece == BlackRook && move.source == H8 {
		p.ClearTag(BlackCanCastleKingSide)
	} else if movingPiece == WhiteRook && move.source == A1 {
		p.ClearTag(WhiteCanCastleQueenSide)
	} else if movingPiece == WhiteRook && move.source == H1 {
		p.ClearTag(WhiteCanCastleKingSide)
	}

	p.ToggleTurn()
}

func (p *Position) UnMakeMove(move Move, tag PositionTag, enPassant Square, capturedPiece Piece) {
	p.tag = tag
	p.enPassant = enPassant
	p.board.Move(move.destination, move.source)

	// Undo enpassant
	if move.HasTag(EnPassant) && move.HasTag(Capture) {
		p.board.UpdateSquare(enPassant, capturedPiece)
	} else if move.HasTag(Capture) { // Undo capture
		p.board.UpdateSquare(move.destination, capturedPiece)
	}

	// Undo promotion
	if move.promoType != NoType {
		p.board.UpdateSquare(move.source, getPiece(Pawn, p.Turn()))
	}
	if move.HasTag(QueenSideCastle) {
		// white
		if move.destination == C1 {
			p.board.Move(D1, A1)
		} else { // black
			p.board.Move(D8, A8)
		}
	} else if move.HasTag(KingSideCastle) {
		// white
		if move.destination == G1 {
			p.board.Move(F1, H1)
		} else { // black
			p.board.Move(F8, H8)
		}
	}
}

func (p *Position) Hash() uint64 {
	return generateZobristHash(p)
}

func findEnPassantSquare(move Move, movingPiece Piece) Square {
	if !move.HasTag(EnPassant) {
		return NoSquare
	} else if movingPiece == WhitePawn &&
		move.source.Rank() == Rank2 && move.destination.Rank() == Rank4 {
		return SquareOf(move.source.File(), Rank3)
	} else if movingPiece == BlackPawn &&
		move.source.Rank() == Rank7 && move.destination.Rank() == Rank5 {
		return SquareOf(move.source.File(), Rank6)
	}
	return NoSquare
}

func (p *Position) copy() *Position {
	return &Position{
		*p.board.copy(),
		p.enPassant,
		p.tag,
	}
}
