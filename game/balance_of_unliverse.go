package game

import (
	"fmt"
	"strconv"
	"time"

	"judge/player"
)

type BalanceOfUniverse struct {
	n            int
	used1, used2 int
	round        int
}

func NewBalanceOfUniverse(n int) *BalanceOfUniverse {
	return &BalanceOfUniverse{
		n:     n,
		used1: 0,
		used2: 0,
		round: 0,
	}
}

func (g *BalanceOfUniverse) Name() string {
	return "balance_of_universe"
}

func (g *BalanceOfUniverse) validate(output string) error {
	return nil
}

func (g *BalanceOfUniverse) playRound(c int, player1, player2 *player.Player, verbose bool) (int, error) {
	if c > 0 {
		return 0, nil
	}
	for {
		if g.round == 0 {
			if err := player1.Send(strconv.Itoa(g.n)); err != nil {
				return 1, fmt.Errorf("failed to send n to player 1: %v", err)
			}
			if err := player2.Send(strconv.Itoa(g.n)); err != nil {
				return 2, fmt.Errorf("failed to send n to player 2: %v", err)
			}
		}
		g.round++

		if g.used1 == g.n && g.used2 == g.n {
			if verbose {
				fmt.Printf("both players have used all their bullions\n")
			}
			break
		}

		aStr, err := player1.Receive(500 * time.Millisecond)
		if err != nil {
			return 1, fmt.Errorf("player 1 failed to respond: %v", err)
		}
		a, err := strconv.Atoi(aStr)
		if err != nil || a < 0 || a+g.used1 > g.n {
			return 1, fmt.Errorf("invalid input from player 1: %s", aStr)
		}

		bStr, err := player2.Receive(500 * time.Millisecond)
		if err != nil {
			return 2, fmt.Errorf("player 2 failed to respond: %v", err)
		}
		b, err := strconv.Atoi(bStr)
		if err != nil || b < 0 || b+g.used2 > g.n {
			return 2, fmt.Errorf("invalid input from player 2: %s", bStr)
		}

		if err := player1.Send(strconv.Itoa(b)); err != nil {
			return 1, fmt.Errorf("failed to send choise of player 2 to player 1: %v", err)
		}
		if err := player2.Send(strconv.Itoa(a)); err != nil {
			return 2, fmt.Errorf("failed to send choise of player 1 to player 2: %v", err)
		}

		g.used1 += a
		g.used2 += b

		if verbose {
			fmt.Printf("Round %d: player 1 put %d (used %d/%d), player 2 put %d (used %d/%d). ",
				g.round, a, g.used1, g.n, b, g.used2, g.n)
		}
		if abs(a-b) <= 1 {
			player1.AddScore(a + b)
			player2.AddScore(a + b)
			if verbose {
				fmt.Printf("Balance! Both get %d\n", a+b)
			}
		} else {
			if a > b {
				player1.AddScore(a + b)
				if verbose {
					fmt.Printf("Imbalance, player 1 gets %d\n", a+b)
				}
			} else {
				player2.AddScore(a + b)
				if verbose {
					fmt.Printf("Imbalance, player gets %d\n", a+b)
				}
			}
		}
	}

	return 0, nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
