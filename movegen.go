package main

import (
	"fmt"
)

/**
Still to support:
- Rook, Bishop and Queen moves and captures
- Check if King is at check (square attacked by opponent)
- Check pinned pieces
- Check double checks (King needs to move)
- Check discovered checks
- Check if promotion results in a check
- Check if draw
- Check if checkmate
- Populated taboo squares for king movements
*/

func (p *Position) LegalMoves() []Move {
	allMoves := make([]Move, 0, 350)
	add := func(m ...Move) {
		allMoves = append(allMoves, m...)
	}

	color := p.Turn()

	if color == White {
		bbPawnMoves(p.board.whitePawn, p.board.whitePieces, p.board.blackPieces,
			p.board.blackKing, color, p.enPassant, add)
		bbKnightMoves(p.board.whiteKnight, p.board.whitePieces, p.board.blackPieces,
			p.board.blackKing, add)
		bbKingMoves(p.board.whiteKing, p.board.whitePieces, p.board.blackPieces,
			0, p.HasTag(WhiteCanCastleKingSide), p.HasTag(WhiteCanCastleQueenSide), add)
	} else if color == Black {
		bbPawnMoves(p.board.blackPawn, p.board.blackPieces, p.board.whitePieces,
			p.board.whiteKing, color, p.enPassant, add)
		bbKnightMoves(p.board.blackKnight, p.board.blackPieces, p.board.whitePieces,
			p.board.whiteKing, add)
		bbKingMoves(p.board.blackKing, p.board.blackPieces, p.board.whitePieces,
			0, p.HasTag(BlackCanCastleKingSide), p.HasTag(BlackCanCastleQueenSide), add)
	}
	fmt.Println("LEN", len(allMoves))
	return allMoves
}

// Pawns

func bbPawnMoves(bbPawn uint64, ownPieces uint64, otherPieces uint64, otherKing uint64,
	color Color, enPassant Square, add func(m ...Move)) {
	both := (otherPieces | ownPieces) ^ universal
	if color == White {
		for bbPawn != 0 {
			src := bitScanReverse(bbPawn)
			srcSq := Square(src)
			pawn := uint64(1 << src)
			dbl := wDblPushTargets(pawn, both)
			if dbl != 0 {
				dest := Square(bitScanReverse(dbl))
				var tag MoveTag = 0
				if wPawnsAble2CaptureAny(dbl, otherKing) != 0 {
					tag |= Check
				}
				add(Move{srcSq, dest, NoType, tag})
			}
			sngl := wSinglePushTargets(pawn, both)
			if sngl != 0 {
				dest := Square(bitScanReverse(sngl))
				if dest.Rank() == Rank8 {
					// TODO: set check flags
					add(Move{srcSq, dest, Queen, 0},
						Move{srcSq, dest, Rook, 0},
						Move{srcSq, dest, Bishop, 0},
						Move{srcSq, dest, Knight, 0})
				} else {
					var tag MoveTag = 0
					if wPawnsAble2CaptureAny(sngl, otherKing) != 0 {
						tag |= Check
					}
					add(Move{srcSq, dest, NoType, tag})
				}
			}
			for _, sq := range getIndicesOfOnes(wPawnsAble2CaptureAny(pawn, otherPieces)) {
				dest := Square(sq)
				if dest.Rank() == Rank8 {
					// TODO: set check flags
					add(Move{srcSq, dest, Queen, Capture},
						Move{srcSq, dest, Rook, Capture},
						Move{srcSq, dest, Bishop, Capture},
						Move{srcSq, dest, Knight, Capture})
				} else {
					var tag MoveTag = Capture
					if wPawnsAble2CaptureAny(uint64(1<<sq), otherKing) != 0 {
						tag |= Check
					}
					add(Move{srcSq, dest, NoType, tag})
				}
			}
			if srcSq.Rank() == Rank5 && enPassant != NoSquare && enPassant.Rank() == Rank6 {
				ep := uint64(1 << enPassant)
				r := wPawnsAble2CaptureAny(pawn, ep)
				if r != 0 {
					dest := Square(bitScanReverse(r))
					var tag MoveTag = Capture | EnPassant
					if wPawnsAble2CaptureAny(uint64(1<<r), otherKing) != 0 {
						tag |= Check
					}
					add(Move{srcSq, dest, NoType, tag})
				}
			}
			bbPawn ^= pawn
		}
	} else if color == Black {
		for bbPawn != 0 {
			src := bitScanReverse(bbPawn)
			srcSq := Square(src)
			pawn := uint64(1 << src)
			dbl := bDoublePushTargets(pawn, both)
			if dbl != 0 {
				dest := Square(bitScanReverse(dbl))
				var tag MoveTag = Capture
				if bPawnsAble2CaptureAny(uint64(1<<dbl), otherKing) != 0 {
					tag |= Check
				}
				add(Move{srcSq, dest, NoType, tag})
			}
			sngl := bSinglePushTargets(pawn, both)
			if sngl != 0 {
				dest := Square(bitScanReverse(sngl))
				if dest.Rank() == Rank8 {
					// TODO: set check flags
					add(Move{srcSq, dest, Queen, 0},
						Move{srcSq, dest, Rook, 0},
						Move{srcSq, dest, Bishop, 0},
						Move{srcSq, dest, Knight, 0})
				} else {
					tag := Capture
					if bPawnsAble2CaptureAny(sngl, otherKing) != 0 {
						tag |= Check
					}
					add(Move{srcSq, dest, NoType, tag})
				}
			}
			for _, sq := range getIndicesOfOnes(bPawnsAble2CaptureAny(pawn, otherPieces)) {
				dest := Square(sq)
				if dest.Rank() == Rank1 {
					// TODO: set check flags
					add(Move{srcSq, dest, Queen, Capture},
						Move{srcSq, dest, Rook, Capture},
						Move{srcSq, dest, Bishop, Capture},
						Move{srcSq, dest, Knight, Capture})
				} else {
					var tag MoveTag = Capture
					if bPawnsAble2CaptureAny(uint64(1<<sq), otherKing) != 0 {
						tag |= Check
					}
					add(Move{srcSq, dest, NoType, tag})
				}
			}
			if srcSq.Rank() == Rank4 && enPassant != NoSquare && enPassant.Rank() == Rank3 {
				ep := uint64(1 << enPassant)
				r := bPawnsAble2CaptureAny(pawn, ep)
				if r != 0 {
					dest := Square(bitScanReverse(r))
					var tag MoveTag = Capture | EnPassant
					if bPawnsAble2CaptureAny(r, otherKing) != 0 {
						tag |= Check
					}
					add(Move{srcSq, dest, NoType, tag})
				}
			}
			bbPawn ^= pawn
		}
	}
}


