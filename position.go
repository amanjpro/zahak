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
	var whiteChecksRook uint64 = 0
	var whiteChecksBishop uint64 = 0
	var whiteChecksQueen uint64 = 0
	var whiteChecksKnight uint64 = 0
	var whiteChecksPawn uint64 = 0
	var whiteChecksKing uint64 = 0
	var blackChecksRook uint64 = 0
	var blackChecksBishop uint64 = 0
	var blackChecksQueen uint64 = 0
	var blackChecksKnight uint64 = 0
	var blackChecksPawn uint64 = 0
	var blackChecksKing uint64 = 0

	board := p.board
	allPieces := board.AllPieces()
	var blackKingIndex int8 = 0
	var whiteKingIndex int8 = 0

	for sq, piece := range allPieces {
		if piece == WhiteRook {
			checks, moves := rookLikeMoves(&board, WhiteRook, &sq, allMoves)
			whiteChecksRook |= checks
			allMoves = moves
		} else if piece == BlackRook {
			checks, moves := rookLikeMoves(&board, BlackRook, &sq, allMoves)
			blackChecksRook |= checks
			allMoves = moves
		} else if piece == WhiteQueen {
			checks1, moves1 := rookLikeMoves(&board, WhiteQueen, &sq, allMoves)
			checks2, moves2 := bishopLikeMoves(&board, WhiteQueen, &sq, moves1)
			whiteChecksQueen |= (checks1 | checks2)
			allMoves = moves2
		} else if piece == BlackQueen {
			checks1, moves1 := rookLikeMoves(&board, BlackQueen, &sq, allMoves)
			checks2, moves2 := bishopLikeMoves(&board, BlackQueen, &sq, moves1)
			blackChecksQueen |= (checks1 | checks2)
			allMoves = moves2
		} else if piece == WhiteBishop {
			checks, moves := bishopLikeMoves(&board, WhiteBishop, &sq, allMoves)
			whiteChecksBishop |= checks
			allMoves = moves
		} else if piece == BlackBishop {
			checks, moves := bishopLikeMoves(&board, BlackBishop, &sq, allMoves)
			blackChecksBishop |= checks
			allMoves = moves
		} else if piece == WhiteKnight {
			checks, moves := knightMoves(&board, White, &sq, allMoves)
			whiteChecksKnight |= checks
			allMoves = moves
		} else if piece == BlackKnight {
			checks, moves := knightMoves(&board, Black, &sq, allMoves)
			blackChecksKnight |= checks
			allMoves = moves
		} else if piece == WhiteKing {
			whiteKingIndex = sq.BitboardIndex()
		} else if piece == BlackKing {
			blackKingIndex = sq.BitboardIndex()
		}
	}

	fmt.Println(blackKingIndex, whiteKingIndex)
	return allMoves
}

func rookLikeMoves(b *Bitboard, p Piece, srcSq *Square,
	allMoves []Move) (uint64, []Move) {
	color := p.Color()
	var checkSet uint64 = 0

	src := srcSq.BitboardIndex()

	// vertical up moves
	for index := src; index >= 0; index -= 8 {
		destSq := SquareFromIndex(index)
		piece := b.PieceAt(&destSq)
		if piece == NoPiece {
			tag := rookLikeCheckTag(b, color, &destSq)
			if tag != Check && p.Type() == Queen {
				tag = bishopLikeCheckTag(b, color, &destSq)
			}
			allMoves = append(allMoves, Move{srcSq, &destSq, NoType, tag})
			setCheckSet(checkSet, index)
		} else if piece.Color() == color { // Rook & Queen cannot jump
			break
		} else if piece.Type() != King { // This is a capture
			tag := rookLikeCheckTag(b, color, &destSq)
			if tag != Check && p.Type() == Queen {
				tag = bishopLikeCheckTag(b, color, &destSq)
			}
			allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), Capture | tag})
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
			tag := rookLikeCheckTag(b, color, &destSq)
			if tag != Check && p.Type() == Queen {
				tag = bishopLikeCheckTag(b, color, &destSq)
			}
			allMoves = append(allMoves, Move{srcSq, &destSq, NoType, tag})
			setCheckSet(checkSet, index)
		} else if piece.Color() == color { // Rook & Queen cannot jump
			break
		} else if piece.Type() != King { // This is a capture
			tag := rookLikeCheckTag(b, color, &destSq)
			if tag != Check && p.Type() == Queen {
				tag = bishopLikeCheckTag(b, color, &destSq)
			}
			allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), tag | Capture})
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
			tag := rookLikeCheckTag(b, color, &destSq)
			if tag != Check && p.Type() == Queen {
				tag = bishopLikeCheckTag(b, color, &destSq)
			}
			allMoves = append(allMoves, Move{srcSq, &destSq, NoType, tag})
			setCheckSet(checkSet, index)
		} else if piece.Color() == color { // Rook & Queen cannot jump
			break
		} else if piece.Type() != King { // This is a capture
			tag := rookLikeCheckTag(b, color, &destSq)
			if tag != Check && p.Type() == Queen {
				tag = bishopLikeCheckTag(b, color, &destSq)
			}
			allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), Capture | tag})
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
			tag := rookLikeCheckTag(b, color, &destSq)
			if tag != Check && p.Type() == Queen {
				tag = bishopLikeCheckTag(b, color, &destSq)
			}
			allMoves = append(allMoves, Move{srcSq, &destSq, NoType, tag})
			setCheckSet(checkSet, index)
		} else if piece.Color() == color { // Rook & Queen cannot jump
			break
		} else if piece.Type() != King { // This is a capture
			tag := rookLikeCheckTag(b, color, &destSq)
			if tag != Check && p.Type() == Queen {
				tag = bishopLikeCheckTag(b, color, &destSq)
			}
			allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), Capture | tag})
			break
		} else if piece.Type() == King {
			// This is a check? what is that? don't add anything yet
			// allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), Capture})
			break
		}
	}
	return checkSet, allMoves
}

