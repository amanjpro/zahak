package engine

type Position struct {
	Board     Bitboard
	EnPassant Square
	Tag       PositionTag
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

func (p *Position) SetTag(tag PositionTag)      { p.Tag |= tag }
func (p *Position) ClearTag(tag PositionTag)    { p.Tag &= ^tag }
func (p *Position) ToggleTag(tag PositionTag)   { p.Tag ^= tag }
func (p *Position) HasTag(tag PositionTag) bool { return p.Tag&tag != 0 }

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

func (p *Position) MakeMove(move *Move) (Piece, Square, PositionTag) {
	ep := p.EnPassant
	tag := p.Tag
	movingPiece := p.Board.PieceAt(move.Source)
	capturedPiece := p.Board.PieceAt(move.Destination)
	p.Board.Move(move.Source, move.Destination)

	// EnPassant flag is a form of capture, captures do not result in enpassant allowance
	if move.HasTag(EnPassant) {
		p.EnPassant = NoSquare
		ep := findEnPassantSquare(move, movingPiece)
		capturedPiece = p.Board.PieceAt(ep)
		p.Board.Clear(ep)
	} else {
		if movingPiece == WhitePawn &&
			move.Source.Rank() == Rank2 && move.Destination.Rank() == Rank4 {
			p.EnPassant = SquareOf(move.Source.File(), Rank3)
		} else if movingPiece == BlackPawn &&
			move.Source.Rank() == Rank7 && move.Destination.Rank() == Rank5 {
			p.EnPassant = SquareOf(move.Source.File(), Rank6)
		} else {
			p.EnPassant = NoSquare
		}
	}

	// Do promotion
	if move.PromoType != NoType {
		p.Board.UpdateSquare(move.Destination, getPiece(move.PromoType, p.Turn()))
	}

	if movingPiece == BlackKing {
		p.ClearTag(BlackCanCastleKingSide)
		p.ClearTag(BlackCanCastleQueenSide)
	} else if movingPiece == WhiteKing {
		p.ClearTag(WhiteCanCastleKingSide)
		p.ClearTag(WhiteCanCastleQueenSide)
	} else if movingPiece == BlackRook && move.Source == A8 {
		p.ClearTag(BlackCanCastleQueenSide)
	} else if movingPiece == BlackRook && move.Source == H8 {
		p.ClearTag(BlackCanCastleKingSide)
	} else if movingPiece == WhiteRook && move.Source == A1 {
		p.ClearTag(WhiteCanCastleQueenSide)
	} else if movingPiece == WhiteRook && move.Source == H1 {
		p.ClearTag(WhiteCanCastleKingSide)
	}

	// capturing rook nullifies castling right for the opponent on the rooks side
	if move.Destination == A8 && p.Turn() == White {
		p.ClearTag(BlackCanCastleQueenSide)
	} else if move.Destination == H8 && p.Turn() == White {
		p.ClearTag(BlackCanCastleKingSide)
	} else if move.Destination == A1 && p.Turn() == Black {
		p.ClearTag(WhiteCanCastleQueenSide)
	} else if move.Destination == H1 && p.Turn() == Black {
		p.ClearTag(WhiteCanCastleKingSide)
	}

	p.ToggleTurn()
	return capturedPiece, ep, tag
}

func (p *Position) UnMakeMove(move *Move, tag PositionTag, enPassant Square, capturedPiece Piece) {
	p.Tag = tag
	p.EnPassant = enPassant

	movingPiece := p.Board.PieceAt(move.Destination)
	p.Board.Move(move.Destination, move.Source)

	// Undo enpassant
	if move.HasTag(EnPassant) && move.HasTag(Capture) {
		p.Board.UpdateSquare(findEnPassantSquare(move, movingPiece), capturedPiece)
	} else if move.HasTag(Capture) { // Undo capture
		p.Board.UpdateSquare(move.Destination, capturedPiece)
	}

	// Undo promotion
	if move.PromoType != NoType {
		p.Board.UpdateSquare(move.Source, getPiece(Pawn, p.Turn()))
	}
	if move.HasTag(QueenSideCastle) {
		// white
		if move.Destination == C1 {
			p.Board.Move(D1, A1)
		} else { // black
			p.Board.Move(D8, A8)
		}
	} else if move.HasTag(KingSideCastle) {
		// white
		if move.Destination == G1 {
			p.Board.Move(F1, H1)
		} else { // black
			p.Board.Move(F8, H8)
		}
	}
}

type Status uint8

const (
	Checkmate Status = iota
	Draw
	Unknown
)

func (p *Position) IsInCheck() bool {
	return isInCheck(p.Board, p.Turn())
}

func (p *Position) Status() Status {
	if p.IsInCheck() {
		if !p.HasLegalMoves() {
			return Checkmate
		}
	} else {
		if !p.HasLegalMoves() {
			return Draw
		} else if p.Board.blackPawn != 0 || p.Board.whitePawn != 0 ||
			p.Board.blackRook != 0 || p.Board.whiteRook != 0 ||
			p.Board.blackQueen != 0 || p.Board.whiteQueen != 0 {
			return Unknown
		} else {
			whiteKnights := getIndicesOfOnes(p.Board.whiteKnight)
			blackKnights := getIndicesOfOnes(p.Board.blackKnight)
			whiteBishops := getIndicesOfOnes(p.Board.whiteBishop)
			blackBishops := getIndicesOfOnes(p.Board.blackBishop)
			all := len(whiteKnights) + len(blackKnights) + len(whiteBishops) + len(blackBishops)
			// both sides have a bare king
			// one side has a king and a minor piece against a bare king

			if all <= 1 {
				return Draw
			}
			// both sides have a king and a bishop, the bishops being the same color
			if p.Board.whiteKnight == 0 && p.Board.blackKnight == 0 {
				if len(blackBishops) == 1 && len(whiteBishops) == 1 &&
					Square(blackBishops[0]).GetColor() == Square(whiteBishops[0]).GetColor() {
					return Draw
				}
			}
		}
	}

	return Unknown
}

func (p *Position) Hash() uint64 {
	return generateZobristHash(p)
}

func findEnPassantSquare(move *Move, movingPiece Piece) Square {
	rank := move.Source.Rank()
	file := move.Destination.File()
	return SquareOf(file, rank)
}

func (p *Position) copy() *Position {
	return &Position{
		*p.Board.copy(),
		p.EnPassant,
		p.Tag,
	}
}
