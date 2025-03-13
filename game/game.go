package game

import (
	"errors"
	"fmt"
	"judge/player"
	"time"
)

type game interface {
	validate(output string) error
	playRound(c int, player1, player2 *player.Player, verbose bool) (int, error)
	Name() string
}

func newGame(name string, count int) (game, error) {
	switch name {
	case "prisoners_dilemma":
		return NewPrisonersDilemma(), nil
	case "good_deal":
		return NewGoodDeal(100, 50), nil
	case "tug_of_war":
		return NewTugOfWar(count), nil
	case "balance_of_universe":
		return NewBalanceOfUniverse(100), nil
	default:
		return nil, errors.New("unsupported game")
	}
}

func Play(name string, count int, player1 *player.Player, player2 *player.Player, verbose bool) error {
	g, err := newGame(name, count)
	if err != nil {
		return err
	}

	player1.StartGame()
	player2.StartGame()

	for c := range count {
		if k, err := g.playRound(c, player1, player2, verbose); err != nil {
			player1.StopGame()
			player2.StopGame()
			return fmt.Errorf("error with game %s: player %d: %v", g.Name(), k, err)
		}
	}

	player1.StopGame()
	player2.StopGame()

	return nil
}

func getPlayerChoice(p *player.Player, verbose bool) (string, error) {
	choice, err := p.Receive(500 * time.Millisecond) // Таймаут 500 мс, можно поменять
	if err != nil {
		return "", fmt.Errorf("failed to get choice from player: %v", err)
	}
	if verbose {
		fmt.Printf("Player choice: %s\n", choice)
	}
	return choice, nil
}
