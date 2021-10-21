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
const QPrecisionIn int16 = 16
const QPrecisionOut int16 = 64

var NetHiddenSize = 128
var CurrentHiddenWeights []int16
var CurrentHiddenBiases []int16
var CurrentOutputWeights []int16
var CurrentOutputBias int32
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
	HiddenOutputs     [][]int16
	EmptyHiddenOutput []int16
	CurrentHidden     int
	HiddenWeights     []int16
	HiddenBiases      []int16
	OutputWeights     []int16
	OutputBias        int32
}

func NewNetworkState() *NetworkState {
	net := NetworkState{
		HiddenWeights: CurrentHiddenWeights,
		HiddenBiases:  CurrentHiddenBiases,
		OutputWeights: CurrentOutputWeights,
		OutputBias:    CurrentOutputBias,
	}

	net.EmptyHiddenOutput = make([]int16, NetHiddenSize)
	net.HiddenOutputs = make([][]int16, MaximumDepth)
	for i := 0; i < MaximumDepth; i++ {
		net.HiddenOutputs[i] = make([]int16, NetHiddenSize)
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
	if buf[0] != 66 || buf[1] != 90 {
		return fmt.Errorf("Magic word does not match expected, exiting")
	}

	if buf[2] != 2 || buf[3] != 0 {
		return fmt.Errorf("Network binary format version is not supported")
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

	CurrentHiddenWeights = make([]int16, NetHiddenSize*NetInputSize)
	for i := 0; i < NetHiddenSize*NetInputSize; i++ {
		_, err := io.ReadFull(f, buf)
		if err != nil {
			panic(err)
		}
		CurrentHiddenWeights[i] = quantize(math.Float32frombits(binary.LittleEndian.Uint32(buf)), false)
	}

	CurrentHiddenBiases = make([]int16, NetHiddenSize)
	for i := 0; i < NetHiddenSize; i++ {
		_, err := io.ReadFull(f, buf)
		if err != nil {
			panic(err)
		}
		CurrentHiddenBiases[i] = quantize(math.Float32frombits(binary.LittleEndian.Uint32(buf)), false)
	}

	CurrentOutputWeights = make([]int16, NetHiddenSize)
	for i := 0; i < NetOutputSize*NetHiddenSize; i++ {
		_, err := io.ReadFull(f, buf)
		if err != nil {
			panic(err)
		}
		CurrentOutputWeights[i] = quantize(math.Float32frombits(binary.LittleEndian.Uint32(buf)), true)
	}

	_, err = io.ReadFull(f, buf)
	if err != nil {
		panic(err)
	}
	ob := math.Float32frombits(binary.LittleEndian.Uint32(buf))
	CurrentOutputBias = int32(math.Round(float64(ob) * float64(QPrecisionOut)))
	return nil
}

func ReLu(x int16) int16 {
	if x < 0 {
		return 0
	}
	return x
}

func quantize(x float32, outputLayer bool) int16 {
	q := float64(QPrecisionIn)
	if outputLayer {
		q = float64(QPrecisionOut)
	}
	return int16(math.Round(float64(x) * q))
}
