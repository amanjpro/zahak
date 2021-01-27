package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/notnil/chess"
)

func main() {
	game := chess.NewGame()
	// generate moves until game is over
	for game.Outcome() == chess.NoOutcome {
		// select a random move
		if game.Position().Turn() == chess.White {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter your next move: ")
			move, _ := reader.ReadString('\n')
			err := game.MoveStr(strings.TrimSpace(move))
			if err != nil {
				fmt.Println(err)
			}
		} else {
			move := search(game.Position(), 5)
			game.Move(move)
		}
		// Position currently
		fmt.Println(game.Position().Board().Draw())
	}
	fmt.Printf("Game completed. %s by %s.\n", game.Outcome(), game.Method())
	fmt.Println(game.String())
	/*
		Output:

		 A B C D E F G H
		8- - - - - - - -
		7- - - - - - ♚ -
		6- - - - ♗ - - -
		5- - - - - - - -
		4- - - - - - - -
		3♔ - - - - - - -
		2- - - - - - - -
		1- - - - - - - -

		Game completed. 1/2-1/2 by InsufficientMaterial.

		1.Nc3 b6 2.a4 e6 3.d4 Bb7 ...
	*/
}
