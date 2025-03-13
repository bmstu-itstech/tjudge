package game

import (
	"fmt"
	"judge/player"
	"time"
)

type PrisonersDilemma struct{}

func NewPrisonersDilemma() *PrisonersDilemma {
	return &PrisonersDilemma{}
}

func (g *PrisonersDilemma) Name() string {
	return "prisoners_dilemma"
}

func (g *PrisonersDilemma) playRound(c int, player1, player2 *player.Player, verbose bool) (int, error) {
	if verbose {
		fmt.Printf("Iteration %d:\n", c+1)
	}

	choice1, err := player1.Receive(500 * time.Millisecond)
	if err != nil {
		return 1, fmt.Errorf("failed to get choice from player 1: %v", err)
	}
	if verbose {
		fmt.Printf("Player 1 choice: %s\n", choice1)
	}

	choice2, err := player2.Receive(500 * time.Millisecond)
	if err != nil {
		return 2, fmt.Errorf("failed to get choice from player 2: %v", err)
	}
	if verbose {
		fmt.Printf("Player 2 choice: %s\n", choice1)
	}

	if err := g.validate(choice1); err != nil {
		return 1, fmt.Errorf("player %d invalid choice: %v", 1, err)
	}
	if err := g.validate(choice2); err != nil {
		return 2, fmt.Errorf("player %d invalid choice: %v", 2, err)
	}

	if err := player1.Send(choice2); err != nil {
		return 1, fmt.Errorf("failed to send input to player %d: %v", 1, err)
	}

	if err := player2.Send(choice1); err != nil {
		return 2, fmt.Errorf("failed to send input to player %d: %v", 2, err)
	}

	s1, s2 := calculateScoresPrisonersDilemma(choice1, choice2)
	player1.AddScore(s1)
	player2.AddScore(s2)

	if verbose {
		fmt.Printf("Player %d choice: %s, Player %d choice: %s\n", 1, choice1, 2, choice2)
		fmt.Printf("Scores after iteration %d: Player %d: %d, Player %d: %d\n", c, 1, player1.GetScore(), 2, player2.GetScore())
	}
	return -1, nil
}

func calculateScoresPrisonersDilemma(choice1, choice2 string) (int, int) {
	if choice1 == "Y" && choice2 == "Y" {
		return 3, 3
	} else if choice1 == "Y" && choice2 == "N" {
		return 0, 5
	} else if choice1 == "N" && choice2 == "Y" {
		return 5, 0
	} else if choice1 == "N" && choice2 == "N" {
		return 1, 1
	}
	return 0, 0
}

func (g PrisonersDilemma) validate(output string) error {
	validChoices := map[string]struct{}{
		"Y": {},
		"N": {},
	}

	if _, ok := validChoices[output]; !ok {
		return fmt.Errorf(
			"invalid choice '%s', expected Y or N",
			output,
		)
	}
	return nil
}
