package tuning

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

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
var RANGE_END = int16(250)
var NUM_PROCESSORS = 8
var initialK = 1.0
var answers = make(chan float64)
var ml = NewMoveList(500)

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
	guesses = append(guesses, EarlyPawnPst[:]...)                 // 0-63
	guesses = append(guesses, LatePawnPst[:]...)                  // 64-127
	guesses = append(guesses, EarlyKnightPst[:]...)               // 128-191
	guesses = append(guesses, LateKnightPst[:]...)                // 192-255
	guesses = append(guesses, EarlyBishopPst[:]...)               // 256-319
	guesses = append(guesses, LateBishopPst[:]...)                // 320-383
	guesses = append(guesses, EarlyRookPst[:]...)                 // 384-447
	guesses = append(guesses, LateRookPst[:]...)                  // 448-511
	guesses = append(guesses, EarlyQueenPst[:]...)                // 512-575
	guesses = append(guesses, LateQueenPst[:]...)                 // 576-639
	guesses = append(guesses, EarlyKingPst[:]...)                 // 640-703
	guesses = append(guesses, LateKingPst[:]...)                  // 704-767
	guesses = append(guesses, BackwardPawnAward)                  // 768
	guesses = append(guesses, IsolatedPawnAward)                  // 769
	guesses = append(guesses, DoublePawnAward)                    // 770
	guesses = append(guesses, EndgamePassedPawnAward)             // 771
	guesses = append(guesses, MiddlegamePassedPawnAward)          // 772
	guesses = append(guesses, EndgameCandidatePassedPawnAward)    // 773
	guesses = append(guesses, MiddlegameCandidatePassedPawnAward) // 774
	guesses = append(guesses, RookOpenFileAward)                  // 775
	guesses = append(guesses, RookSemiOpenFileAward)              // 776
	guesses = append(guesses, VeritcalDoubleRookAward)            // 777
	guesses = append(guesses, HorizontalDoubleRookAward)          // 778
	guesses = append(guesses, PawnFactorCoeff)                    // 779
	guesses = append(guesses, AggressivityFactorCoeff)            // 780
	guesses = append(guesses, MiddlegameAggressivityFactorCoeff)  // 781
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
	BackwardPawnAward = guesses[768]
	IsolatedPawnAward = guesses[769]
	DoublePawnAward = guesses[770]
	EndgamePassedPawnAward = guesses[771]
	MiddlegamePassedPawnAward = guesses[772]
	EndgameCandidatePassedPawnAward = guesses[773]
	MiddlegameCandidatePassedPawnAward = guesses[774]
	RookOpenFileAward = guesses[775]
	RookSemiOpenFileAward = guesses[776]
	VeritcalDoubleRookAward = guesses[777]
	HorizontalDoubleRookAward = guesses[778]
	PawnFactorCoeff = guesses[779]
	AggressivityFactorCoeff = guesses[780]
	MiddlegameAggressivityFactorCoeff = guesses[781]
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

	fmt.Printf("var BackwardPawnAward int16 = %d\n", guesses[768])
	fmt.Printf("var IsolatedPawnAward int16 = %d\n", guesses[769])
	fmt.Printf("var DoublePawnAward int16 = %d\n", guesses[770])
	fmt.Printf("var EndgamePassedPawnAward int16 = %d\n", guesses[771])
	fmt.Printf("var MiddlegamePassedPawnAward int16 = %d\n", guesses[772])
	fmt.Printf("var EndgameCandidatePassedPawnAward int16 = %d\n", guesses[773])
	fmt.Printf("var MiddlegameCandidatePassedPawnAward int16 = %d\n", guesses[774])
	fmt.Printf("var RookOpenFileAward int16 = %d\n", guesses[775])
	fmt.Printf("var RookSemiOpenFileAward int16 = %d\n", guesses[776])
	fmt.Printf("var VeritcalDoubleRookAward int16 = %d\n", guesses[777])
	fmt.Printf("var HorizontalDoubleRookAward int16 = %d\n", guesses[778])
	fmt.Printf("var PawnFactorCoeff int16 = %d\n", guesses[779])
	fmt.Printf("var AggressivityFactorCoeff int16 = %d\n", guesses[780])
	fmt.Printf("var MiddlegameAggressivityFactorCoeff int16 = %d\n", guesses[781])
	fmt.Println("===================================================")
}

func updateSingleGradient(testPosition TestPosition, params []int16, K float64, gradient *[]float64) {

	updateEvalParams(params)
	err := linearEvaluation(testPosition)
	sig := sigmoid(err, K)
	x := (testPosition.outcome - sig) * sig * (1 - sig)

	for i := 0; i < len(params); i++ {
		(*gradient)[i] += x * float64(params[i])
	}
}