// Knights
func bbKnightMoves(bbPiece uint64, ownPieces uint64, otherPieces uint64,
	otherKing uint64, add func(m ...Move)) {
	both := (otherPieces | ownPieces)
	for bbPiece != 0 {
		src := bitScanReverse(bbPiece)
		srcSq := Square(src)
		knight := uint64(1 << src)
		moves := knightMovesNoCaptures(knight, both)
		fmt.Println("MOVES: ", moves)
		for moves != 0 {
			sq := bitScanReverse(moves)
			dest := Square(sq)
			var tag MoveTag = 0
			if knightCaptures(knight, otherKing) != 0 {
				tag |= Check
			}
			add(Move{srcSq, dest, NoType, tag})
			moves ^= (1 << sq)
		}
		captures := knightCaptures(knight, otherPieces)
		fmt.Println("CAPTURE: ", captures)
		for captures != 0 {
			sq := bitScanReverse(captures)
			dest := Square(sq)
			var tag MoveTag = Capture
			if knightCaptures(knight, otherKing) != 0 {
				tag |= Check
			}
			add(Move{srcSq, dest, NoType, tag})
			captures ^= (1 << sq)
		}
		bbPiece ^= knight
	}
}

// Kings
func bbKingMoves(bbPiece uint64, ownPieces uint64, otherPieces uint64,
	tabooSquares uint64, kingSideCastle bool, queenSideCastle bool,
	add func(m ...Move)) {
	both := (otherPieces | ownPieces)
	if bbPiece != 0 {
		src := bitScanReverse(bbPiece)
		srcSq := Square(src)
		king := uint64(1 << src)
		moves := kingMovesNoCaptures(king, both, tabooSquares)
		for moves != 0 {
			sq := bitScanReverse(moves)
			dest := Square(sq)
			add(Move{srcSq, dest, NoType, 0})
			moves ^= (1 << sq)
		}
		captures := kingCaptures(king, otherPieces, tabooSquares)
		for captures != 0 {
			sq := bitScanReverse(captures)
			dest := Square(sq)
			add(Move{srcSq, dest, NoType, Capture})
			captures ^= (1 << sq)
		}

		E := E1
		F := F1
		G := G1
		D := D1
		C := C1
		B := B1
		if srcSq.Rank() == Rank8 {
			E = E8
			F = F8
			G = G8
			D = D8
			C = C8
			B = B8
		}

		kingSide := uint64(1<<F | 1<<G)
		queenSide := uint64(1<<D | 1<<C | 1<<B)

		if kingSideCastle &&
			(ownPieces&kingSide == 0) && // are empty
			(tabooSquares&(kingSide|1<<E) == 0) { // Not in check
			add(Move{srcSq, G, NoType, 0})
		}

		if queenSideCastle &&
			(ownPieces&queenSide == 0) && // are empty
			(tabooSquares&(queenSide|1<<E) == 0) { // Not in check
			add(Move{srcSq, C, NoType, 0})
		}
	}
}


