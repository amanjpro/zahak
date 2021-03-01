package engine

import (
	"github.com/brentp/intintmap"
)

type Position struct {
	Board         Bitboard
	EnPassant     Square
	Tag           PositionTag
	hash          uint64
	Positions     intintmap.Map
	HalfMoveClock uint8
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

func (p *Position) MakeNullMove() Square {
	ep := p.EnPassant
	p.EnPassant = NoSquare
	p.HalfMoveClock += 1
	p.ToggleTurn()
	updateHashForNullMove(p, NoSquare, ep)
	v, ok := p.Positions.Get(int64(p.Hash()))
	if ok {
		p.Positions.Put(int64(p.Hash()), v+1)
	} else {
		p.Positions.Put(int64(p.Hash()), 1)
	}
	return ep
}

func (p *Position) UnMakeNullMove(ep Square) {
	v, ok := p.Positions.Get(int64(p.Hash()))
	if ok {
		if v <= 1 {
			p.Positions.Del(int64(p.Hash()))
		} else {
			p.Positions.Put(int64(p.Hash()), v-1)
		}
	}
	p.EnPassant = ep
	p.HalfMoveClock -= 1
	p.ToggleTurn()
}

func (p *Position) ToggleTurn() {
	p.ToggleTag(BlackToMove)
	p.ToggleTag(WhiteToMove)
}

// only for movegen
func (p *Position) partialMakeMove(move Move) Piece {
	capturedPiece := p.Board.PieceAt(move.Destination)
	p.Board.Move(move.Source, move.Destination)

	// EnPassant flag is a form of capture, captures do not result in enpassant allowance
	if move.HasTag(EnPassant) {
		ep := findEnPassantCaptureSquare(move)
		capturedPiece = p.Board.PieceAt(ep)
		p.Board.Clear(ep)
	}

	// Do promotion
	if move.PromoType != NoType {
		promoPiece := GetPiece(move.PromoType, p.Turn())
		p.Board.UpdateSquare(move.Destination, promoPiece)
	}

	p.ToggleTurn()
	return capturedPiece
}

// only for movegen
func (p *Position) partialUnMakeMove(move Move, capturedPiece Piece) {
	p.Board.Move(move.Destination, move.Source)
	// Undo enpassant
	if move.HasTag(EnPassant) {
		cp := findEnPassantCaptureSquare(move)
		p.Board.UpdateSquare(cp, capturedPiece)
	} else if move.HasTag(Capture) { // Undo capture
		p.Board.UpdateSquare(move.Destination, capturedPiece)
	}

	p.ToggleTurn()
	// Undo promotion
	if move.PromoType != NoType {
		movingPiece := GetPiece(Pawn, p.Turn())
		p.Board.UpdateSquare(move.Source, movingPiece)
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

func (p *Position) MakeMove(move Move) (Piece, Square, PositionTag, uint8) {
	hc := p.HalfMoveClock
	ep := p.EnPassant
	tag := p.Tag
	movingPiece := p.Board.PieceAt(move.Source)
	capturedPiece := p.Board.PieceAt(move.Destination)
	p.Board.Move(move.Source, move.Destination)
	captureSquare := NoSquare
	promoPiece := NoPiece

	if movingPiece.Type() == Pawn || capturedPiece != NoPiece {
		p.HalfMoveClock = 0
	} else {
		p.HalfMoveClock += 1
	}

	// EnPassant flag is a form of capture, captures do not result in enpassant allowance
	if move.HasTag(EnPassant) {
		p.EnPassant = NoSquare
		ep := findEnPassantCaptureSquare(move)
		capturedPiece = p.Board.PieceAt(ep)
		captureSquare = ep
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

	if move.HasTag(Capture) && !move.HasTag(EnPassant) {
		captureSquare = move.Destination
	}

	// Do promotion
	if move.PromoType != NoType {
		promoPiece = GetPiece(move.PromoType, p.Turn())
		p.Board.UpdateSquare(move.Destination, promoPiece)
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
	updateHash(p, move, movingPiece, capturedPiece, captureSquare, p.EnPassant, ep, promoPiece, tag)
	v, ok := p.Positions.Get(int64(p.Hash()))
	if ok {
		p.Positions.Put(int64(p.Hash()), v+1)
	} else {
		p.Positions.Put(int64(p.Hash()), 1)
	}
	return capturedPiece, ep, tag, hc
}

func (p *Position) UnMakeMove(move Move, tag PositionTag, enPassant Square, capturedPiece Piece,
	halfClock uint8) {
	oldTag := p.Tag
	oldEnPassant := p.EnPassant
	movingPiece := p.Board.PieceAt(move.Destination)
	promoPiece := movingPiece
	p.Tag = tag
	p.HalfMoveClock = halfClock
	p.EnPassant = enPassant

	v, ok := p.Positions.Get(int64(p.Hash()))
	if ok {
		if v <= 1 {
			p.Positions.Del(int64(p.Hash()))
		} else {
			p.Positions.Put(int64(p.Hash()), v-1)
		}
	}

	captureSquare := NoSquare
	p.Board.Move(move.Destination, move.Source)
	// Undo enpassant
	if move.HasTag(EnPassant) {
		cp := findEnPassantCaptureSquare(move)
		captureSquare = cp
		p.Board.UpdateSquare(cp, capturedPiece)
	} else if move.HasTag(Capture) { // Undo capture
		p.Board.UpdateSquare(move.Destination, capturedPiece)
		captureSquare = move.Destination
	}

	// Undo promotion
	if move.PromoType != NoType {
		movingPiece = GetPiece(Pawn, p.Turn())
		p.Board.UpdateSquare(move.Source, movingPiece)
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
	updateHash(p, move, movingPiece, capturedPiece, captureSquare, p.EnPassant, oldEnPassant, promoPiece, oldTag)
}

type Status uint8

const (
	Checkmate Status = iota
	Draw
	Unknown
)

func (p *Position) IsEndGame() bool {
	return p.Board.IsEndGame()
}

func (p *Position) IsInCheck() bool {
	return isInCheck(p.Board, p.Turn())
}

func (p *Position) Status() Status {
	value, ok := p.Positions.Get(int64(p.Hash()))
	if ok && value >= 3 {
		return Draw
	}
	if p.IsInCheck() {
		if !p.HasLegalMoves() {
			return Checkmate
		}
	} else if p.HalfMoveClock >= 100 {
		return Draw
	} else {
		if !p.HasLegalMoves() {
			return Draw
		} else if p.Board.blackPawn != 0 || p.Board.whitePawn != 0 ||
			p.Board.blackRook != 0 || p.Board.whiteRook != 0 ||
			p.Board.blackQueen != 0 || p.Board.whiteQueen != 0 {
			return Unknown
		} else {
			wKnights := bitScanForward(p.Board.whiteKnight)
			bKnights := bitScanForward(p.Board.blackKnight)
			wBishops := bitScanForward(p.Board.blackBishop)
			bBishops := bitScanForward(p.Board.blackBishop)

			wKnightsNum := 0
			bKnightsNum := 0
			wBishopsNum := 0
			bBishopsNum := 0

			if wKnights != 64 {
				wKnightsNum = 1
			}

			if bKnights != 64 {
				bKnightsNum = 1
			}

			if wBishops != 64 {
				wBishopsNum = 1
			}

			if bBishops != 64 {
				bBishopsNum = 1
			}

			all := wKnightsNum + bKnightsNum + wBishopsNum + bBishopsNum

			// both sides have a bare king
			// one side has a king and a minor piece against a bare king

			if all <= 1 {
				return Draw
			}
			// both sides have a king and a bishop, the bishops being the same color
			if wKnightsNum == 0 && bKnightsNum == 0 {
				otherWB := wBishops ^ (1 << wBishops)
				otherBB := bBishops ^ (1 << bBishops)
				if otherWB == 0 && otherBB == 0 &&
					Square(1<<bBishops).GetColor() == Square(1<<wBishops).GetColor() {
					return Draw
				}
			}
		}
	}

	return Unknown
}

func (p *Position) IsFIDEDrawRule() bool {
	if p.HalfMoveClock >= 100 {
		return true
	}
	value, ok := p.Positions.Get(int64(p.Hash()))
	return (ok && value >= 3)
}

func (p *Position) Hash() uint64 {
	if p.hash == 0 {
		hash := generateZobristHash(p)
		p.hash = hash
	}
	return p.hash
}

func findEnPassantCaptureSquare(move Move) Square {
	rank := move.Source.Rank()
	file := move.Destination.File()
	return SquareOf(file, rank)
}

func (p *Position) copy() *Position {
	copyMap := intintmap.New(10000, 0.5)
	for item := range p.Positions.Items() {
		copyMap.Put(item[0], item[1])
	}
	return &Position{
		*p.Board.copy(),
		p.EnPassant,
		p.Tag,
		p.hash,
		*copyMap,
		p.HalfMoveClock,
	}
}