func bishopLikeMoves(b *Bitboard, p Piece, srcSq *Square,
	allMoves []Move) (uint64, []Move) {
	color := p.Color()
	var checkSet uint64 = 0

	src := srcSq.BitboardIndex()

	// up-right
	for index := src; index >= 0; index -= 7 {
		destSq := SquareFromIndex(index)
		piece := b.PieceAt(&destSq)
		if piece == NoPiece {
			tag := bishopLikeCheckTag(b, color, &destSq)
			if tag != Check && p.Type() == Queen {
				tag = rookLikeCheckTag(b, color, &destSq)
			}
			allMoves = append(allMoves, Move{srcSq, &destSq, NoType, tag})
			setCheckSet(checkSet, index)
		} else if piece.Color() == color { // Bishop & Queen cannot jump
			break
		} else if piece.Type() != King { // This is a capture
			tag := bishopLikeCheckTag(b, color, &destSq)
			if tag != Check && p.Type() == Queen {
				tag = rookLikeCheckTag(b, color, &destSq)
			}
			allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), Capture | tag})
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
			tag := bishopLikeCheckTag(b, color, &destSq)
			if tag != Check && p.Type() == Queen {
				tag = rookLikeCheckTag(b, color, &destSq)
			}
			allMoves = append(allMoves, Move{srcSq, &destSq, NoType, tag})
			setCheckSet(checkSet, index)
		} else if piece.Color() == color { // Bishop & Queen cannot jump
			break
		} else if piece.Type() != King { // This is a capture
			tag := bishopLikeCheckTag(b, color, &destSq)
			if tag != Check && p.Type() == Queen {
				tag = rookLikeCheckTag(b, color, &destSq)
			}
			allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), tag | Capture})
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
			tag := bishopLikeCheckTag(b, color, &destSq)
			if tag != Check && p.Type() == Queen {
				tag = rookLikeCheckTag(b, color, &destSq)
			}
			allMoves = append(allMoves, Move{srcSq, &destSq, NoType, tag})
			setCheckSet(checkSet, index)
		} else if piece.Color() == color { // Bishop & Queen cannot jump
			break
		} else if piece.Type() != King { // This is a capture
			tag := bishopLikeCheckTag(b, color, &destSq)
			if tag != Check && p.Type() == Queen {
				tag = rookLikeCheckTag(b, color, &destSq)
			}
			allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), tag | Capture})
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
			tag := bishopLikeCheckTag(b, color, &destSq)
			if tag != Check && p.Type() == Queen {
				tag = rookLikeCheckTag(b, color, &destSq)
			}
			allMoves = append(allMoves, Move{srcSq, &destSq, NoType, tag})
			setCheckSet(checkSet, index)
		} else if piece.Color() == color { // Bishop & Queen cannot jump
			break
		} else if piece.Type() != King { // This is a capture
			tag := bishopLikeCheckTag(b, color, &destSq)
			if tag != Check && p.Type() == Queen {
				tag = rookLikeCheckTag(b, color, &destSq)
			}
			allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), Capture | tag})
			break
		} else if piece.Type() == King {
			// This is a check? what is that? don't add anything yet
			// allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), Capture})
			break
		}
	}
	return checkSet, allMoves
}

