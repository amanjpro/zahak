//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"

	. "github.com/amanjpro/zahak/engine"
	// . "github.com/amanjpro/zahak/zahak"
)

var netPath = "default.nn"
var Version = "dev"

func main() {
	LoadNetwork(netPath)

	v, err := os.Create(fmt.Sprintf("zahak%cversion.go", os.PathSeparator))
	if err != nil {
		panic(err)
	}
	defer v.Close()

	v.WriteString("package main\n\n")
	v.WriteString("// Code generated by go generate; DO NOT EDIT.\n\n")
	v.WriteString("func init() {\n")
	v.WriteString(fmt.Sprintf("version = \"%s\"\n", Version))
	v.WriteString("}\n")

	f, err := os.Create(fmt.Sprintf("engine%cnn.go", os.PathSeparator))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("package engine\n\n")
	f.WriteString("// Code generated by go generate; DO NOT EDIT.\n\n")
	f.WriteString("func init() {\n")

	f.WriteString("CurrentHiddenWeights = []float32 {\n")
	for i := 0; i < len(CurrentHiddenWeights); i++ {
		if i%10 == 0 {
			f.WriteString("\n")
		}
		f.WriteString(fmt.Sprintf("%g,", CurrentHiddenWeights[i]))
	}
	f.WriteString("\n}\n")

	f.WriteString("CurrentHiddenBiases = []float32 {\n")
	for i := 0; i < len(CurrentHiddenBiases); i++ {
		if i%10 == 0 {
			f.WriteString("\n")
		}
		f.WriteString(fmt.Sprintf("%g,", CurrentHiddenBiases[i]))
	}
	f.WriteString("\n}\n")

	f.WriteString("CurrentOutputWeights = []float32 {\n")
	for i := 0; i < len(CurrentOutputWeights); i++ {
		if i%10 == 0 {
			f.WriteString("\n")
		}
		f.WriteString(fmt.Sprintf("%g,", CurrentOutputWeights[i]))
	}
	f.WriteString("\n}\n")

	f.WriteString(fmt.Sprintf("CurrentOutputBias = %g\n", CurrentOutputBias))
	f.WriteString(fmt.Sprintf("CurrentNetworkId = %d\n", CurrentNetworkId))
	f.WriteString(fmt.Sprintf("NetHiddenSize = %d\n", NetHiddenSize))

	f.WriteString("}\n")
}
