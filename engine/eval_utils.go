package engine

import (
	"math/bits"
)

/**
A bunch of bittwiddling functions mainly meant for evaluation and search
This is not about move generation
*/

func (b *Bitboard) attacksTo(occupied uint64, sq Square) uint64 {
	knights := b.blackKnight | b.whiteKnight
	kings := b.blackKing | b.whiteKing
	rooksQueens := b.blackQueen | b.whiteQueen
	bishopsQueens := rooksQueens
	rooksQueens |= b.blackRook | b.whiteRook
	bishopsQueens |= b.blackBishop | b.whiteBishop

	sqMask := SquareMask[sq]
	return wPawnsAble2CaptureAny(sqMask, b.blackPawn) |
		bPawnsAble2CaptureAny(sqMask, b.whitePawn) |
		(computedKnightAttacks[sq] & knights) |
		(computedKingAttacks[sq] & kings) |
		(bishopAttacks(sq, occupied, empty) & bishopsQueens) |
		(rookAttacks(sq, occupied, empty) & rooksQueens)
}

func (b *Bitboard) getLeastValuablePiece(attacks uint64, color Color) (uint64, Piece) {
	shift := int8(0)
	if color == Black {
		shift = int8(6)
	}
	start := int8(WhitePawn) + shift
	finish := int8(WhiteKing) + shift

	for piece := start; piece <= finish; piece++ {
		bb := b.GetBitboardOf(Piece(piece))
		subset := attacks & bb
		if subset != 0 {
			return subset & -subset, Piece(piece) // The piece and its location on the board
		}
	}
	return 0, NoPiece
}

func (b *Bitboard) StaticExchangeEval(toSq Square, target Piece, frSq Square, aPiece Piece) int16 {

	gain := make([]int16, 32)
	d := 0

	mayXray := b.blackBishop | b.whiteBishop |
		b.blackRook | b.whiteRook | b.blackQueen | b.whiteQueen

	fromSet := SquareMask[frSq]
	occupied := b.whitePieces | b.blackPieces
	attacks := b.attacksTo(occupied, toSq)

	// Ray Attacks, to update the attack def
	rooksQueens := b.blackQueen | b.whiteQueen
	bishopsQueens := rooksQueens
	rooksQueens |= b.blackRook | b.whiteRook
	bishopsQueens |= b.blackBishop | b.whiteBishop

	gain[d] = target.Weight()

	for fromSet != 0 {
		d++ // next depth and side
		color := aPiece.Color()
		gain[d] = aPiece.Weight() - gain[d-1] // speculative store, if defended
		if max(-gain[d-1], gain[d]) < 0 {
			break // pruning does not influence the result
		}
		attacks ^= fromSet  // reset bit in set to traverse
		occupied ^= fromSet // reset bit in temporary occupancy (for x-Rays)
		if fromSet&mayXray != 0 {
			bishopsQueens &^= fromSet // reset bit in temporary occupancy for bishops/queens
			rooksQueens &^= fromSet   // reset bit in temporary occupancy for rooks/queens
			attacks |= (bishopAttacks(toSq, occupied, empty) & bishopsQueens)
			attacks |= (rookAttacks(toSq, occupied, empty) & rooksQueens)
		}
		fromSet, aPiece = b.getLeastValuablePiece(attacks, color.Other())
	}
	for d--; d > 0; d-- {
		gain[d-1] = -max(-gain[d-1], gain[d])
	}
	return gain[0]
}

func (b *Bitboard) IsBackwardPawn(pawn uint64, bb uint64, color Color) bool {
	if color == White {
		return (noEaOne(pawn)&bb) != 0 && (noWeOne(pawn)&bb) != 0
	}
	if color == Black {
		return (soEaOne(pawn)&bb) != 0 && (soWeOne(pawn)&bb) != 0
	}
	return false
}

func (b *Bitboard) IsHorizontalDoubleRook(sq Square, otherRooks uint64, occupied uint64) bool {
	horizontalAttacks := getPositiveRayAttacks(sq, occupied, North) |
		getNegativeRayAttacks(sq, occupied, South)
	return (horizontalAttacks & otherRooks) != 0
}

func (b *Bitboard) IsVerticalDoubleRook(sq Square, otherRooks uint64, occupied uint64) bool {
	horizontalAttacks := getPositiveRayAttacks(sq, occupied, East) |
		getNegativeRayAttacks(sq, occupied, West)
	return (horizontalAttacks & otherRooks) != 0
}

