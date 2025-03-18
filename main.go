package main

import (
	"flag"
	"fmt"
	"judge/game"
	"judge/player"
	"os"
)

func Errorf(format string, a ...interface{}) {
	format = fmt.Sprintf("error: %s\n", format)
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
}

func main() {
	count := flag.Uint("c", 1, "Count of round (short)")
	flag.UintVar(count, "count", 10, "Count of round")
	verbose := flag.Bool("v", false, "Verbose mode (short)")
	flag.BoolVar(verbose, "verbose", false, "Verbose mode")
	flag.Parse()
	flag.Usage = func() {
		fmt.Printf("Usage of %s [args] game prorgam1 program2:\n", os.Args[0])
		fmt.Println("Args:")
		flag.PrintDefaults()
	}

	if flag.NArg() != 3 {
		flag.Usage()
		return
	}

	player1, err := player.NewPlayer(flag.Arg(1))
	if err != nil {
		Errorf("failed to create player 1: %s", err)
		os.Exit(3)
	}

	player2, err := player.NewPlayer(flag.Arg(2))
	if err != nil {
		Errorf("failed to create player 2: %s", err)
		os.Exit(3)
	}

	k, err := game.Play(flag.Arg(0), int(*count), player1, player2, *verbose)
	if err != nil {
		Errorf("error while playing round: %s", err)
		os.Exit(k)
	}

	if *verbose {
		fmt.Println("Final scores")
	}

	fmt.Println(player1.GetScore())
	fmt.Println(player2.GetScore())
}