// Pawn Pushes

func wSinglePushTargets(wpawns uint64, empty uint64) uint64 {
	return nortOne(wpawns) & empty
}

func wDblPushTargets(wpawns uint64, empty uint64) uint64 {
	const rank4 = uint64(0x00000000FF000000)
	singlePushs := wSinglePushTargets(wpawns, empty)
	return nortOne(singlePushs) & empty & rank4
}

func bSinglePushTargets(bpawns uint64, empty uint64) uint64 {
	return soutOne(bpawns) & empty
}

func bDoublePushTargets(bpawns uint64, empty uint64) uint64 {
	const rank5 = uint64(0x000000FF00000000)
	singlePushs := bSinglePushTargets(bpawns, empty)
	return soutOne(singlePushs) & empty & rank5
}

// Pawns Attacks
// func wPawnEastAttacks(wpawns uint64) uint64 {
// 	return noEaOne(wpawns)
// }
//
// func wPawnWestAttacks(wpawns uint64) uint64 {
// 	return noWeOne(wpawns)
// }
//
// func bPawnEastAttacks(bpawns uint64) uint64 {
// 	return soEaOne(bpawns)
// }
//
// func bPawnWestAttacks(bpawns uint64) uint64 {
// 	return soWeOne(bpawns)
// }

func wPawnAnyAttacks(wpawns uint64) uint64 {
	return noEaOne(wpawns) | noWeOne(wpawns)
}

// func wPawnDblAttacks(wpawns uint64) uint64 {
// 	return noEaOne(wpawns) & noWeOne(wpawns)
// }

// func wPawnSingleAttacks(wpawns uint64) uint64 {
// 	return noEaOne(wpawns) ^ noWeOne(wpawns)
// }

func bPawnAnyAttacks(bpawns uint64) uint64 {
	return soEaOne(bpawns) | soWeOne(bpawns)
}

// func bPawnDblAttacks(bpawns uint64) uint64 {
// 	return soEaOne(bpawns) & soWeOne(bpawns)
// }

// func bPawnSingleAttacks(bpawns uint64) uint64 {
// 	return soEaOne(bpawns) ^ soWeOne(bpawns)
// }

// func wPawnsAble2CaptureEast(wpawns uint64, bpieces uint64) uint64 {
// 	return wpawns & bPawnWestAttacks(bpieces)
// }
//
// func wPawnsAble2CaptureWest(wpawns uint64, bpieces uint64) uint64 {
// 	return wpawns & bPawnEastAttacks(bpieces)
// }

func wPawnsAble2CaptureAny(wpawns uint64, bpieces uint64) uint64 {
	return wPawnAnyAttacks(wpawns) & bpieces
}

// func bPawnsAble2CaptureEast(bpawns uint64, wpieces uint64) uint64 {
// 	return bpawns & bPawnWestAttacks(wpieces)
// }
//
// func bPawnsAble2CaptureWest(bpawns uint64, wpieces uint64) uint64 {
// 	return bpawns & bPawnEastAttacks(wpieces)
// }
//
func bPawnsAble2CaptureAny(bpawns uint64, wpieces uint64) uint64 {
	return bPawnAnyAttacks(bpawns) & wpieces
}

// The mighty knight

func knightMovesNoCaptures(b uint64, other uint64) uint64 {
	attacks := knightAttacks(b)
	return (attacks ^ other) & attacks
}

func knightCaptures(b uint64, other uint64) uint64 {
	return knightAttacks(b) & other
}

