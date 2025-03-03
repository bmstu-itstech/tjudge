package main

import (
	"flag"
	"fmt"
	"judge/game"
	"judge/player"
)

const (
	Dilemma = "prisoners_dilemma"
)

func main() {
	count := flag.Uint("c", 1, "Count of round (short)")
	flag.UintVar(count, "count", 1, "Count of round")
	verbose := flag.Bool("v", false, "Verbose mode (short)")
	flag.BoolVar(verbose, "verbose", false, "Verbose mode")
	flag.Parse()

	if flag.NArg() != 3 {
		fmt.Println("judge <args> <game> <program1> <program2>")
		fmt.Println("games:\n", Dilemma)
		flag.Usage()
		return
	}

	player1, err := player.NewPlayer(flag.Arg(1))
	if err != nil {
		fmt.Println(fmt.Errorf("error creating player 1: %w", err))
		return
	}

	player2, err := player.NewPlayer(flag.Arg(1))
	if err != nil {
		fmt.Println(fmt.Errorf("error creating player 2: %w", err))
		return
	}

	err = game.Play(flag.Arg(0), int(*count), player1, player2, *verbose)
	if err != nil {
		fmt.Println("error playing round:", err)
	}

	if *verbose {
		fmt.Println("Final scores")
	}

	fmt.Println(player1.GetScore())
	fmt.Println(player2.GetScore())
}