func (b *Bitboard) AllAttacks(color Color) (uint64, uint64, uint64) {
	var ownPawns, ownKnights, ownR, ownB, ownQ, ownKing, ownPieces uint64
	occupiedBB := b.whitePieces | b.blackPieces
	if color == Black {
		ownPieces = b.blackPieces
		ownPawns = bPawnsAble2CaptureAny(b.blackPawn, universal)
		ownKnights = b.blackKnight
		ownR = b.blackRook
		ownB = b.blackBishop
		ownQ = b.blackQueen
		ownKing = b.blackKing
	} else {
		ownPieces = b.whitePieces
		ownPawns = wPawnsAble2CaptureAny(b.whitePawn, universal)
		ownKnights = b.whiteKnight
		ownR = b.whiteRook
		ownB = b.whiteBishop
		ownQ = b.whiteQueen
		ownKing = b.whiteKing
	}
	pawnAttacks := ownPawns
	minorAttacks := (knightAttacks(ownKnights))
	otherAttacks := kingAttacks(ownKing)
	for ownB != 0 {
		sq := bitScanForward(ownB)
		minorAttacks |= bishopAttacks(Square(sq), occupiedBB, ownPieces)
		ownB ^= (1 << sq)
	}

	for ownR != 0 {
		sq := bitScanForward(ownR)
		otherAttacks |= rookAttacks(Square(sq), occupiedBB, ownPieces)
		ownR ^= (1 << sq)
	}

	for ownQ != 0 {
		sq := bitScanForward(ownQ)
		otherAttacks |= queenAttacks(Square(sq), occupiedBB, ownPieces)
		ownQ ^= (1 << sq)
	}

	return pawnAttacks, minorAttacks, otherAttacks
}

func QueenAttacks(sq Square, occ uint64, own uint64) uint64 {
	return queenAttacks(sq, occ, own)
}

func max(x int16, y int16) int16 {
	if x > y {
		return x
	}
	return y
}

var SquareInnerRingMask [64]uint64 = initInnerRingMask()
var SquareOuterRingMask [64]uint64 = initOuterRingMask()

func initInnerRingMask() [64]uint64 {
	var masks [64]uint64
	for i := 0; i < 64; i++ {
		masks[i] = kingAttacks(uint64(1 << i))
	}
	return masks
}

func initOuterRingMask() [64]uint64 {
	var masks [64]uint64
	for i := 0; i < 64; i++ {
		innerMask := SquareInnerRingMask[i]
		iter := innerMask
		masks[i] = 0
		for iter != 0 {
			sq := bitScanForward(iter)
			masks[i] |= kingAttacks(1 << sq)
			iter ^= (1 << sq)
		}
		masks[i] = masks[i] &^ innerMask
	}
	return masks
}

// Pawn attacks/structure

func (p *Position) CountBackwardPawns(color Color) int16 {
	switch color {
	case White:
		return int16(bits.OnesCount64(wBackward(p.Board.whitePawn, p.Board.blackPawn)))
	case Black:
		return int16(bits.OnesCount64(bBackward(p.Board.blackPawn, p.Board.whitePawn)))
	}
	return 0
}

func (p *Position) CountCandidatePawns(color Color) int16 {
	switch color {
	case White:
		return int16(bits.OnesCount64(wCandidatesOn5th(p.Board.whitePawn, p.Board.blackPawn)))
	case Black:
		return int16(bits.OnesCount64(bCandidatesOn4th(p.Board.blackPawn, p.Board.whitePawn)))
	}
	return 0
}

func (p *Position) CountDoublePawns(color Color) int16 {
	switch color {
	case White:
		return int16(bits.OnesCount64(wPawnsBehindOwn(p.Board.whitePawn)))
	case Black:
		return int16(bits.OnesCount64(bPawnsBehindOwn(p.Board.blackPawn)))
	}
	return 0
}

func (p *Position) CountIsolatedPawns(color Color) int16 {
	switch color {
	case White:
		return int16(bits.OnesCount64(isolanis(p.Board.whitePawn)))
	case Black:
		return int16(bits.OnesCount64(isolanis(p.Board.blackPawn)))
	}
	return 0
}

const lowHalf uint64 = 0x00000000FFFFFFFF
const hiHalf uint64 = 0xFFFFFFFF00000000

func (p *Position) CountPassedPawns(color Color) (int16, int16) {
	switch color {
	case White:
		passers := wPassedPawns(p.Board.whitePawn, p.Board.blackPawn)
		return int16(bits.OnesCount64(passers & lowHalf)), int16(bits.OnesCount64(passers & hiHalf))
	case Black:
		passers := bPassedPawns(p.Board.blackPawn, p.Board.whitePawn)
		return int16(bits.OnesCount64(passers & hiHalf)), int16(bits.OnesCount64(passers & lowHalf))
	}
	return 0, 0
}

