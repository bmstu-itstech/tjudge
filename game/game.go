package game

import (
	"errors"
	"fmt"
	"judge/player"
	"time"
)

type Game struct {
	name    string
	players *[]*player.Player
}

func NewGame(name string, players *[]*player.Player) (*Game, error) {
	switch name {
	case "prisoners_dilemma":
		break
	default:
		return nil, fmt.Errorf("unsupported game")
	}
	return &Game{name: name, players: players}, nil
}

func (g *Game) Play(count int, verbose bool) error {
	var ignore map[int]error = make(map[int]error)

	for i, player1 := range *g.players {
		for j, player2 := range *g.players {
			if i == j {
				continue
			}

			// Игроки не нарушали правил
			if _, ok := ignore[i]; ok {
				continue
			}
			if _, ok := ignore[j]; ok {
				continue
			}

			player1.StartGame()
			player2.StartGame()

			flag := false
			for c := range count {
				if flag {
					break
				}
				switch g.name {
				case "prisoners_dilemma":
					if k, err := playRoundPrisonersDilemma(c, i, j, player1, player2, verbose); err != nil {
						ignore[k] = err
						flag = true
					}
				}
			}

			player1.StopGame()
			player2.StopGame()
		}
	}

	var err error

	for i := range ignore {
		err = errors.Join(ignore[i])
	}

	return err
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
