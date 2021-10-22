//go:build !amd64
// +build !amd64

package engine

func (n *NetworkState) QuickFeed(turn Color) float32 {
	// apply output layer
	output := float32(0)
	var hiddenOutputs [][]float32
	if turn == White {
		hiddenOutputs = n.WhiteHiddenOutputs[n.CurrentHidden]
	} else {
		hiddenOutputs = n.BlackHiddenOutputs[n.CurrentHidden]
	}
	for i := 0; i < len(n.OutputWeights); i++ {
		output += ReLu(hiddenOutputs[i]) * n.OutputWeights[i]
	}
	return output + n.OutputBias
}

func (n *NetworkState) UpdateHidden(wUpdates *Updates, bUpdates *Updates) {
	n.CurrentHidden += 1
	wHiddenOutputs := n.WhiteHiddenOutputs[n.CurrentHidden]
	bHiddenOutputs := n.BlackHiddenOutputs[n.CurrentHidden]
	copy(wHiddenOutputs, n.WhiteHiddenOutputs[n.CurrentHidden-1])
	copy(bHiddenOutputs, n.BlackHiddenOutputs[n.CurrentHidden-1])

	for i := 0; i < updates.Size; i++ {
		weights := n.HiddenWeights
		for j := 0; j < len(hiddenOutputs); j++ {
			wHiddenOutputs[j] += float32(wUpdates.Coeffs[i]) * weights[int(wUpdates.Indices[i])*NetHiddenSize+j]
			bHiddenOutputs[j] += float32(bUpdates.Coeffs[i]) * weights[int(bUpdates.Indices[i])*NetHiddenSize+j]
		}
	}
}
