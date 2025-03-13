package game

import (
	"fmt"
	"strconv"
	"time"

	"judge/player"
)

type TugOfWar struct {
	m                int
	energy1, energy2 int
}

func NewTugOfWar(m int) *TugOfWar {
	return &TugOfWar{
		m:       m,
		energy1: 100,
		energy2: 100,
	}
}

func (g *TugOfWar) Name() string {
	return "tug_of_war"
}

func (g *TugOfWar) validate(output string) error {
	return nil
}

func (g *TugOfWar) playRound(c int, player1, player2 *player.Player, verbose bool) (int, error) {
	if c == 0 {
		if err := player1.Send(strconv.Itoa(g.m)); err != nil {
			return 1, fmt.Errorf("failed to send m to player 1: %v", err)
		}
		if err := player2.Send(strconv.Itoa(g.m)); err != nil {
			return 2, fmt.Errorf("failed to send m to player 2: %v", err)
		}
	}

	if g.energy1 <= 0 && g.energy2 <= 0 {
		if verbose {
			fmt.Printf("Round %d: player 1 has %d energy and player 2 has %d energy, so round is skipped.\n",
				c+1, g.energy1, g.energy2)
		}
		return 0, nil
	}
	if g.energy1 <= 0 {
		player2.AddScore(1)
		if verbose {
			fmt.Printf("Round %d: player 1 has spent all his energy, so player 2 wins.\n",
				c+1)
		}
		return 0, nil
	}
	if g.energy2 <= 0 {
		player1.AddScore(1)
		if verbose {
			fmt.Printf("Round %d: player 2 has spent all his energy, so player 1 wins.\n",
				c+1)
		}
		return 0, nil
	}

	var total1, total2, y int

	for {
		if err := player1.Send(strconv.Itoa(y)); err != nil {
			return 1, fmt.Errorf("failed to send offer to player 1: %v", err)
		}
		xStr, err := player1.Receive(500 * time.Millisecond)
		if err != nil {
			return 1, fmt.Errorf("player 1 failed to respond: %v", err)
		}
		x, err := strconv.Atoi(xStr)
		if err != nil {
			return 1, fmt.Errorf("invalid input from player 1: %v", err)
		}
		if x < 0 {
			player2.AddScore(1)
			if verbose {
				fmt.Printf("Round %d: player 1 has surrendered, so player 2 wins.\n",
					c+1)
			}
			if err := player2.Send(strconv.Itoa(-1)); err != nil {
				return 2, fmt.Errorf("failed to send results to player 2: %v", err)
			}
			break
		}
		x = lessOrEqual(x, g.energy1)
		g.energy1 -= x
		total1 += x

		if x < y {
			player2.AddScore(1)
			if verbose {
				fmt.Printf("Round %d: player 1 lost by spending less energy (player 1: %d vs player 2: %d), so player 2 wins.\n",
					c+1, total1, total2)
			}
			if err := player1.Send(strconv.Itoa(-1)); err != nil {
				return 1, fmt.Errorf("failed to send results to player 1: %v", err)
			}
			if err := player2.Send(strconv.Itoa(-1)); err != nil {
				return 2, fmt.Errorf("failed to send results to player 2: %v", err)
			}
			break
		}

		x -= y

		if err := player2.Send(strconv.Itoa(x)); err != nil {
			return 2, fmt.Errorf("failed to send offer to player 2: %v", err)
		}
		yStr, err := player2.Receive(500 * time.Millisecond)
		if err != nil {
			return 2, fmt.Errorf("player 2 failed to respond: %v", err)
		}
		y, err = strconv.Atoi(yStr)
		if err != nil {
			return 2, fmt.Errorf("invalid input from player 2: %v", err)
		}

		if y < 0 {
			player1.AddScore(1)
			if verbose {
				fmt.Printf("Round %d: player 2 has surrendered, so player 1 wins.\n",
					c+1)
			}
			if err := player1.Send(strconv.Itoa(-1)); err != nil {
				return 1, fmt.Errorf("failed to send results to player 1: %v", err)
			}
			break
		}
		y = lessOrEqual(y, g.energy2)
		g.energy2 -= y
		total2 += y

		if y < x {
			player1.AddScore(1)
			if verbose {
				fmt.Printf("Round %d: player 2 lost by spending less energy (player 1: %d vs player 2: %d), so player 1 wins.\n",
					c+1, total1, total2)
			}
			if err := player1.Send(strconv.Itoa(-1)); err != nil {
				return 1, fmt.Errorf("failed to send results to player 1: %v", err)
			}
			if err := player2.Send(strconv.Itoa(-1)); err != nil {
				return 2, fmt.Errorf("failed to send results to player 2: %v", err)
			}
			break
		}

		if x == 0 && y == 0 {
			if verbose {
				fmt.Printf("Round %d: players prefer to save their energy and agreed on a draw.\n",
					c+1)
			}
			if err := player1.Send(strconv.Itoa(-1)); err != nil {
				return 1, fmt.Errorf("failed to send results to player 1: %v", err)
			}
			if err := player2.Send(strconv.Itoa(-1)); err != nil {
				return 2, fmt.Errorf("failed to send results to player 2: %v", err)
			}
			break
		}

		y -= x

		if verbose {
			fmt.Printf("Round %d: player 2 spent more or equal than player 1 (player 1: %d vs player 2: %d), so player 1 gets a chance to recoup.\n",
				c+1, total1, total2)
		}
	}

	if verbose {
		fmt.Printf("Round %d: fineshed, player 1 has %d energy, player 2 has %d energy.\n",
			c+1, g.energy1, g.energy2)
	}

	return 0, nil
}

func lessOrEqual(val, max int) int {
	if val > max {
		return max
	}
	return val
}
