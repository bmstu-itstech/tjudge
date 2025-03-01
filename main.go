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

	if flag.NArg() < 3 {
		fmt.Println("judge <args> <game> <program1> <program2>")
		fmt.Println("games:\n", Dilemma)
		flag.Usage()
		return
	}

	var players []*player.Player

	for i := 1; i < flag.NArg(); i++ {
		playeri, err := player.NewPlayer(flag.Arg(i))
		if err != nil {
			fmt.Println(fmt.Errorf("error creating player %d: %w", i, err))
			return
		}
		players = append(players, playeri)
	}

	// g, err := game.NewGame()
	// if err != nil {
	// 	fmt.Println("error creating game:", err)
	// 	return
	// }

	err := game.Play(flag.Arg(0), int(*count), players, *verbose)
	if err != nil {
		fmt.Println("error playing round:", err)
	}

	if *verbose {
		fmt.Println("Final scores")
	}
	for i := range players {
		fmt.Println(players[i].GetScore())
	}
}
