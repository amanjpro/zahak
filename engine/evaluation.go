package engine

import "math/bits"

const A_FileFill = uint64(1<<A1 | 1<<A2 | 1<<A3 | 1<<A4 | 1<<A5 | 1<<A6 | 1<<A7 | 1<<A8)
const H_FileFill = uint64(1<<H1 | 1<<H2 | 1<<H3 | 1<<H4 | 1<<H5 | 1<<H6 | 1<<H7 | 1<<H8)
const Rank2Fill = uint64(1<<A2 | 1<<B2 | 1<<C2 | 1<<D2 | 1<<E2 | 1<<F2 | 1<<G2 | 1<<H2)
const Rank7Fill = uint64(1<<A7 | 1<<B7 | 1<<C7 | 1<<D7 | 1<<E7 | 1<<F7 | 1<<G7 | 1<<H7)

func (p *Position) Evaluate() int16 {
	board := p.Board
	// turn := p.Turn()
	var drawDivider int16 = 0

	whitePawnsCount := bits.OnesCount64(board.whitePawn)
	whiteKnightsCount := bits.OnesCount64(board.whiteKnight)
	whiteBishopsCount := bits.OnesCount64(board.whiteBishop)
	whiteRooksCount := bits.OnesCount64(board.whiteRook)
	whiteQueensCount := bits.OnesCount64(board.whiteQueen)

	blackPawnsCount := bits.OnesCount64(board.blackPawn)
	blackKnightsCount := bits.OnesCount64(board.blackKnight)
	blackBishopsCount := bits.OnesCount64(board.blackBishop)
	blackRooksCount := bits.OnesCount64(board.blackRook)
	blackQueensCount := bits.OnesCount64(board.blackQueen)

	allPiecesCount :=
		whitePawnsCount +
			blackPawnsCount +
			whiteKnightsCount +
			blackKnightsCount +
			whiteBishopsCount +
			blackBishopsCount +
			whiteRooksCount +
			blackRooksCount +
			whiteQueensCount +
			blackQueensCount

	// bbWhitePawn := board.whitePawn
	// bbBlackPawn := board.blackPawn
	//
	// // Special endgame scenarios
	// // KPK endgame, ask the bitbase
	// if allPiecesCount == 1 && (whitePawnsCount == 1 || blackPawnsCount == 1) {
	// 	if whitePawnsCount == 1 {
	// 		return KpkProbe(board, White, turn)
	// 	}
	// 	return KpkProbe(board, Black, turn)
	// } else if allPiecesCount == 2 && ((whitePawnsCount == 1 && whiteBishopsCount == 1) || (blackPawnsCount == 1 && blackBishopsCount == 1)) {
	// 	queeningColor := NoColor
	// 	// Is it wrong bishop?
	// 	if whiteBishopsCount == 1 { // White is winning
	// 		bishopColor := Square(bits.TrailingZeros64(board.GetBitboardOf(WhiteBishop))).GetColor()
	// 		if bbWhitePawn&H_FileFill != 0 {
	// 			queeningColor = Black
	// 		} else if bbWhitePawn&A_FileFill != 0 {
	// 			queeningColor = White
	// 		}
	// 		if queeningColor != NoColor && queeningColor != bishopColor {
	// 			res := KpkProbe(board, White, turn)
	// 			if res == 0 {
	// 				return 0
	// 			} else if turn == White {
	// 				return res + BlackBishop.Weight()
	// 			} else {
	// 				return res - BlackBishop.Weight()
	// 			}
	// 		}
	// 	} else { // Black is winning
	// 		bishopColor := Square(bits.TrailingZeros64(board.GetBitboardOf(BlackBishop))).GetColor()
	// 		if bbBlackPawn&H_FileFill != 0 {
	// 			queeningColor = White
	// 		} else if bbBlackPawn&A_FileFill != 0 {
	// 			queeningColor = Black
	// 		}
	// 		if queeningColor != NoColor && queeningColor != bishopColor {
	// 			res := KpkProbe(board, Black, turn)
	// 			if res == 0 {
	// 				return 0
	// 			} else if turn == Black {
	// 				return res + BlackBishop.Weight()
	// 			} else {
	// 				return res - BlackBishop.Weight()
	// 			}
	// 		}
	// 		res := KpkProbe(board, Black, turn)
	// 		if res == 0 {
	// 			return 0
	// 		} else if turn == Black {
	// 			return res + BlackBishop.Weight()
	// 		} else {
	// 			return res - BlackBishop.Weight()
	// 		}
	// 	}
	// }

	// Other Draw scenarios
	{

		if (allPiecesCount == 2 && whiteRooksCount == 1 && (blackKnightsCount == 1 || blackBishopsCount == 1)) ||
			(allPiecesCount == 2 && blackRooksCount == 1 && (whiteKnightsCount == 1 || whiteBishopsCount == 1)) ||
			(allPiecesCount == 2 && (blackKnightsCount == 1 || blackBishopsCount == 1) && whitePawnsCount == 1) ||
			(allPiecesCount == 2 && (whiteKnightsCount == 1 || whiteBishopsCount == 1) && blackPawnsCount == 1) ||
			(allPiecesCount == 3 && blackRooksCount == 1 && whiteRooksCount == 1 && (whiteKnightsCount == 1 || blackKnightsCount == 1 || blackBishopsCount == 1 || whiteBishopsCount == 1)) {
			drawDivider = 3
			// } else if whiteRooksCount+blackRooksCount+whiteKnightsCount+blackKnightsCount+whiteQueensCount+blackQueensCount == 0 &&
			// 	abs16(int16(whitePawnsCount-blackPawnsCount)) <= 1 && whiteBishopsCount == 1 && blackBishopsCount == 1 {
			// 	wColor := Square(bits.TrailingZeros64(board.GetBitboardOf(WhiteBishop))).GetColor()
			// 	bColor := Square(bits.TrailingZeros64(board.GetBitboardOf(BlackBishop))).GetColor()
			// 	// Opposite Color Bishop
			// 	if wColor != bColor {
			// 		promotingWhitePawns := bbWhitePawn & Rank7Fill
			// 		promotingBlackPawns := bbBlackPawn & Rank2Fill
			// 		if promotingBlackPawns == 0 && promotingWhitePawns == 0 {
			// 			drawDivider = 2
			// 		}
			// 	}
		}
	}

	output := p.Net.QuickFeed()
	if p.Turn() == Black {
		return -toEval(output)
	}
	return toEval(output) >> drawDivider

}

func toEval(eval float32) int16 {
	if eval >= MAX_NON_CHECKMATE {
		return int16(MAX_NON_CHECKMATE)
	} else if eval <= MIN_NON_CHECKMATE {
		return int16(MIN_NON_CHECKMATE)
	}
	return int16(eval)
}

func abs16(x int16) int16 {
	if x >= 0 {
		return x
	}
	return -x
}