// pawn utils

var A_FileFill = FileFill(uint64(1 << A1))
var B_FileFill = FileFill(uint64(1 << B1))
var C_FileFill = FileFill(uint64(1 << C1))
var F_FileFill = FileFill(uint64(1 << F1))
var G_FileFill = FileFill(uint64(1 << G1))
var H_FileFill = FileFill(uint64(1 << H1))

func nortFill(gen uint64) uint64 {
	gen |= (gen << 8)
	gen |= (gen << 16)
	gen |= (gen << 32)
	return gen
}

func soutFill(gen uint64) uint64 {
	gen |= (gen >> 8)
	gen |= (gen >> 16)
	gen |= (gen >> 32)
	return gen
}

func wFrontFill(wpawns uint64) uint64 {
	return nortFill(wpawns)
}

func wRearFill(wpawns uint64) uint64 {
	return soutFill(wpawns)
}

func bFrontFill(bpawns uint64) uint64 {
	return soutFill(bpawns)
}

func bRearFill(bpawns uint64) uint64 {
	return nortFill(bpawns)
}

func fileFill(gen uint64) uint64 {
	return nortFill(gen) | soutFill(gen)
}

func wStop(wpawns uint64) uint64 {
	return nortOne(wpawns)
}

func bStop(bpawns uint64) uint64 {
	return soutOne(bpawns)
}

func wFrontSpan(wpawns uint64) uint64 {
	return nortFill(wStop(wpawns))
}

func bFrontSpan(bpawns uint64) uint64 {
	return soutFill(bStop(bpawns))
}

func wEastAttackFrontSpans(wpawns uint64) uint64 {
	return eastOne(wFrontSpan(wpawns))
}

func wWestAttackFrontSpans(wpawns uint64) uint64 {
	return westOne(wFrontSpan(wpawns))
}

func bEastAttackFrontSpans(bpawns uint64) uint64 {
	return eastOne(bFrontSpan(bpawns))
}

func bWestAttackFrontSpans(bpawns uint64) uint64 {
	return westOne(bFrontSpan(bpawns))
}

func wEastAttackRearSpans(wpawns uint64) uint64 {
	return eastOne(wRearFill(wpawns))
}

func wWestAttackRearSpans(wpawns uint64) uint64 {
	return westOne(wRearFill(wpawns))
}

func bEastAttackRearSpans(bpawns uint64) uint64 {
	return eastOne(bRearFill(bpawns))
}

func bWestAttackRearSpans(bpawns uint64) uint64 {
	return westOne(bRearFill(bpawns))
}

func eastAttackFileFill(pawns uint64) uint64 {
	return eastOne(fileFill(pawns))
}

func westAttackFileFill(pawns uint64) uint64 {
	return westOne(fileFill(pawns))
}

func wBackward(wpawns uint64, bpawns uint64) uint64 {
	stops := wStop(wpawns)
	wAttackSpans := wEastAttackFrontSpans(wpawns) | wWestAttackFrontSpans(wpawns)
	bAttacks := bPawnAnyAttacks(bpawns)
	return (stops & bAttacks &^ wAttackSpans) >> 8
}

func bBackward(bpawns uint64, wpawns uint64) uint64 {
	stops := bStop(bpawns)
	bAttackSpans := bEastAttackFrontSpans(bpawns) | bWestAttackFrontSpans(bpawns)
	wAttacks := wPawnAnyAttacks(wpawns)
	return (stops & wAttacks &^ bAttackSpans) << 8
}

func bCandidatesOn4th(bpawns uint64, wpawns uint64) uint64 {
	wPawnAnyAttacks := wPawnAnyAttacks(wpawns)
	bSafeSquares := bSafePawnSquares(bpawns, wpawns)
	bSafeAttacked := wPawnAnyAttacks & bSafeSquares
	whiteFrontSpan := (wpawns << 8) | (wpawns << 16) // only for 5th rank
	return bpawns & rank4 &^ whiteFrontSpan & (bSafeAttacked << 8)
}

func wCandidatesOn5th(wpawns uint64, bpawns uint64) uint64 {
	bPawnAnyAttacks := bPawnAnyAttacks(bpawns)
	wSafeSquares := wSafePawnSquares(wpawns, bpawns)
	wSafeAttacked := bPawnAnyAttacks & wSafeSquares
	blackFrontSpan := (bpawns >> 8) | (bpawns >> 16) // only for 5th rank
	return wpawns & rank5 &^ blackFrontSpan & (wSafeAttacked >> 8)
}