func computeGradient(testPositions []TestPosition, gradient []float64, guesses []int16, K float64) []int16 {

	var local = make([]float64, len(guesses))

	for i := 0; i < len(guesses); i++ {
		updateSingleGradient(testPositions[i], guesses, K, &local)
	}

	for i := 0; i < len(gradient); i++ {
		gradient[i] += local[i]
	}

	return guesses
}

//
// var MAXEPOCHS = 10000
// func findOptimalParams(guesses []int16) {
// 	for epoch := 1; epoch <= MAXEPOCHS; epoch++ {
//         for int batch = 0; batch < NPOSITIONS / BATCHSIZE; batch++) {
// 	var gradient = make([]float64, len(guesses))
//
//             computeGradient(entries, gradient, guesses, K);
//
//             for (int i = 0; i < NTERMS; i++) {
//                 adagrad[i][MG] += pow((K / 200.0) * gradient[i][MG] / BATCHSIZE, 2.0);
//                 adagrad[i][EG] += pow((K / 200.0) * gradient[i][EG] / BATCHSIZE, 2.0);
//                 params[i][MG] += (K / 200.0) * (gradient[i][MG] / BATCHSIZE) * (rate / sqrt(1e-8 + adagrad[i][MG]));
//                 params[i][EG] += (K / 200.0) * (gradient[i][EG] / BATCHSIZE) * (rate / sqrt(1e-8 + adagrad[i][EG]));
//             }
//         }
//
//         error = TunedEvaluationErrors(entries, params, K);
//         printf("\rEpoch [%d] Error = [%.8f], Rate = [%g]", epoch, error, rate);
//
//         // Pre-scheduled Learning Rate drops
//         if (epoch % LRSTEPRATE == 0) rate = rate / LRDROPRATE;
//         if (epoch % REPORTING == 0) PrintParameters(params, currentParams);
//     }
// }

func localOptimize(initialGuess []int16, K float64) []int16 {
	nParams := len(initialGuess)
	bestE := E(initialGuess, K)
	bestParValues := initialGuess
	improved := true
	for improved {
		improved = false
		for pi := 0; pi < nParams; pi++ {
			if _, ok := futileIndices[pi]; ok {
				continue
			}
			newParValues := bestParValues
			newParValues[pi] += 1
			newE := E(newParValues, K)
			if newE < bestE {
				bestE = newE
				bestParValues = newParValues
				fmt.Println("Best parameters so far")
				printOptimalGuesses(bestParValues)
				improved = true
			} else {
				newParValues[pi] -= 2
				newE = E(newParValues, K)
				if newE < bestE {
					bestE = newE
					bestParValues = newParValues
					fmt.Println("Best parameters so far")
					printOptimalGuesses(bestParValues)
					improved = true
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

func linearEvaluation(testPosition TestPosition) int16 {
	pos := testPosition.pos
	eval := Evaluate(pos)
	if pos.Turn() == Black {
		eval = -eval
	}
	return eval
}

func processor(start int, end int, K float64) {
	var acc float64 = 0
	for i := start; i < end; i++ {
		eval := linearEvaluation(testPositions[i])
		e := testPositions[i].outcome - sigmoid(eval, K)
		acc += e * e
	}

	answers <- acc
}

func E(guess []int16, K float64) float64 {
	acc := float64(0)
	updateEvalParams(guess)

	chunk := len(testPositions) / NUM_PROCESSORS

	for i := 0; i < NUM_PROCESSORS; i++ {
		start := i * chunk
		end := (i+1)*chunk - 1
		if i == NUM_PROCESSORS-1 {
			end = len(testPositions)
		}
		go processor(start, end, K)
	}

	for i := 0; i < NUM_PROCESSORS; i++ {
		ans := <-answers
		acc += ans
	}
	return 1 * acc / float64(len(testPositions))
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
	fen := strings.Trim(strings.Join(fields[:6], " "), ";")

	outcomeStr := strings.Replace(fields[8], "pgn=", "", -1)
	outcome, e := strconv.ParseFloat(outcomeStr, 64)
	if e != nil {
		panic(e)
	}
	if fields[1] == "b" && outcomeStr != "0.5" {
		if outcomeStr == "0.0" {
			outcome = 1
		} else {
			outcome = 0
		}
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
		fmt.Println(line)
	})
}

func Tune(path string) {
	loadPositions(path, func(line string) {
		fen, outcome := parseLine(line)
		game := FromFen(fen, true)
		pos := game.Position()
		tp := TestPosition{pos, outcome}
		testPositions = append(testPositions, tp)
	})
	fmt.Printf("%d positions loaded\n", len(testPositions))
	K := findK()
	fmt.Printf("Optimal K is %f\n", K)
	optimalGuesses := localOptimize(initialGuesses, K)
	fmt.Println("Optimal Parameters have been found!!")
	fmt.Println("===================================================")
	printOptimalGuesses(optimalGuesses)
	close(answers)
}