func knightAttacks(b uint64) uint64 {
	return noNoEa(b) | noEaEa(b) | soEaEa(b) | soSoEa(b) |
		noNoWe(b) | noWeWe(b) | soWeWe(b) | soSoWe(b)
}

func noNoEa(b uint64) uint64 {
	return (b << 17) & notAFile
}
func noEaEa(b uint64) uint64 {
	return (b << 10) & notABFile
}
func soEaEa(b uint64) uint64 {
	return (b >> 6) & notABFile
}
func soSoEa(b uint64) uint64 {
	return (b >> 15) & notAFile
}
func noNoWe(b uint64) uint64 {
	return (b << 15) & notHFile
}
func noWeWe(b uint64) uint64 {
	return (b << 6) & notGHFile
}
func soWeWe(b uint64) uint64 {
	return (b >> 10) & notGHFile
}
func soSoWe(b uint64) uint64 {
	return (b >> 17) & notHFile
}

// Utilites
func getIndicesOfOnes(bb uint64) []uint8 {
	indices := make([]uint8, 0, 8)
	for bb != 0 {
		index := bitScanReverse(bb)
		bb ^= (1 << index)
		indices = append(indices, index)
	}
	return indices
}

var index64 [64]uint8 = [64]uint8{
	0, 47, 1, 56, 48, 27, 2, 60,
	57, 49, 41, 37, 28, 16, 3, 61,
	54, 58, 35, 52, 50, 42, 21, 44,
	38, 32, 29, 23, 17, 11, 4, 62,
	46, 55, 26, 59, 40, 36, 15, 53,
	34, 51, 20, 43, 31, 22, 10, 45,
	25, 39, 14, 33, 19, 30, 9, 24,
	13, 18, 8, 12, 7, 6, 5, 63,
}

/**
 * bitScanReverse
 * @authors Kim Walisch, Mark Dickinson
 * @param bb bitboard to scan
 * @precondition bb != 0
 * @return index (0..63) of most significant one bit
 */
func bitScanReverse(bb uint64) uint8 {
	const debruijn64 = uint64(0x03f79d71b4cb0a89)
	bb |= bb >> 1
	bb |= bb >> 2
	bb |= bb >> 4
	bb |= bb >> 8
	bb |= bb >> 16
	bb |= bb >> 32
	return index64[(bb*debruijn64)>>58]
}

// func diagonal(file File, rank Rank) uint64 {
// 	return (uint64(rank) - uint64(file)) & 15
// }
//
// func antiDiagonal(file File, rank Rank) uint64 {
// 	return (uint64(file) + uint64(rank)) ^ 15
// }

// King & Kingslayer
func kingMovesNoCaptures(b uint64, others uint64, tabooSquares uint64) uint64 {
	attacks := kingAttacks(b)
	return (attacks ^ others) & attacks & (tabooSquares ^ universal)
}

func kingCaptures(b uint64, others uint64, tabooSquares uint64) uint64 {
	return kingAttacks(b) & others & (tabooSquares ^ universal)
}

func kingAttacks(b uint64) uint64 {
	return soutOne(b) | nortOne(b) | eastOne(b) | noEaOne(b) |
		soEaOne(b) | westOne(b) | soWeOne(b) | noWeOne(b)
}

func soutOne(b uint64) uint64 {
	return b >> 8
}
func nortOne(b uint64) uint64 {
	return b << 8
}

func eastOne(b uint64) uint64 {
	return (b << 1) & notAFile
}
func noEaOne(b uint64) uint64 {
	return (b << 9) & notAFile
}
func soEaOne(b uint64) uint64 {
	return (b >> 7) & notAFile
}
func westOne(b uint64) uint64 {
	return (b >> 1) & notHFile
}
func soWeOne(b uint64) uint64 {
	return (b >> 9) & notHFile
}
func noWeOne(b uint64) uint64 {
	return (b << 7) & notHFile
}

const empty = uint64(0)
const universal = uint64(0xffffffffffffffff)
const notAFile = uint64(0xfefefefefefefefe) // ~0x0101010101010101
const notBFile = uint64(0xbfbfbfbfbfbfbfbf)
const notGFile = uint64(0xfdfdfdfdfdfdfdfd)
const notHFile = uint64(0x7f7f7f7f7f7f7f7f) // ~0x8080808080808080
const notABFile = notAFile & notBFile
const notGHFile = notGFile & notHFile
