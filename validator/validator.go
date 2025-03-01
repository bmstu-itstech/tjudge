package validator

import (
	"fmt"
	"strings"
)

func Validate(output string, game string) error {
	normalized := strings.TrimSpace(output)

	switch game {
	case "prisoners_dilemma":
		return validatePrisonersDilemma(normalized)
	default:
		return fmt.Errorf("unsupported game: %s", game)
	}
}

// Валидатор для Дилеммы Заключенного
func validatePrisonersDilemma(output string) error {
	validChoices := map[string]bool{
		"Y": true,
		"N": true,
	}

	if !validChoices[output] {
		return fmt.Errorf(
			"invalid choice '%s', expected Y or N",
			output,
		)
	}
	return nil
}
