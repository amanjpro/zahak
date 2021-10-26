package engine

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
)

const NetInputSize = 769
const NetOutputSize = 1
const NetLayers = 1
const MaximumDepth = 128

var NetHiddenSize = 128
var CurrentHiddenWeights []float32
var CurrentHiddenBiases []float32
var CurrentOutputWeights []float32
var CurrentOutputBias float32
var CurrentNetworkId uint32

type Updates struct {
	Indices []int16
	Coeffs  []int8
	Size    int
}

func (u *Updates) Add(index int16, coeff int8) {
	u.Indices[u.Size] = index
	u.Coeffs[u.Size] = coeff
	u.Size += 1
}

type NetworkState struct {
	HiddenOutputs     [][]float32
	EmptyHiddenOutput []float32
	CurrentHidden     int
	HiddenWeights     []float32
	HiddenBiases      []float32
	OutputWeights     []float32
	OutputBias        float32
}

func NewNetworkState() *NetworkState {
	net := NetworkState{
		HiddenWeights: CurrentHiddenWeights,
		HiddenBiases:  CurrentHiddenBiases,
		OutputWeights: CurrentOutputWeights,
		OutputBias:    CurrentOutputBias,
	}

	net.EmptyHiddenOutput = make([]float32, NetHiddenSize)
	net.HiddenOutputs = make([][]float32, MaximumDepth)
	for i := 0; i < MaximumDepth; i++ {
		net.HiddenOutputs[i] = make([]float32, NetHiddenSize)
	}
	return &net
}

const Remove int8 = -1
const Add int8 = 1

func calculateNetInputIndex(sq Square, piece Piece) int16 {
	return int16(piece-1)*64 + int16(sq)
}

func (n *NetworkState) RevertHidden() {
	n.CurrentHidden -= 1
}

func (n *NetworkState) Recalculate(input []int16) {
	n.CurrentHidden = 0
	// apply hidden layer
	hiddenOutputs := n.HiddenOutputs[n.CurrentHidden]
	copy(hiddenOutputs, n.EmptyHiddenOutput)

	for index := 0; index < len(input); index++ {
		i := int(input[index])
		weights := n.HiddenWeights
		for j := 0; j < len(hiddenOutputs); j++ {
			hiddenOutputs[j] += weights[i*NetHiddenSize+j]
		}
	}
	for i := 0; i < len(hiddenOutputs); i++ {
		hiddenOutputs[i] += n.HiddenBiases[i]
	}
}

// load a neural network from file
func LoadNetwork(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Read headers
	buf := make([]byte, 4)
	_, err = io.ReadFull(f, buf)
	if err != nil {
		return err
	}
	if buf[0] != 66 || buf[1] != 90 || buf[2] != 2 || buf[3] != 0 {
		return fmt.Errorf("Magic word does not match expected, exiting")
	}

	_, err = io.ReadFull(f, buf)
	if err != nil {
		return err
	}
	id := binary.LittleEndian.Uint32(buf)

	// Read Topology Header
	buf = make([]byte, 12)
	_, err = io.ReadFull(f, buf)
	if err != nil {
		return err
	}
	inputs := binary.LittleEndian.Uint32(buf[:4])
	outputs := binary.LittleEndian.Uint32(buf[4:8])
	layers := binary.LittleEndian.Uint32(buf[8:])

	if inputs != NetInputSize || outputs != NetOutputSize || layers != NetLayers {
		return fmt.Errorf("Topology is not supported, exiting")
	}

	buf = make([]byte, 4)
	_, err = io.ReadFull(f, buf)
	if err != nil {
		return err
	}
	neurons := binary.LittleEndian.Uint32(buf)
	NetHiddenSize = int(neurons)

	CurrentNetworkId = id

	CurrentHiddenWeights = make([]float32, NetHiddenSize*NetInputSize)
	for i := 0; i < NetHiddenSize*NetInputSize; i++ {
		_, err := io.ReadFull(f, buf)
		if err != nil {
			panic(err)
		}
		CurrentHiddenWeights[i] = math.Float32frombits(binary.LittleEndian.Uint32(buf))
	}

	CurrentHiddenBiases = make([]float32, NetHiddenSize)
	for i := 0; i < NetHiddenSize; i++ {
		_, err := io.ReadFull(f, buf)
		if err != nil {
			panic(err)
		}
		CurrentHiddenBiases[i] = math.Float32frombits(binary.LittleEndian.Uint32(buf))
	}

	CurrentOutputWeights = make([]float32, NetHiddenSize)
	for i := 0; i < NetOutputSize*NetHiddenSize; i++ {
		_, err := io.ReadFull(f, buf)
		if err != nil {
			panic(err)
		}
		CurrentOutputWeights[i] = math.Float32frombits(binary.LittleEndian.Uint32(buf))
	}

	_, err = io.ReadFull(f, buf)
	if err != nil {
		panic(err)
	}
	CurrentOutputBias = math.Float32frombits(binary.LittleEndian.Uint32(buf))
	return nil
}

func ReLu(x float32) float32 {
	if x < 0 {
		return 0
	}
	return x
}
