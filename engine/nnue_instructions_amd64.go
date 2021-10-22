//go:build amd64
// +build amd64

package engine

import "unsafe"

//go:noescape
func _update_hidden(white_previous_outputs, black_previous_outputs, white_update_indices, white_update_coeffs, black_update_indices, black_update_coeffs, update_size, weights, white_outputs, black_outputs, outputs_len unsafe.Pointer)

//go:noescape
func _quick_feed(hidden_outputs, hidden_outputs_len, weights, weights_len, res unsafe.Pointer)

func (n *NetworkState) UpdateHidden(wUpdates *Updates, bUpdates *Updates) {
	n.CurrentHidden += 1

	p1 := unsafe.Pointer(&n.WhiteHiddenOutputs[n.CurrentHidden-1][0])
	p2 := unsafe.Pointer(&n.BlackHiddenOutputs[n.CurrentHidden-1][0])
	p3 := unsafe.Pointer(&wUpdates.Indices[0])
	p4 := unsafe.Pointer(&wUpdates.Coeffs[0])
	p5 := unsafe.Pointer(&bUpdates.Indices[0])
	p6 := unsafe.Pointer(&bUpdates.Coeffs[0])
	p7 := unsafe.Pointer(uintptr(bUpdates.Size))
	p8 := unsafe.Pointer(&n.HiddenWeights[0])
	p9 := unsafe.Pointer(&n.WhiteHiddenOutputs[n.CurrentHidden][0])
	p10 := unsafe.Pointer(&n.BlackHiddenOutputs[n.CurrentHidden][0])
	p11 := unsafe.Pointer(uintptr(NetHiddenSize))

	_update_hidden(p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11)
}

func (n *NetworkState) QuickFeed(turn Color) float32 {
	var p1 unsafe.Pointer
	if turn == White {
		p1 = unsafe.Pointer(&n.WhiteHiddenOutputs[n.CurrentHidden][0])
	} else {
		p1 = unsafe.Pointer(&n.BlackHiddenOutputs[n.CurrentHidden][0])
	}
	p2 := unsafe.Pointer(uintptr(NetHiddenSize))
	p3 := unsafe.Pointer(&n.OutputWeights[0])
	p4 := unsafe.Pointer(uintptr(NetHiddenSize))
	var res float32

	_quick_feed(p1, p2, p3, p4, unsafe.Pointer(&res))
	return res + n.OutputBias
}
