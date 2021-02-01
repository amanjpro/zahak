package main

import (
	"fmt"
)

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

func (p *Position) ValidMoves() []Move {
	var allMoves = make([]Move, 0, 256)
	var whiteChecks uint64 = 0
	var blackChecks uint64 = 0

	board := p.board
	allPieces := board.AllPieces()
	blackKingIndex := 0
	whiteKingIndex := 0

	for sq, piece := range allPieces {
		if piece == WhiteRook {
			checks, moves := rookLikeMoves(&board, White, &sq, allMoves)
			whiteChecks |= checks
			allMoves = moves
		} else if piece == BlackRook {
			checks, moves := rookLikeMoves(&board, Black, &sq, allMoves)
			blackChecks |= checks
			allMoves = moves
		} else if piece == WhiteQueen {
			checks1, moves1 := rookLikeMoves(&board, White, &sq, allMoves)
			checks2, moves2 := bishopLikeMoves(&board, White, &sq, moves1)
			whiteChecks |= (checks1 | checks2)
			allMoves = moves2
		} else if piece == BlackQueen {
			checks1, moves1 := rookLikeMoves(&board, Black, &sq, allMoves)
			checks2, moves2 := bishopLikeMoves(&board, Black, &sq, moves1)
			blackChecks |= (checks1 | checks2)
			allMoves = moves2
		} else if piece == WhiteBishop {
			checks, moves := bishopLikeMoves(&board, White, &sq, allMoves)
			whiteChecks |= checks
			allMoves = moves
		} else if piece == BlackBishop {
			checks, moves := bishopLikeMoves(&board, Black, &sq, allMoves)
			blackChecks |= checks
			allMoves = moves
		}
	}

	fmt.Println(blackKingIndex, whiteKingIndex)
	return allMoves
}

func rookLikeMoves(b *Bitboard, color Color, srcSq *Square,
	allMoves []Move) (uint64, []Move) {
	var checkSet uint64 = 0

	src := srcSq.BitboardIndex()

	// vertical up moves
	for index := src; index >= 0; index -= 8 {
		destSq := SquareFromIndex(index)
		piece := b.PieceAt(&destSq)
		if piece == NoPiece {
			allMoves = append(allMoves, Move{srcSq, &destSq, NoType, 0})
			setCheckSet(checkSet, index)
		} else if piece.Color() == color { // Rook & Queen cannot jump
			break
		} else if piece.Type() != King { // This is a capture
			allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), Capture})
			break
		} else if piece.Type() == King {
			// This is a check? what is that? don't add anything yet
			// allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), Capture})
			break
		}
	}

	// vertical down moves
	for index := src; index < 64; index += 8 {
		destSq := SquareFromIndex(index)
		piece := b.PieceAt(&destSq)
		if piece == NoPiece {
			allMoves = append(allMoves, Move{srcSq, &destSq, NoType, 0})
			setCheckSet(checkSet, index)
		} else if piece.Color() == color { // Rook & Queen cannot jump
			break
		} else if piece.Type() != King { // This is a capture
			allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), Capture})
			break
		} else if piece.Type() == King {
			// This is a check? what is that? don't add anything yet
			// allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), Capture})
			break
		}
	}

	// horizontal up moves
	for index := src; index >= 0; index-- {
		destSq := SquareFromIndex(index)
		piece := b.PieceAt(&destSq)
		if piece == NoPiece {
			allMoves = append(allMoves, Move{srcSq, &destSq, NoType, 0})
			setCheckSet(checkSet, index)
		} else if piece.Color() == color { // Rook & Queen cannot jump
			break
		} else if piece.Type() != King { // This is a capture
			allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), Capture})
			break
		} else if piece.Type() == King {
			// This is a check? what is that? don't add anything yet
			// allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), Capture})
			break
		}
	}

	// horizontal down moves
	for index := src; index < 8; index++ {
		destSq := SquareFromIndex(index)
		piece := b.PieceAt(&destSq)
		if piece == NoPiece {
			allMoves = append(allMoves, Move{srcSq, &destSq, NoType, 0})
			setCheckSet(checkSet, index)
		} else if piece.Color() == color { // Rook & Queen cannot jump
			break
		} else if piece.Type() != King { // This is a capture
			allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), Capture})
			break
		} else if piece.Type() == King {
			// This is a check? what is that? don't add anything yet
			// allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), Capture})
			break
		}
	}
	return checkSet, allMoves
}

func bishopLikeMoves(b *Bitboard, color Color, srcSq *Square,
	allMoves []Move) (uint64, []Move) {
	var checkSet uint64 = 0

	src := srcSq.BitboardIndex()

	// up-right
	for index := src; index >= 0; index -= 7 {
		destSq := SquareFromIndex(index)
		piece := b.PieceAt(&destSq)
		if piece == NoPiece {
			allMoves = append(allMoves, Move{srcSq, &destSq, NoType, 0})
			setCheckSet(checkSet, index)
		} else if piece.Color() == color { // Bishop & Queen cannot jump
			break
		} else if piece.Type() != King { // This is a capture
			allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), Capture})
			break
		} else if piece.Type() == King {
			break
			// This is a check? what is that? don't add anything yet
			// allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), Capture})
		}
	}

	// down-left
	for index := src; index < 64; index += 7 {
		destSq := SquareFromIndex(index)
		piece := b.PieceAt(&destSq)
		if piece == NoPiece {
			allMoves = append(allMoves, Move{srcSq, &destSq, NoType, 0})
			setCheckSet(checkSet, index)
		} else if piece.Color() == color { // Bishop & Queen cannot jump
			break
		} else if piece.Type() != King { // This is a capture
			allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), Capture})
			break
		} else if piece.Type() == King {
			break
			// This is a check? what is that? don't add anything yet
			// allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), Capture})
		}
	}

	// up-left
	for index := src; index >= 0; index -= 9 {
		destSq := SquareFromIndex(index)
		piece := b.PieceAt(&destSq)
		if piece == NoPiece {
			allMoves = append(allMoves, Move{srcSq, &destSq, NoType, 0})
			setCheckSet(checkSet, index)
		} else if piece.Color() == color { // Bishop & Queen cannot jump
			break
		} else if piece.Type() != King { // This is a capture
			allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), Capture})
			break
		} else if piece.Type() == King {
			// This is a check? what is that? don't add anything yet
			// allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), Capture})
			break
		}
	}

	// up-right
	for index := src; index < 8; index += 9 {
		destSq := SquareFromIndex(index)
		piece := b.PieceAt(&destSq)
		if piece == NoPiece {
			allMoves = append(allMoves, Move{srcSq, &destSq, NoType, 0})
			setCheckSet(checkSet, index)
		} else if piece.Color() == color { // Bishop & Queen cannot jump
			break
		} else if piece.Type() != King { // This is a capture
			allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), Capture})
			break
		} else if piece.Type() == King {
			// This is a check? what is that? don't add anything yet
			// allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), Capture})
			break
		}
	}
	return checkSet, allMoves
}

func setCheckSet(checkBitSet uint64, bitIndex int8) {
	checkBitSet |= (1 << bitIndex)
}
