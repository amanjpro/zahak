package tuning

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	// "sync"

	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
	. "github.com/amanjpro/zahak/search"
)

type TestPosition struct {
	pos     *Position
	outcome float64
}

var testPositions []TestPosition
var initialGuesses = computeInitialGuesses()
var K_PRECISION = 10
var NUM_PROCESSORS = 8
var initialK = 1.0
var answers = make(chan float64)
var ml = NewMoveList(500)
var MAXEPOCHS = 10000
var BATCHSIZE int
var REPORTING = 50
var LRSTEPRATE = 250
var LRATE = 1.00
var LRDROPRATE = 1.00

func initEngines() []*Engine {
	res := make([]*Engine, NUM_PROCESSORS)
	for i := 0; i < NUM_PROCESSORS; i++ {
		res[i] = NewEngine(nil)
	}
	return res
}

var futileIndices map[int]bool = initFutileIndices()

func initFutileIndices() map[int]bool {
	futileIndices := []int{
		0, 1, 2, 3, 4, 5, 6, 7,
		56, 57, 58, 59, 60, 61, 62, 63,
		64, 65, 66, 67, 68, 69, 70, 71,
		120, 121, 122, 123, 124, 125, 126, 127,
	}

	res := make(map[int]bool, len(futileIndices))

	for _, v := range futileIndices {
		res[v] = true
	}

	return res
}

func computeInitialGuesses() []int16 {
	var guesses = make([]int16, 0, 800)
	guesses = append(guesses, EarlyPawnPst[:]...)                  // 0-63
	guesses = append(guesses, LatePawnPst[:]...)                   // 64-127
	guesses = append(guesses, EarlyKnightPst[:]...)                // 128-191
	guesses = append(guesses, LateKnightPst[:]...)                 // 192-255
	guesses = append(guesses, EarlyBishopPst[:]...)                // 256-319
	guesses = append(guesses, LateBishopPst[:]...)                 // 320-383
	guesses = append(guesses, EarlyRookPst[:]...)                  // 384-447
	guesses = append(guesses, LateRookPst[:]...)                   // 448-511
	guesses = append(guesses, EarlyQueenPst[:]...)                 // 512-575
	guesses = append(guesses, LateQueenPst[:]...)                  // 576-639
	guesses = append(guesses, EarlyKingPst[:]...)                  // 640-703
	guesses = append(guesses, LateKingPst[:]...)                   // 704-767
	guesses = append(guesses, MiddlegameBackwardPawnAward)         // 768
	guesses = append(guesses, EndgameBackwardPawnAward)            // 769
	guesses = append(guesses, MiddlegameIsolatedPawnAward)         // 770
	guesses = append(guesses, EndgameIsolatedPawnAward)            // 771
	guesses = append(guesses, MiddlegameDoublePawnAward)           // 772
	guesses = append(guesses, EndgameDoublePawnAward)              // 773
	guesses = append(guesses, MiddlegamePassedPawnAward)           // 774
	guesses = append(guesses, EndgamePassedPawnAward)              // 775
	guesses = append(guesses, MiddlegameCandidatePassedPawnAward)  // 776
	guesses = append(guesses, EndgameCandidatePassedPawnAward)     // 777
	guesses = append(guesses, MiddlegameRookOpenFileAward)         // 778
	guesses = append(guesses, EndgameRookOpenFileAward)            // 779
	guesses = append(guesses, MiddlegameRookSemiOpenFileAward)     // 780
	guesses = append(guesses, EndgameRookSemiOpenFileAward)        // 781
	guesses = append(guesses, MiddlegameVeritcalDoubleRookAward)   // 782
	guesses = append(guesses, EndgameVeritcalDoubleRookAward)      // 783
	guesses = append(guesses, MiddlegameHorizontalDoubleRookAward) // 784
	guesses = append(guesses, EndgameHorizontalDoubleRookAward)    // 785
	guesses = append(guesses, MiddlegamePawnFactorCoeff)           // 786
	guesses = append(guesses, EndgamePawnFactorCoeff)              // 787
	guesses = append(guesses, MiddlegameAggressivityFactorCoeff)   // 788
	guesses = append(guesses, EndgameAggressivityFactorCoeff)      // 789
	guesses = append(guesses, MiddlegameCastlingAward)             // 790
	return guesses
}

