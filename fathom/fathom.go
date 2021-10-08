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
	if MaxPieceCount == 0 &&
		pos.HalfMoveClock != 0 &&
		pos.HasCastling() &&
		!depthCardinalityCheck(pos, depth) {
		return TB_RESULT_FAILED
	}

	board := pos.Board
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
		C.uint(pos.EnPassant),
		C.bool(pos.Turn() == White),
	))
}

func depthCardinalityCheck(pos *Position, depth int8) bool {
	board := pos.Board
	cardinality := bits.OnesCount64(board.GetWhitePieces() | board.GetBlackPieces())
	return cardinality < MaxPieceCount || (cardinality == MaxPieceCount && depth >= MinProbeDepth)
}

var promoPieces = [5]PieceType{NoType, Queen, Rook, Bishop, Knight}

func ProbeDTZ(pos *Position) Move {

	board := pos.Board

	cardinality := bits.OnesCount64(board.GetWhitePieces() | board.GetBlackPieces())
	if pos.HasCastling() && cardinality > MaxPieceCount {
		return EmptyMove
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
		C.uint(pos.EnPassant),
		C.bool(pos.Turn() == White),
		nil,
	))

	if res == TB_RESULT_FAILED || res == TB_RESULT_STALEMATE || res == TB_RESULT_CHECKMATE {
		return EmptyMove
	}

	start := Square(TB_GET_FROM(res))
	end := Square(TB_GET_TO(res))
	ep := Square(TB_GET_EP(res))
	promo := TB_GET_PROMOTES(res)
	piece := board.PieceAt(start)
	capturedPiece := board.PieceAt(end)
	var tags MoveTag = 0
	if capturedPiece != NoPiece {
		tags |= Capture
	}

	if ep != 0 {
		tags |= EnPassant
	}
	return NewMove(start, end, piece, capturedPiece, promoPieces[promo], tags)
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
