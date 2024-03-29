package engine

import (
	"math/bits"
)

const CHECKMATE_EVAL int16 = 30000
const MAX_NON_CHECKMATE float32 = 25000
const MIN_NON_CHECKMATE float32 = -MAX_NON_CHECKMATE
const CASTLING_FLAG = WhiteCanCastleQueenSide | WhiteCanCastleKingSide | BlackCanCastleQueenSide | BlackCanCastleKingSide

var phaseValues = []int{0, 1, 1, 2, 4, 0, 0, 1, 1, 2, 4, 0}
var maxPhase int = 24

type Position struct {
	Board         Bitboard
	Net           NetworkState
	Updates       Updates
	EnPassant     Square
	Tag           PositionTag
	hash          uint64
	Positions     map[uint64]int
	HalfMoveClock uint8
	phase         int
	hashAcc       []uint64
	hashHeight    int
}

type PositionTag uint8

const (
	WhiteCanCastleKingSide PositionTag = 1 << iota
	WhiteCanCastleQueenSide
	BlackCanCastleKingSide
	BlackCanCastleQueenSide
	InCheck
	BlackToMove
	WhiteToMove
)

func (p *Position) SetTag(tag PositionTag)      { p.Tag |= tag }
func (p *Position) ClearTag(tag PositionTag)    { p.Tag &= ^tag }
func (p *Position) ToggleTag(tag PositionTag)   { p.Tag ^= tag }
func (p *Position) HasTag(tag PositionTag) bool { return p.Tag&tag != 0 }

func (p *Position) HasCastling() bool {
	return p.HasTag(CASTLING_FLAG)
}

func (p *Position) Turn() Color {
	if p.HasTag(WhiteToMove) {
		return White
	}
	return Black
}

func (p *Position) MakeNullMove() Square {
	ep := p.EnPassant
	p.Updates.Size = 0
	p.EnPassant = NoSquare
	p.HalfMoveClock += 1
	p.ToggleTurn()
	updateHashForNullMove(p, NoSquare, ep)
	if p.Turn() == White {
		p.Updates.Add(768, Add)
	} else {
		p.Updates.Add(768, Remove)
	}
	p.Net.UpdateHidden(&p.Updates)
	return ep
}

func (p *Position) UnMakeNullMove(ep Square) {
	updateHashForNullMove(p, NoSquare, ep)
	p.EnPassant = ep
	p.HalfMoveClock -= 1
	p.ToggleTurn()
	p.Net.RevertHidden()
}

func (p *Position) ToggleTurn() {
	p.ToggleTag(BlackToMove)
	p.ToggleTag(WhiteToMove)
}

// only for movegen
func (p *Position) partialMakeMove(move Move) {
	source := move.Source()
	dest := move.Destination()
	movingPiece := move.MovingPiece()
	cp := move.CapturedPiece()

	// EnPassant flag is a form of capture, captures do not result in enpassant allowance
	if move.IsEnPassant() {
		ep := findEnPassantCaptureSquare(move)
		p.Board.Move(source, dest, movingPiece, NoPiece)
		p.Board.Clear(ep, cp)
	} else {
		p.Board.Move(source, dest, movingPiece, cp)
	}

	// Do promotion
	promoType := move.PromoType()
	if promoType != NoType {
		promoPiece := GetPiece(promoType, p.Turn())
		p.Board.UpdateSquare(dest, promoPiece, movingPiece)
	}

	p.ToggleTurn()
}

// only for movegen
func (p *Position) partialUnMakeMove(move Move) {
	p.ToggleTurn()
	movingPiece := move.MovingPiece()
	capturedPiece := move.CapturedPiece()
	source := move.Source()
	dest := move.Destination()

	// Undo promotion
	promoType := move.PromoType()
	if promoType != NoType {
		promoPiece := GetPiece(promoType, p.Turn())
		p.Board.UpdateSquare(dest, movingPiece, promoPiece)
	}

	p.Board.Move(dest, source, movingPiece, NoPiece)
	// Undo enpassant
	if move.IsEnPassant() {
		cp := findEnPassantCaptureSquare(move)
		p.Board.UpdateSquare(cp, capturedPiece, NoPiece)
	} else if move.IsCapture() { // Undo capture
		p.Board.UpdateSquare(dest, capturedPiece, NoPiece)
	}

	if move.IsQueenSideCastle() {
		// white
		if dest == C1 {
			p.Board.Move(D1, A1, WhiteRook, NoPiece)
		} else { // black
			p.Board.Move(D8, A8, BlackRook, NoPiece)
		}
	} else if move.IsKingSideCastle() {
		// white
		if dest == G1 {
			p.Board.Move(F1, H1, WhiteRook, NoPiece)
		} else { // black
			p.Board.Move(F8, H8, BlackRook, NoPiece)
		}
	}
}

