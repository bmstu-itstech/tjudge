package game

import (
	"fmt"
	"judge/player"
)

type PrisonersDilemma struct{}

func NewPrisonersDilemma() *PrisonersDilemma {
	return &PrisonersDilemma{}
}

func (g *PrisonersDilemma) Name() string {
	return "prisoners_dilemma"
}

func (g *PrisonersDilemma) playRound(c, i, j int, player1, player2 *player.Player, verbose bool) (int, error) {
	if verbose {
		fmt.Printf("Iteration %d:\n", c+1)
	}

	choice1, err := getPlayerChoice(player1, verbose)
	if err != nil {
		return i, fmt.Errorf("player %d error: %v", i, err)
	}

	choice2, err := getPlayerChoice(player2, verbose)
	if err != nil {
		return j, fmt.Errorf("player %d error: %v", j, err)
	}

	if err := g.validate(choice1); err != nil {
		return i, fmt.Errorf("player %d invalid choice: %v", i, err)
	}
	if err := g.validate(choice2); err != nil {
		return j, fmt.Errorf("player %d invalid choice: %v", j, err)
	}

	if err := player1.Send(choice2); err != nil {
		return i, fmt.Errorf("failed to send input to player %d: %v", i, err)
	}

	if err := player2.Send(choice1); err != nil {
		return j, fmt.Errorf("failed to send input to player %d: %v", j, err)
	}

	s1, s2 := calculateScoresPrisonersDilemma(choice1, choice2)
	player1.AddScore(s1)
	player2.AddScore(s2)

	if verbose {
		fmt.Printf("Player %d choice: %s, Player %d choice: %s\n", i, choice1, j, choice2)
		fmt.Printf("Scores after iteration %d: Player %d: %d, Player %d: %d\n", c, i, player1.GetScore(), j, player2.GetScore())
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
