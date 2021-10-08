package fathom

// #cgo CFLAGS: -O3 -std=gnu11 -w
// #include "tbprobe.h"
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"math/bits"
	"strings"
	"unsafe"

	. "github.com/amanjpro/zahak/engine"
)

const DefaultProbeDepth = 0

var MaxPieceCount = 0
var MinProbeDepth int8 = DefaultProbeDepth

const (
	TB_LOSS          uint32 = 0 /* LOSS */
	TB_BLESSED_LOSS  uint32 = 1 /* LOSS but 50-move draw */
	TB_DRAW          uint32 = 2 /* DRAW */
	TB_CURSED_WIN    uint32 = 3 /* WIN but 50-move draw  */
	TB_WIN           uint32 = 4 /* WIN  */
	TB_RESULT_FAILED uint32 = 0xFFFFFFFF
)

func SetSyzygyPath(path string) {
	cPath := C.CString(strings.TrimSpace(path))
	defer C.free(unsafe.Pointer(cPath))
	C.tb_init(cPath)
	MaxPieceCount = int(C.TB_LARGEST)
	if MaxPieceCount != 0 {
		fmt.Printf("string info loaded syzygy tablebase %d men\n", MaxPieceCount)
	} else {
		fmt.Print("string info no syzygy tablebase was found\n")
	}
}

func ClearSyzygy() {
	C.tb_free()
}

func ProbeWDL(pos *Position, depth int8) uint32 {

	board := pos.Board
	allPiecesCount := bits.OnesCount64(board.GetWhitePieces() | board.GetBlackPieces())
	if MaxPieceCount == 0 ||
		pos.HalfMoveClock != 0 ||
		pos.HasCastling() ||
		depth < MinProbeDepth ||
		allPiecesCount > MaxPieceCount {
		return TB_RESULT_FAILED
	}

	ep := pos.EnPassant
	if ep == NoSquare {
		ep = A1
	}

	return uint32(C.tb_probe_wdl(
		C.uint64_t(board.GetWhitePieces()),
		C.uint64_t(board.GetBlackPieces()),
		C.uint64_t(board.Kings()),
		C.uint64_t(board.Queens()),
		C.uint64_t(board.Rooks()),
		C.uint64_t(board.Bishops()),
		C.uint64_t(board.Knights()),
		C.uint64_t(board.Pawns()),
		C.uint(pos.HalfMoveClock),
		C.uint(0),
		C.uint(ep),
		C.bool(pos.Turn() == White),
	))
}

var promoPieces = [5]PieceType{NoType, Queen, Rook, Bishop, Knight}

func ProbeDTZ(pos *Position) Move {

	board := pos.Board

	allPiecesCount := bits.OnesCount64(board.GetWhitePieces() | board.GetBlackPieces())
	if pos.HasCastling() || allPiecesCount > MaxPieceCount {
		return EmptyMove
	}

	ep := pos.EnPassant
	if ep == NoSquare {
		ep = A1
	}

	res := uint32(C.tb_probe_root(
		C.uint64_t(board.GetWhitePieces()),
		C.uint64_t(board.GetBlackPieces()),
		C.uint64_t(board.Kings()),
		C.uint64_t(board.Queens()),
		C.uint64_t(board.Rooks()),
		C.uint64_t(board.Bishops()),
		C.uint64_t(board.Knights()),
		C.uint64_t(board.Pawns()),
		C.uint(pos.HalfMoveClock),
		C.uint(0),
		C.uint(ep),
		C.bool(pos.Turn() == White),
		nil,
	))

	if res == TB_RESULT_FAILED {
		return EmptyMove
	}

	src := Square(TB_GET_FROM(res))
	dest := Square(TB_GET_TO(res))
	promo := TB_GET_PROMOTES(res)
	piece := board.PieceAt(src)
	capturedPiece := board.PieceAt(dest)

	return NewMove(src, dest, piece, capturedPiece, promoPieces[promo], 0)
}

// From fathom source, translitrated to Go

const TB_RESULT_WDL_MASK uint32 = 0x0000000F
const TB_RESULT_TO_MASK uint32 = 0x000003F0
const TB_RESULT_FROM_MASK uint32 = 0x0000FC00
const TB_RESULT_PROMOTES_MASK uint32 = 0x00070000
const TB_RESULT_EP_MASK uint32 = 0x00080000
const TB_RESULT_DTZ_MASK uint32 = 0xFFF00000
const TB_RESULT_WDL_SHIFT uint32 = 0
const TB_RESULT_TO_SHIFT uint32 = 4
const TB_RESULT_FROM_SHIFT uint32 = 10
const TB_RESULT_PROMOTES_SHIFT uint32 = 16
const TB_RESULT_EP_SHIFT uint32 = 19
const TB_RESULT_DTZ_SHIFT uint32 = 20

func TB_GET_WDL(res uint32) uint32 {
	return ((res & TB_RESULT_WDL_MASK) >> TB_RESULT_WDL_SHIFT)
}
func TB_GET_TO(res uint32) uint32 {
	return ((res & TB_RESULT_TO_MASK) >> TB_RESULT_TO_SHIFT)
}

func TB_GET_FROM(res uint32) uint32 {
	return ((res & TB_RESULT_FROM_MASK) >> TB_RESULT_FROM_SHIFT)
}

func TB_GET_PROMOTES(res uint32) uint32 {
	return ((res & TB_RESULT_PROMOTES_MASK) >> TB_RESULT_PROMOTES_SHIFT)
}

func TB_GET_EP(res uint32) uint32 {
	return ((res & TB_RESULT_EP_MASK) >> TB_RESULT_EP_SHIFT)
}

func TB_GET_DTZ(res uint32) uint32 {
	return ((res & TB_RESULT_DTZ_MASK) >> TB_RESULT_DTZ_SHIFT)
}

var TB_RESULT_STALEMATE = TB_SET_WDL(0, TB_DRAW)
var TB_RESULT_CHECKMATE = TB_SET_WDL(0, TB_WIN)

func TB_SET_WDL(_res, _wdl uint32) uint32 {
	return (((_res) &^ TB_RESULT_WDL_MASK) | (((_wdl) << TB_RESULT_WDL_SHIFT) & TB_RESULT_WDL_MASK))
}