func (p *Position) NetInput() []int16 {
	input := make([]int16, 0, 33)

	for j := 0; j < 64; j++ {
		sq := Square(j)

		piece := p.Board.PieceAt(sq)
		if piece != NoPiece {
			input = append(input, calculateNetInputIndex(sq, piece))
		}
	}
	if p.Turn() == White {
		input = append(input, 768)
	}
	return input
}

func (p *Position) GameMakeMove(move Move) (Square, PositionTag, uint8, bool) {
	return p.makeMoveHelper(move, false)
}

func (p *Position) MakeMove(move Move) (Square, PositionTag, uint8, bool) {
	return p.makeMoveHelper(move, true)
}

func (p *Position) makeMoveHelper(move Move, updateHidden bool) (Square, PositionTag, uint8, bool) {
	p.Updates.Size = 0
	hc := p.HalfMoveClock
	ep := p.EnPassant
	tag := p.Tag
	movingPiece := move.MovingPiece()
	capturedPiece := move.CapturedPiece()
	source := move.Source()
	dest := move.Destination()
	captureSquare := NoSquare
	promoPiece := NoPiece

	p.Board.Move(source, dest, movingPiece, NoPiece)
	p.Updates.Add(calculateNetInputIndex(source, movingPiece), Remove)

	if movingPiece.Type() == Pawn || capturedPiece != NoPiece {
		p.HalfMoveClock = 0
	} else {
		p.HalfMoveClock += 1
	}

	// EnPassant flag is a form of capture, captures do not result in enpassant allowance
	if move.IsEnPassant() {
		p.EnPassant = NoSquare
		ep := findEnPassantCaptureSquare(move)
		captureSquare = ep
		p.Board.Clear(ep, capturedPiece)
		p.Updates.Add(calculateNetInputIndex(captureSquare, capturedPiece), Remove)
		// p.phase -= phaseValues[capturedPiece-1]
	} else if move.IsCapture() {
		captureSquare = dest
		p.Board.Clear(dest, capturedPiece)
		p.Updates.Add(calculateNetInputIndex(captureSquare, capturedPiece), Remove)
		p.phase -= phaseValues[capturedPiece-1]
	}

	if movingPiece == WhitePawn &&
		source.Rank() == Rank2 && dest.Rank() == Rank4 {
		p.EnPassant = SquareOf(source.File(), Rank3)
	} else if movingPiece == BlackPawn &&
		source.Rank() == Rank7 && dest.Rank() == Rank5 {
		p.EnPassant = SquareOf(source.File(), Rank6)
	} else {
		p.EnPassant = NoSquare
	}

	// Do promotion
	turn := p.Turn()
	promoType := move.PromoType()
	if promoType != NoType {
		promoPiece = GetPiece(promoType, turn)
		p.Board.UpdateSquare(dest, promoPiece, movingPiece)
		p.Updates.Add(calculateNetInputIndex(dest, promoPiece), Add)
		p.phase += phaseValues[promoPiece-1]
		// p.phase -= phaseValues[movingPiece-1]
	} else {
		p.Updates.Add(calculateNetInputIndex(dest, movingPiece), Add)
	}

	if move.IsQueenSideCastle() {
		if turn == White {
			p.Updates.Add(calculateNetInputIndex(A1, WhiteRook), Remove)
			p.Updates.Add(calculateNetInputIndex(D1, WhiteRook), Add)
		} else {
			p.Updates.Add(calculateNetInputIndex(A8, BlackRook), Remove)
			p.Updates.Add(calculateNetInputIndex(D8, BlackRook), Add)
		}
	} else if move.IsKingSideCastle() {
		if turn == White {
			p.Updates.Add(calculateNetInputIndex(H1, WhiteRook), Remove)
			p.Updates.Add(calculateNetInputIndex(F1, WhiteRook), Add)
		} else {
			p.Updates.Add(calculateNetInputIndex(H8, BlackRook), Remove)
			p.Updates.Add(calculateNetInputIndex(F8, BlackRook), Add)
		}
	}

	if movingPiece == BlackKing {
		p.ClearTag(BlackCanCastleKingSide)
		p.ClearTag(BlackCanCastleQueenSide)
	} else if movingPiece == WhiteKing {
		p.ClearTag(WhiteCanCastleKingSide)
		p.ClearTag(WhiteCanCastleQueenSide)
	} else if movingPiece == BlackRook && source == A8 {
		p.ClearTag(BlackCanCastleQueenSide)
	} else if movingPiece == BlackRook && source == H8 {
		p.ClearTag(BlackCanCastleKingSide)
	} else if movingPiece == WhiteRook && source == A1 {
		p.ClearTag(WhiteCanCastleQueenSide)
	} else if movingPiece == WhiteRook && source == H1 {
		p.ClearTag(WhiteCanCastleKingSide)
	}

	// capturing rook nullifies castling right for the opponent on the rooks side
	if dest == A8 && p.Turn() == White {
		p.ClearTag(BlackCanCastleQueenSide)
	} else if dest == H8 && p.Turn() == White {
		p.ClearTag(BlackCanCastleKingSide)
	} else if dest == A1 && p.Turn() == Black {
		p.ClearTag(WhiteCanCastleQueenSide)
	} else if dest == H1 && p.Turn() == Black {
		p.ClearTag(WhiteCanCastleKingSide)
	}

	if isInCheck(&p.Board, p.Turn()) { // Is the move legal (bad way to determine legality
		p.unMakeMoveHelper(move, tag, ep, hc, false)
		return NoSquare, 0, 0, false
	}

	// movingSide := p.Turn()
	p.ToggleTurn()

	// Set check tag
	if isInCheck(&p.Board, p.Turn()) {
		p.SetTag(InCheck)
	} else {
		p.ClearTag(InCheck)
	}

	if p.Turn() == White {
		p.Updates.Add(768, Add)
	} else {
		p.Updates.Add(768, Remove)
	}

	if updateHidden {
		p.hashAcc[p.hashHeight] = p.hash
		p.Net.UpdateHidden(&p.Updates)
		p.hashHeight += 1
	}

	updateHash(p, move, captureSquare, p.EnPassant, ep, promoPiece, tag)
	return ep, tag, hc, true
}