func wPawnEastAttacks(wpawns uint64) uint64 {
	return noEaOne(wpawns)
}

func wPawnWestAttacks(wpawns uint64) uint64 {
	return noWeOne(wpawns)
}

func bPawnEastAttacks(bpawns uint64) uint64 {
	return soEaOne(bpawns)
}

func bPawnWestAttacks(wpawns uint64) uint64 {
	return soWeOne(wpawns)
}

func wSafePawnSquares(wpawns uint64, bpawns uint64) uint64 {
	wPawnEastAttacks := wPawnEastAttacks(wpawns)
	wPawnWestAttacks := wPawnWestAttacks(wpawns)
	bPawnEastAttacks := bPawnEastAttacks(bpawns)
	bPawnWestAttacks := bPawnWestAttacks(bpawns)
	wPawnDblAttacks := wPawnEastAttacks & wPawnWestAttacks
	wPawnOddAttacks := wPawnEastAttacks ^ wPawnWestAttacks
	bPawnDblAttacks := bPawnEastAttacks & bPawnWestAttacks
	bPawnAnyAttacks := bPawnEastAttacks | bPawnWestAttacks
	return wPawnDblAttacks | ^bPawnAnyAttacks | (wPawnOddAttacks &^ bPawnDblAttacks)
}

func bSafePawnSquares(bpawns uint64, wpawns uint64) uint64 {
	bPawnEastAttacks := bPawnEastAttacks(bpawns)
	bPawnWestAttacks := bPawnWestAttacks(bpawns)
	wPawnEastAttacks := wPawnEastAttacks(wpawns)
	wPawnWestAttacks := wPawnWestAttacks(wpawns)
	bPawnDblAttacks := bPawnEastAttacks & bPawnWestAttacks
	bPawnOddAttacks := bPawnEastAttacks ^ bPawnWestAttacks
	wPawnDblAttacks := wPawnEastAttacks & wPawnWestAttacks
	wPawnAnyAttacks := wPawnEastAttacks | wPawnWestAttacks
	return bPawnDblAttacks | ^wPawnAnyAttacks | (bPawnOddAttacks &^ wPawnDblAttacks)
}

func wFrontSpans(wpawns uint64) uint64 {
	return nortOne(nortFill(wpawns))
}

func bRearSpans(bpawns uint64) uint64 {
	return nortOne(nortFill(bpawns))
}

func bFrontSpans(bpawns uint64) uint64 {
	return soutOne(soutFill(bpawns))
}

func wRearSpans(wpawns uint64) uint64 {
	return soutOne(soutFill(wpawns))
}

// pawns with at least one pawn in front on the same file
func wPawnsBehindOwn(wpawns uint64) uint64 {
	return wpawns & wRearSpans(wpawns)
}

// pawns with at least one pawn behind on the same file
func wPawnsInfrontOwn(wpawns uint64) uint64 {
	return wpawns & wFrontSpans(wpawns)
}

// pawns with at least one pawn in front on the same file
func bPawnsBehindOwn(bpawns uint64) uint64 {
	return bpawns & bRearSpans(bpawns)
}

// pawns with at least one pawn behind on the same file
func bPawnsInfrontOwn(bpawns uint64) uint64 {
	return bpawns & bFrontSpans(bpawns)
}

func noNeighborOnEastFile(pawns uint64) uint64 {
	return pawns &^ westAttackFileFill(pawns)
}

func noNeighborOnWestFile(pawns uint64) uint64 {
	return pawns &^ eastAttackFileFill(pawns)
}

func isolanis(pawns uint64) uint64 {
	return noNeighborOnEastFile(pawns) & noNeighborOnWestFile(pawns)
}

func halfIsolanis(pawns uint64) uint64 {
	return noNeighborOnEastFile(pawns) ^ noNeighborOnWestFile(pawns)
}

func wPassedPawns(wpawns uint64, bpawns uint64) uint64 {
	allFrontSpans := bFrontSpans(bpawns)
	allFrontSpans |= eastOne(allFrontSpans) | westOne(allFrontSpans)
	return wpawns &^ allFrontSpans
}

func bPassedPawns(bpawns uint64, wpawns uint64) uint64 {
	allFrontSpans := wFrontSpans(wpawns)
	allFrontSpans |= eastOne(allFrontSpans) | westOne(allFrontSpans)
	return bpawns &^ allFrontSpans
}

func FileFill(fset uint64) uint64 {
	return nortFill(fset) | soutFill(fset)
}