func updateEvalParams(guesses []int16) {
	for i := 0; i < 64; i++ {
		EarlyPawnPst[i] = guesses[i+0*64]
		LatePawnPst[i] = guesses[i+1*64]
		EarlyKnightPst[i] = guesses[i+2*64]
		LateKnightPst[i] = guesses[i+3*64]
		EarlyBishopPst[i] = guesses[i+4*64]
		LateBishopPst[i] = guesses[i+5*64]
		EarlyRookPst[i] = guesses[i+6*64]
		LateRookPst[i] = guesses[i+7*64]
		EarlyQueenPst[i] = guesses[i+8*64]
		LateQueenPst[i] = guesses[i+9*64]
		EarlyKingPst[i] = guesses[i+10*64]
		LateKingPst[i] = guesses[i+11*64]
	}
	MiddlegameBackwardPawnAward = guesses[768]
	EndgameBackwardPawnAward = guesses[769]
	MiddlegameIsolatedPawnAward = guesses[770]
	EndgameIsolatedPawnAward = guesses[771]
	MiddlegameDoublePawnAward = guesses[772]
	EndgameDoublePawnAward = guesses[773]
	MiddlegamePassedPawnAward = guesses[774]
	EndgamePassedPawnAward = guesses[775]
	MiddlegameCandidatePassedPawnAward = guesses[776]
	EndgameCandidatePassedPawnAward = guesses[777]
	MiddlegameRookOpenFileAward = guesses[778]
	EndgameRookOpenFileAward = guesses[779]
	MiddlegameRookSemiOpenFileAward = guesses[780]
	EndgameRookSemiOpenFileAward = guesses[781]
	MiddlegameVeritcalDoubleRookAward = guesses[782]
	EndgameVeritcalDoubleRookAward = guesses[783]
	MiddlegameHorizontalDoubleRookAward = guesses[784]
	EndgameHorizontalDoubleRookAward = guesses[785]
	MiddlegamePawnFactorCoeff = guesses[786]
	EndgamePawnFactorCoeff = guesses[787]
	MiddlegameAggressivityFactorCoeff = guesses[788]
	EndgameAggressivityFactorCoeff = guesses[789]
	MiddlegameCastlingAward = guesses[790]
}

func toEvalParams(guesses []float64) []int16 {
	params := make([]int16, len(guesses))
	for i, v := range guesses {
		params[i] = int16(v)
	}
	return params
}

func printPst(pst []int16, varname string) {
	acc := ""
	for i := 0; i < 64; i++ {
		if i%8 == 0 {
			acc = fmt.Sprintf("%s\n", acc)
		}
		acc = fmt.Sprintf("%s %d,", acc, pst[i])
	}
	fmt.Printf("var %s = [64]int16 { %s \n}\n\n", varname, acc)
}

func printOptimalGuesses(guesses []int16) {
	fmt.Println("// Middle-game")
	printPst(guesses[0:64], "EarlyPawnPst")
	printPst(guesses[128:192], "EarlyKnightPst")
	printPst(guesses[256:320], "EarlyBishopPst")
	printPst(guesses[384:448], "EarlyRookPst")
	printPst(guesses[512:576], "EarlyQueenPst")
	printPst(guesses[640:704], "EarlyKingPst")
	fmt.Println("// Endgame")
	printPst(guesses[64:128], "LatePawnPst")
	printPst(guesses[192:256], "LateKnightPst")
	printPst(guesses[320:384], "LateBishopPst")
	printPst(guesses[448:512], "LateRookPst")
	printPst(guesses[576:640], "LateQueenPst")
	printPst(guesses[704:768], "LateKingPst")

	fmt.Printf("var MiddlegameBackwardPawnAward int16 = %d\n", guesses[768])
	fmt.Printf("var EndgameBackwardPawnAward int16 = %d\n", guesses[769])
	fmt.Printf("var MiddlegameIsolatedPawnAward int16 = %d\n", guesses[770])
	fmt.Printf("var EndgameIsolatedPawnAward int16 = %d\n", guesses[771])
	fmt.Printf("var MiddlegameDoublePawnAward int16 = %d\n", guesses[772])
	fmt.Printf("var EndgameDoublePawnAward int16 = %d\n", guesses[773])
	fmt.Printf("var MiddlegamePassedPawnAward int16 = %d\n", guesses[774])
	fmt.Printf("var EndgamePassedPawnAward int16 = %d\n", guesses[775])
	fmt.Printf("var MiddlegameCandidatePassedPawnAward int16 = %d\n", guesses[776])
	fmt.Printf("var EndgameCandidatePassedPawnAward int16 = %d\n", guesses[777])
	fmt.Printf("var MiddlegameRookOpenFileAward int16 = %d\n", guesses[778])
	fmt.Printf("var EndgameRookOpenFileAward int16 = %d\n", guesses[779])
	fmt.Printf("var MiddlegameRookSemiOpenFileAward int16 = %d\n", guesses[780])
	fmt.Printf("var EndgameRookSemiOpenFileAward int16 = %d\n", guesses[781])
	fmt.Printf("var MiddlegameVeritcalDoubleRookAward int16 = %d\n", guesses[782])
	fmt.Printf("var EndgameVeritcalDoubleRookAward int16 = %d\n", guesses[783])
	fmt.Printf("var MiddlegameHorizontalDoubleRookAward int16 = %d\n", guesses[784])
	fmt.Printf("var EndgameHorizontalDoubleRookAward int16 = %d\n", guesses[785])
	fmt.Printf("var MiddlegamePawnFactorCoeff int16 = %d\n", guesses[786])
	fmt.Printf("var EndgamePawnFactorCoeff int16 = %d\n", guesses[787])
	fmt.Printf("var MiddlegameAggressivityFactorCoeff int16 = %d\n", guesses[788])
	fmt.Printf("var EndgameAggressivityFactorCoeff int16 = %d\n", guesses[789])
	fmt.Printf("var MiddlegameCastlingAward int16 = %d\n", guesses[790])
	fmt.Println("===================================================")
}