func (p *Position) GameUnMakeMove(move Move, tag PositionTag, enPassant Square, halfClock uint8) {
	p.unMakeMoveHelper(move, tag, enPassant, halfClock, true)
}

func (p *Position) UnMakeMove(move Move, tag PositionTag, enPassant Square, halfClock uint8) {
	p.unMakeMoveHelper(move, tag, enPassant, halfClock, true)
	p.Net.RevertHidden()
	p.hashHeight -= 1
	p.hash = p.hashAcc[p.hashHeight]
}

func (p *Position) unMakeMoveHelper(move Move, tag PositionTag, enPassant Square, halfClock uint8, isLegal bool) {
	movingPiece := move.MovingPiece()
	capturedPiece := move.CapturedPiece()
	promoPiece := NoPiece
	p.Tag = tag
	p.HalfMoveClock = halfClock
	p.EnPassant = enPassant
	source := move.Source()
	dest := move.Destination()
	// Undo promotion
	promoType := move.PromoType()
	if promoType != NoType {
		promoPiece = GetPiece(promoType, p.Turn())
		p.Board.UpdateSquare(dest, movingPiece, promoPiece)
		p.phase -= phaseValues[promoPiece-1]
		// p.phase += phaseValues[movingPiece-1]
	}
	p.Board.Move(dest, source, movingPiece, NoPiece)

	// captureSquare := NoSquare
	// Undo enpassant
	if move.IsEnPassant() {
		cp := findEnPassantCaptureSquare(move)
		// captureSquare = cp
		p.Board.UpdateSquare(cp, capturedPiece, NoPiece)
		// p.phase += phaseValues[capturedPiece-1]
	} else if move.IsCapture() { // Undo capture
		// captureSquare = dest
		p.Board.UpdateSquare(dest, capturedPiece, NoPiece)
		p.phase += phaseValues[capturedPiece-1]
	}

	if move.IsQueenSideCastle() {
		// white
		if dest == C1 {
			p.Board.Move(D1, A1, WhiteRook, NoPiece)
		} else { // black
			p.Board.Move(D8, A8, BlackRook, NoPiece)
		}
	} else if move.IsKingSideCastle() {
		// white
		if dest == G1 {
			p.Board.Move(F1, H1, WhiteRook, NoPiece)
		} else { // black
			p.Board.Move(F8, H8, BlackRook, NoPiece)
		}
	}

	// if isLegal {
	// }
}