func knightMoves(b *Bitboard, color Color, srcSq *Square,
	allMoves []Move) (uint64, []Move) {
	src := srcSq.BitboardIndex()
	var checkSet uint64 = 0
	destinations := []int8{src + 15, src + 17, src + 6, src + 10, src - 15, src - 17, src - 6, src - 10}

	for i := 0; i < len(destinations); i++ {
		sq := destinations[i]
		if sq >= 0 && sq < 64 {
			destSq := SquareFromIndex(sq)
			piece := b.PieceAt(&destSq)
			if piece == NoPiece {
				tag := knightCheckTag(b, color, &destSq)
				allMoves = append(allMoves, Move{srcSq, &destSq, NoType, tag})
				setCheckSet(checkSet, sq)
			} else if piece.Color() == color { // Knight cannot land on own pieces
				continue
			} else if piece.Type() != King { // This is a capture
				tag := knightCheckTag(b, color, &destSq)
				allMoves = append(allMoves, Move{srcSq, &destSq, piece.Type(), tag | Capture})
			}
		}
	}

	return checkSet, allMoves
}

func knightCheckTag(b *Bitboard, color Color, srcSq *Square) MoveTag {
	src := srcSq.BitboardIndex()
	destinations := []int8{src + 15, src + 17, src + 6, src + 10, src - 15, src - 17, src - 6, src - 10}

	for i := 0; i < len(destinations); i++ {
		sq := destinations[i]
		if sq >= 0 && sq < 64 {
			destSq := SquareFromIndex(sq)
			piece := b.PieceAt(&destSq)
			if piece.Color() != color && piece.Type() == King {
				return Check
			}
		}
	}

	return 0
}

func rookLikeCheckTag(b *Bitboard, color Color, srcSq *Square) MoveTag {
	src := srcSq.BitboardIndex()

	// vertical up moves
	for index := src; index >= 0; index -= 8 {
		destSq := SquareFromIndex(index)
		piece := b.PieceAt(&destSq)
		if piece.Color() != color && piece.Type() == King {
			return Check
		}
	}

	// vertical down moves
	for index := src; index < 64; index += 8 {
		destSq := SquareFromIndex(index)
		piece := b.PieceAt(&destSq)
		if piece.Color() != color && piece.Type() == King {
			return Check
		}
	}

	// horizontal up moves
	for index := src; index >= 0; index-- {
		destSq := SquareFromIndex(index)
		piece := b.PieceAt(&destSq)
		if piece.Color() != color && piece.Type() == King {
			return Check
		}
	}

	// horizontal down moves
	for index := src; index < 8; index++ {
		destSq := SquareFromIndex(index)
		piece := b.PieceAt(&destSq)
		if piece.Color() != color && piece.Type() == King {
			return Check
		}
	}
	return 0
}

func bishopLikeCheckTag(b *Bitboard, color Color, srcSq *Square) MoveTag {

	src := srcSq.BitboardIndex()

	// up-right
	for index := src; index >= 0; index -= 7 {
		destSq := SquareFromIndex(index)
		piece := b.PieceAt(&destSq)
		if piece.Color() != color && piece.Type() == King {
			return Check
		}
	}

	// down-left
	for index := src; index < 64; index += 7 {
		destSq := SquareFromIndex(index)
		piece := b.PieceAt(&destSq)
		if piece.Color() != color && piece.Type() == King {
			return Check
		}
	}

	// up-left
	for index := src; index >= 0; index -= 9 {
		destSq := SquareFromIndex(index)
		piece := b.PieceAt(&destSq)
		if piece.Color() != color && piece.Type() == King {
			return Check
		}
	}

	// up-right
	for index := src; index < 8; index += 9 {
		destSq := SquareFromIndex(index)
		piece := b.PieceAt(&destSq)
		if piece.Color() != color && piece.Type() == King {
			return Check
		}
	}
	return 0
}

func setCheckSet(checkBitSet uint64, bitIndex int8) {
	checkBitSet |= (1 << bitIndex)
}