func localOptimize(initialGuess []int16, K float64) []int16 {
	nParams := len(initialGuess)
	bestE := E(initialGuess, K)
	bestParValues := append([]int16{}, initialGuess...)
	improved := true
	for improved {
		improved = false
		for pi := 0; pi < nParams; pi++ {
			// if pi < 768 {
			// 	continue
			// }
			if _, ok := futileIndices[pi]; ok {
				continue
			}
			bestParValues[pi] += 1
			newE := E(bestParValues, K)
			if newE < bestE {
				bestE = newE
				fmt.Println("Best parameters so far")
				printOptimalGuesses(bestParValues)
				improved = true
			} else {
				bestParValues[pi] -= 2
				newE = E(bestParValues, K)
				if newE < bestE {
					bestE = newE
					fmt.Println("Best parameters so far")
					printOptimalGuesses(bestParValues)
					improved = true
				} else {
					bestParValues[pi] += 1 // reset the guess
				}
			}
		}
	}
	return bestParValues
}

func findK() float64 {
	start := 0.0
	var end float64 = 10
	step := 1.0
	curr := start
	var err float64
	best := E(initialGuesses, start)

	for i := 0; i < K_PRECISION; i++ {

		// Find the minimum within [start, end] using the current step
		curr = start - step
		for curr < end {
			curr = curr + step
			err = E(initialGuesses, curr)
			if err <= best {
				best = err
				start = curr
			}
		}

		fmt.Printf("K so far is %f, iteration %d\n", start, i)

		// Adjust the search space
		end = start + step
		start = start - step
		step = step / 10.0
	}

	return start
}

func linearEvaluation(pos *Position) int16 {
	eval := Evaluate(pos)
	if pos.Turn() == Black {
		eval = -eval
	}
	return eval
}

func processor(start int, end int, K float64) {
	var acc float64 = 0
	for i := start; i < end; i++ {
		eval := linearEvaluation(testPositions[i].pos)
		acc = math.Pow(testPositions[i].outcome-sigmoid(eval, K), 2)
	}

	answers <- acc
}

func E(guess []int16, K float64) float64 {
	acc := float64(0)
	updateEvalParams(guess)

	for i := 0; i < NUM_PROCESSORS; i++ {
		start := i * BATCHSIZE
		end := (i + 1) * BATCHSIZE
		if i == NUM_PROCESSORS-1 {
			end = len(testPositions)
		}
		go processor(start, end, K)
	}

	for i := 0; i < NUM_PROCESSORS; i++ {
		ans := <-answers
		acc += ans
	}
	return (1 / float64(len(testPositions))) * acc
}

func sigmoid(E int16, K float64) float64 {
	exp := K * float64(E) / 400
	return 1 / (1 + math.Pow(10, -exp))
}

