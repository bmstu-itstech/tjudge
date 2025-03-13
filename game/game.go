package game

import (
	"errors"
	"fmt"
	"judge/player"
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

	err = player1.StartGame()
	if err != nil {
		return fmt.Errorf("error while game starting player 1: %v", err)
	}
	err = player2.StartGame()
	if err != nil {
		return fmt.Errorf("error while game starting player 2: %v", err)
	}

	for c := range count {
		if k, err := g.playRound(c, player1, player2, verbose); err != nil {
			err2 := player1.StopGame()
			if err2 != nil {
				err2 = fmt.Errorf("error while game stoping player 1: %v", err2)
			}
			err3 := player2.StopGame()
			if err3 != nil {
				err3 = fmt.Errorf("error while game stoping player 2: %v", err3)
			}
			return errors.Join(err2, err3, fmt.Errorf("error with game %s: player %d: %v", g.Name(), k, err))
		}
	}

	err = player1.StopGame()
	if err != nil {
		err = fmt.Errorf("error while game stoping player 1: %v", err)
	}
	err2 := player2.StopGame()
	if err2 != nil {
		err2 = fmt.Errorf("error while game stoping player 2: %v", err)
	}

	return errors.Join(err, err2)
}