type Status uint8

const (
	Checkmate Status = iota
	Draw
	Unknown
)

func (p *Position) IsEndGame() bool {
	return p.Board.IsEndGame(p.Turn())
}

func (p *Position) IsInCheck() bool {
	return p.HasTag(InCheck)
}

func (p *Position) IsDraw() bool {
	if p.HalfMoveClock > 100 {
		return true
	} else {
		if p.Board.blackPawn != 0 || p.Board.whitePawn != 0 ||
			p.Board.blackRook != 0 || p.Board.whiteRook != 0 ||
			p.Board.blackQueen != 0 || p.Board.whiteQueen != 0 {
			return false
		} else {
			wKnights := bitScanForward(p.Board.whiteKnight)
			bKnights := bitScanForward(p.Board.blackKnight)
			wBishops := bitScanForward(p.Board.whiteBishop)
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
				if wKnightsNum != 0 {
					return bits.OnesCount64(p.Board.whiteKnight) == 1
				} else if bKnightsNum != 0 {
					return bits.OnesCount64(p.Board.blackKnight) == 1
				} else if wBishopsNum != 0 {
					return bits.OnesCount64(p.Board.whiteBishop) == 1
				} else if bBishopsNum != 0 {
					return bits.OnesCount64(p.Board.blackBishop) == 1
				}
			}
			// both sides have a king and a bishop, the bishops being the same color
			if wKnightsNum == 0 && bKnightsNum == 0 {
				otherWB := p.Board.whiteBishop ^ (1 << wBishops)
				otherBB := p.Board.blackBishop ^ (1 << bBishops)
				if otherWB == 0 && otherBB == 0 &&
					Square(bBishops).GetColor() == Square(wBishops).GetColor() {
					return true
				}
			}
		}
	}

	return false
}

func (p *Position) IsFIDEDrawRule() bool {
	if p.HalfMoveClock >= 100 {
		return true
	}
	value, ok := p.Positions[p.Hash()]
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
	return move.Destination() ^ 8
}

func (p *Position) recalculatePhase() {
	p.phase = int(bits.OnesCount64(p.Board.whitePawn))*phaseValues[WhitePawn-1] +
		int(bits.OnesCount64(p.Board.whiteKnight))*phaseValues[WhiteKnight-1] +
		int(bits.OnesCount64(p.Board.whiteBishop))*phaseValues[WhiteBishop-1] +
		int(bits.OnesCount64(p.Board.whiteRook))*phaseValues[WhiteRook-1] +
		int(bits.OnesCount64(p.Board.whiteQueen))*phaseValues[WhiteQueen-1] +
		int(bits.OnesCount64(p.Board.whiteKing))*phaseValues[WhiteKing-1] +
		int(bits.OnesCount64(p.Board.blackPawn))*phaseValues[BlackPawn-1] +
		int(bits.OnesCount64(p.Board.blackKnight))*phaseValues[BlackKnight-1] +
		int(bits.OnesCount64(p.Board.blackBishop))*phaseValues[BlackBishop-1] +
		int(bits.OnesCount64(p.Board.blackRook))*phaseValues[BlackRook-1] +
		int(bits.OnesCount64(p.Board.blackQueen))*phaseValues[BlackQueen-1] +
		int(bits.OnesCount64(p.Board.blackKing))*phaseValues[BlackKing-1]
}

func (p *Position) Copy() *Position {
	copyMap := make(map[uint64]int, len(p.Positions))
	for k, v := range p.Positions {
		copyMap[k] = v
	}
	newUpdates := Updates{
		Indices: make([]int16, 5),
		Coeffs:  make([]int8, 5),
		Size:    0,
	}

	newPos := &Position{
		p.Board.copy(),
		NewNetworkState(),
		newUpdates,
		p.EnPassant,
		p.Tag,
		p.hash,
		copyMap,
		p.HalfMoveClock,
		p.phase,
		make([]uint64, MaximumDepth),
		0,
	}
	newPos.Net.Recalculate(newPos.NetInput())
	return newPos
}