func loadPositions(path string, actionFn func(string)) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	testPositions = make([]TestPosition, 0, 14_000_000)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		actionFn(line)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func parseLine(line string) (string, float64) {
	fields := strings.Fields(line)
	// fen := strings.Join(fields[:4], " ")
	// fen = fmt.Sprintf("%s 0 1", fen)
	//
	// outcomeStr := strings.Trim(fields[5], "\";")
	// var outcome float64
	// if outcomeStr == "1/2-1/2" {
	// 	outcome = 0.5
	// } else if outcomeStr == "1-0" {
	// 	outcome = 1.0
	// } else if outcomeStr == "0-1" {
	// 	outcome = 0.0
	// } else {
	// 	panic(fmt.Sprintf("Unexpected output %s", outcomeStr))
	// }

	fen := strings.Trim(strings.Join(fields[:6], " "), ";")

	outcomeStr := strings.Replace(fields[8], "pgn=", "", -1)
	outcome, e := strconv.ParseFloat(outcomeStr, 64)
	if e != nil {
		panic(e)
	}
	if fields[1] == "b" && outcomeStr == "1.0" {
		outcome = 0
	} else if fields[1] == "b" && outcomeStr == "0.0" {
		outcome = 1
	}

	return fen, outcome

}

func PrepareTuningData(path string) {
	loadPositions(path, func(line string) {
		fen, _ := parseLine(line)
		game := FromFen(fen, true)
		pos := game.Position()
		ml.Size = 0
		ml.Next = 0
		if pos.IsInCheck() {
			return
		}

		pos.GetCaptureMoves(ml)
		if ml.Size > 0 {
			return
		}
	})
}

func updateSingleGradient(testPosition TestPosition, gradient *[]float64, params []float64, K float64) {

	updateEvalParams(toEvalParams(params))
	E := linearEvaluation(testPosition.pos)
	S := sigmoid(E, K)
	X := (testPosition.outcome - S) * S * (1 - S)

	for i := 0; i < len(params); i++ {
		(*gradient)[i] += X
	}
}

func computeGradient(gradient *[]float64, params []float64, K float64, batch int) {

	local := make([]float64, len(initialGuesses))
	start := batch * BATCHSIZE
	end := (batch + 1) * BATCHSIZE
	if batch == NUM_PROCESSORS-1 {
		end = len(testPositions)
	}
	start = 0
	for i := start; i < end; i++ {
		entry := testPositions[i]
		updateSingleGradient(entry, &local, params, K)
	}

	for i := 0; i < len(initialGuesses); i++ {
		(*gradient)[i] += local[i]
	}
}

func gradientOptimize(K float64) {
	params := make([]float64, len(initialGuesses))
	adagrad := make([]float64, len(initialGuesses))
	rate := LRATE
	for epoch := 1; epoch <= MAXEPOCHS; epoch++ {
		// var waitGroup sync.WaitGroup
		// waitGroup.Add(NUM_PROCESSORS)
		for batch := 0; batch < NUM_PROCESSORS; batch++ {
			gradient := make([]float64, len(initialGuesses))
			// go func(batch int) {
			// 	defer waitGroup.Done()
			computeGradient(&gradient, params, K, batch)
			for i := 0; i < len(initialGuesses); i++ {
				adagrad[i] += math.Pow((K/200.0)*gradient[i]/float64(BATCHSIZE), 2.0)
				adagrad[i] += math.Pow((K/200.0)*gradient[i]/float64(BATCHSIZE), 2.0)
				params[i] += (K / 200.0) * (gradient[i] / float64(BATCHSIZE)) * (rate / math.Sqrt(1e-8+adagrad[i]))
				params[i] += (K / 200.0) * (gradient[i] / float64(BATCHSIZE)) * (rate / math.Sqrt(1e-8+adagrad[i]))
			}
			//
			// }(batch)
		}

		// waitGroup.Wait()

		err := E(toEvalParams(params), K)
		fmt.Printf("\rEpoch [%d] Error = [%.8f], Rate = [%g]\n", epoch, err, rate)
		// Pre-scheduled Learning Rate drops
		if epoch%LRSTEPRATE == 0 {
			rate = rate / LRDROPRATE
		}
		if epoch%REPORTING == 0 {
			printOptimalGuesses(toEvalParams(params))
		}
	}
}

func Tune(path string) {
	loadPositions(path, func(line string) {
		fen, outcome := parseLine(line)
		game := FromFen(fen, true)
		pos := game.Position()
		tp := TestPosition{pos, outcome}
		testPositions = append(testPositions, tp)
	})

	BATCHSIZE = len(testPositions) / NUM_PROCESSORS
	fmt.Printf("%d positions loaded\n", len(testPositions))
	K := findK()
	fmt.Printf("Optimal K is %f\n", K)
	optimalGuesses := localOptimize(initialGuesses, K)
	// gradientOptimize(K)
	fmt.Println("Optimal Parameters have been found!!")
	fmt.Println("===================================================")
	printOptimalGuesses(optimalGuesses)
	close(answers)
}
