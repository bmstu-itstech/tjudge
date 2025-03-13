package game

import (
	"fmt"
	"strconv"
	"time"

	"judge/player"
)

type GoodDeal struct {
	R, T int
}

func NewGoodDeal(R, T int) *GoodDeal {
	return &GoodDeal{
		R: R,
		T: T,
	}
}

func (g *GoodDeal) Name() string {
	return "good_deal"
}

func (g *GoodDeal) playRound(c int, player1, player2 *player.Player, verbose bool) (int, error) {
	var aPlayer, bPlayer *player.Player
	if c%2 == 0 {
		aPlayer = player1
		bPlayer = player2
	} else {
		aPlayer = player2
		bPlayer = player1
	}

	if err := aPlayer.Send(strconv.Itoa(g.R)); err != nil {
		return c%2 + 1, fmt.Errorf("failed to send R to A: %v", err)
	}
	if err := aPlayer.Send(strconv.Itoa(g.T)); err != nil {
		return c%2 + 1, fmt.Errorf("failed to send T to A: %v", err)
	}
	if err := aPlayer.Send("A"); err != nil {
		return c%2 + 1, fmt.Errorf("failed to send role to A: %v", err)
	}
	mStr, err := aPlayer.Receive(500 * time.Millisecond)
	if err != nil {
		return c%2 + 1, fmt.Errorf("failed to receive m from A: %v", err)
	}
	m, err := strconv.Atoi(mStr)
	if err != nil || m < 0 || m > g.R {
		return c%2 + 1, fmt.Errorf("invalid m from A: %s", mStr)
	}

	if err := bPlayer.Send(strconv.Itoa(g.R)); err != nil {
		return (c+1)%2 + 1, fmt.Errorf("failed to send R to B: %v", err)
	}
	if err := bPlayer.Send(strconv.Itoa(g.T)); err != nil {
		return (c+1)%2 + 1, fmt.Errorf("failed to send T to B: %v", err)
	}
	if err := bPlayer.Send(fmt.Sprintf("B\n%d", m)); err != nil {
		return (c+1)%2 + 1, fmt.Errorf("failed to send role and m to B: %v", err)
	}
	decisionStr, err := bPlayer.Receive(500 * time.Millisecond)
	if err != nil {
		return (c+1)%2 + 1, fmt.Errorf("failed to receive decision from B: %v", err)
	}
	decision, err := strconv.Atoi(decisionStr)
	if err != nil {
		return (c+1)%2 + 1, fmt.Errorf("invalid decision from B: %v", err)
	}
	if g.validate(decisionStr) != nil {
		return (c+1)%2 + 1, fmt.Errorf("invalid decision from B: %s, %v", decisionStr, g.validate(decisionStr))
	}

	if err := aPlayer.Send(decisionStr); err != nil {
		return c%2 + 1, fmt.Errorf("failed to send answer from B to A: %v", err)
	}

	if decision == 1 {
		if aPlayer == player1 {
			player1.AddScore(g.R - m)
			player2.AddScore(m)
		} else {
			player1.AddScore(m)
			player2.AddScore(g.R - m)
		}
	} else {
		player1.AddScore(-g.T)
		player2.AddScore(-g.T)
	}

	if verbose {
		fmt.Printf("Round %d: A offered %d, B decided %d\n",
			c+1, m, decision)
	}

	return 0, nil
}

func (g GoodDeal) validate(output string) error {
	validChoices := map[string]struct{}{
		"1": {},
		"0": {},
	}

	if _, ok := validChoices[output]; !ok {
		return fmt.Errorf(
			"invalid choice '%s', expected 1 or 0",
			output,
		)
	}
	return nil
}
