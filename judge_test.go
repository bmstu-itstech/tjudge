package main_test

import (
	"judge/game"
	"judge/player"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPlayer_CorrectWork(t *testing.T) {
	player, err := player.NewPlayer("tests/hello.py")
	require.NoError(t, err)
	err = player.StartGame()
	require.NoError(t, err)
	answer, err := player.Receive(500 * time.Millisecond)
	require.NoError(t, err)
	require.Equal(t, answer, "Hello!")
	err = player.StopGame()
	require.NoError(t, err)
}

func TestPlayer_FileNotExist(t *testing.T) {
	_, err := player.NewPlayer("notexist.py")
	require.ErrorContains(t, err, "file does not exist")
}

func TestPlayer_ItsDirectory(t *testing.T) {
	_, err := player.NewPlayer("tests")
	require.ErrorContains(t, err, "path is a directory")
}

func TestPlayer_FileIsEmpty(t *testing.T) {
	player, err := player.NewPlayer("tests/empty.py")
	require.NoError(t, err)
	err = player.StartGame()
	require.NoError(t, err)
	_, err = player.Receive(500 * time.Millisecond)
	require.ErrorContains(t, err, "stdout read: EOF")
}

func TestPlayer_TwiceRuned(t *testing.T) {
	player, err := player.NewPlayer("tests/hello.py")
	require.NoError(t, err)
	err = player.StartGame()
	require.NoError(t, err)
	err = player.StartGame()
	require.ErrorContains(t, err, "player is already running")
	err = player.StopGame()
	require.NoError(t, err)
}

func TestPlayer_TwiceStoped(t *testing.T) {
	player, err := player.NewPlayer("tests/hello.py")
	require.NoError(t, err)
	err = player.StartGame()
	require.NoError(t, err)
	err = player.StopGame()
	require.NoError(t, err)
	err = player.StopGame()
	require.ErrorContains(t, err, "player process is not running")
}

func TestGame_Unsupported(t *testing.T) {
	k, err := game.Play("not_existing_game", 10, nil, nil, false)
	require.Equal(t, k, 3)
	require.ErrorContains(t, err, "unsupported game")
}

func TestValidator_PrisonersDilemma(t *testing.T) {
	player1, err := player.NewPlayer("tests/prisoners_dilemma/kind.py")
	require.NoError(t, err)
	player2, err := player.NewPlayer("tests/prisoners_dilemma/evil.py")
	require.NoError(t, err)
	player3, err := player.NewPlayer("tests/hello.py")
	require.NoError(t, err)
	k, err := game.Play("prisoners_dilemma", 10, player1, player2, false)
	require.Equal(t, k, 0)
	require.NoError(t, err)
	k, err = game.Play("prisoners_dilemma", 10, player1, player3, false)
	require.Equal(t, k, 2)
	require.ErrorContains(t, err, "invalid choice")
	k, err = game.Play("prisoners_dilemma", 10, player2, player3, false)
	require.Equal(t, k, 2)
	require.ErrorContains(t, err, "invalid choice")
	require.Equal(t, player1.GetScore(), 0)
	require.Equal(t, player2.GetScore(), 50)
	require.Equal(t, player3.GetScore(), 0)
}

func TestGame_PrisonersDilemma(t *testing.T) {
	player1, err := player.NewPlayer("tests/prisoners_dilemma/kind.py")
	require.NoError(t, err)
	player2, err := player.NewPlayer("tests/prisoners_dilemma/evil.py")
	require.NoError(t, err)
	k, err := game.Play("prisoners_dilemma", 10, player1, player2, false)
	require.Equal(t, k, 0)
	require.NoError(t, err)
	require.Equal(t, player1.GetScore(), 0)
	require.Equal(t, player2.GetScore(), 50)
}

func TestGame_PrisonersDilemma_SilentPlayer(t *testing.T) {
	player1, err := player.NewPlayer("tests/echo.py")
	require.NoError(t, err)
	player2, err := player.NewPlayer("tests/prisoners_dilemma/evil.py")
	require.NoError(t, err)
	timer := time.Now()
	k, err := game.Play("prisoners_dilemma", 10, player1, player2, false)
	require.Equal(t, k, 1)
	require.ErrorContains(t, err, "timeout exceeded")
	require.Equal(t, time.Now().Compare(timer.Add(time.Second)), -1)
	require.Equal(t, player1.GetScore(), 0)
	require.Equal(t, player2.GetScore(), 0)
}

func TestGame_GoodDeal(t *testing.T) {
	player1, err := player.NewPlayer("tests/good_deal/kind.py")
	require.NoError(t, err)
	player2, err := player.NewPlayer("tests/good_deal/evil.py")
	require.NoError(t, err)
	player3, err := player.NewPlayer("tests/good_deal/smart.py")
	require.NoError(t, err)
	player4, err := player.NewPlayer("tests/good_deal/smart.py")
	require.NoError(t, err)
	k, err := game.Play("good_deal", 10, player1, player2, false)
	require.Equal(t, k, 0)
	require.NoError(t, err)
	k, err = game.Play("good_deal", 10, player1, player3, false)
	require.Equal(t, k, 0)
	require.NoError(t, err)
	k, err = game.Play("good_deal", 10, player3, player2, false)
	require.Equal(t, k, 0)
	require.NoError(t, err)
	k, err = game.Play("good_deal", 10, player3, player4, false)
	require.Equal(t, k, 0)
	require.NoError(t, err)
	require.Equal(t, player1.GetScore(), 0)
	require.Equal(t, player2.GetScore(), 500)
	require.Equal(t, player3.GetScore(), 1000)
	require.Equal(t, player4.GetScore(), 500)
}

func TestGame_GoodDeal_IncorrectValue(t *testing.T) {
	player1, err := player.NewPlayer("tests/echo.py")
	require.NoError(t, err)
	player2, err := player.NewPlayer("tests/good_deal/evil.py")
	require.NoError(t, err)
	k, err := game.Play("good_deal", 10, player2, player1, false)
	require.Equal(t, k, 2)
	require.ErrorContains(t, err, "invalid decision from B")
	require.Equal(t, player1.GetScore(), 0)
	require.Equal(t, player2.GetScore(), 0)
	player1, err = player.NewPlayer("tests/hello.py")
	require.NoError(t, err)
	k, err = game.Play("good_deal", 10, player1, player2, false)
	require.Equal(t, k, 1)
	require.ErrorContains(t, err, "invalid m from A")
	require.Equal(t, player1.GetScore(), 0)
	require.Equal(t, player2.GetScore(), 0)
}

func TestGame_GoodDeal_SilentPlayer(t *testing.T) {
	player1, err := player.NewPlayer("tests/silent.py")
	require.NoError(t, err)
	player2, err := player.NewPlayer("tests/good_deal/evil.py")
	require.NoError(t, err)
	timer := time.Now()
	k, err := game.Play("good_deal", 10, player1, player2, false)
	require.Equal(t, k, 1)
	require.ErrorContains(t, err, "timeout exceeded")
	require.Equal(t, time.Now().Compare(timer.Add(time.Second)), -1)
	require.Equal(t, player1.GetScore(), 0)
	require.Equal(t, player2.GetScore(), 0)
}

func TestGame_GoodDeal_TooBigAndTooLitleInput(t *testing.T) {
	player1, err := player.NewPlayer("tests/good_deal/tooBigInput.py")
	require.NoError(t, err)
	player2, err := player.NewPlayer("tests/good_deal/tooLitleInput.py")
	require.NoError(t, err)
	k, err := game.Play("good_deal", 10, player1, player2, false)
	require.Equal(t, k, 1)
	require.ErrorContains(t, err, "invalid m from A")
	require.Equal(t, player1.GetScore(), 0)
	require.Equal(t, player2.GetScore(), 0)
	k, err = game.Play("good_deal", 10, player2, player1, false)
	require.Equal(t, k, 1)
	require.ErrorContains(t, err, "invalid m from A")
	require.Equal(t, player1.GetScore(), 0)
	require.Equal(t, player2.GetScore(), 0)
}

func TestGame_TugOfWar(t *testing.T) {
	player1, err := player.NewPlayer("tests/tug_of_war/default.py")
	require.NoError(t, err)
	player2, err := player.NewPlayer("tests/tug_of_war/default2.py")
	require.NoError(t, err)
	player3, err := player.NewPlayer("tests/tug_of_war/always_lose.py")
	require.NoError(t, err)
	player4, err := player.NewPlayer("tests/tug_of_war/spender.py")
	require.NoError(t, err)
	k, err := game.Play("tug_of_war", 10, player4, player3, false)
	require.Equal(t, k, 0)
	require.NoError(t, err)
	require.Equal(t, player4.GetScore(), 1)
	require.Equal(t, player3.GetScore(), 9) // Транжира потратил всю энергию в первом раунде, то есть проиграл во всех последующих
	k, err = game.Play("tug_of_war", 10, player1, player3, false)
	require.Equal(t, k, 0)
	require.NoError(t, err)
	require.Equal(t, player1.GetScore(), 10) // Выиграл все 10 раундов
	require.Equal(t, player3.GetScore(), 9)
	k, err = game.Play("tug_of_war", 10, player1, player2, false)
	require.Equal(t, k, 0)
	require.NoError(t, err)
	require.Equal(t, player1.GetScore(), 11) // Выиграл 1 раунд,
	require.Equal(t, player2.GetScore(), 9)  // Выиграл 9 раундов - см. примечания к игре в режиме отладки
}

func TestGame_TugOfWar_IncorrectValue(t *testing.T) {
	player1, err := player.NewPlayer("tests/hello.py")
	require.NoError(t, err)
	player2, err := player.NewPlayer("tests/tug_of_war/default.py")
	require.NoError(t, err)
	k, err := game.Play("tug_of_war", 10, player1, player2, false)
	require.Equal(t, k, 1)
	require.ErrorContains(t, err, "invalid syntax")
	require.Equal(t, player1.GetScore(), 0)
	require.Equal(t, player2.GetScore(), 0)
}

func TestGame_TugOfWar_SilentPlayer(t *testing.T) {
	player1, err := player.NewPlayer("tests/silent.py")
	require.NoError(t, err)
	player2, err := player.NewPlayer("tests/tug_of_war/default.py")
	require.NoError(t, err)
	timer := time.Now()
	k, err := game.Play("tug_of_war", 10, player1, player2, false)
	require.Equal(t, k, 1)
	require.ErrorContains(t, err, "timeout exceeded")
	require.Equal(t, time.Now().Compare(timer.Add(time.Second)), -1)
	require.Equal(t, player1.GetScore(), 0)
	require.Equal(t, player2.GetScore(), 0)
}

func TestGame_TugOfWar_TooBigAndTooLitleInput(t *testing.T) {
	player1, err := player.NewPlayer("tests/tug_of_war/tooBigInput.py")
	require.NoError(t, err)
	player2, err := player.NewPlayer("tests/tug_of_war/tooLitleInput.py")
	require.NoError(t, err)
	k, err := game.Play("tug_of_war", 10, player1, player2, false)
	require.Equal(t, k, 0)
	require.NoError(t, err)
	require.Equal(t, player1.GetScore(), 1) // Потратил всю энергию в первом раунде, второй игрок ввёл отрицательное число, то есть сдался
	require.Equal(t, player2.GetScore(), 9) // Хотя он и вводит всегда отрицательные числа, по правилу он выиграл все оставшиеся раунды, так как у него остались силы
}

func TestGame_BalanceOfUniverse(t *testing.T) {
	player1, err := player.NewPlayer("tests/balance_of_universe/default.py")
	require.NoError(t, err)
	player2, err := player.NewPlayer("tests/balance_of_universe/default.py")
	require.NoError(t, err)
	player3, err := player.NewPlayer("tests/balance_of_universe/negative_value.py")
	require.NoError(t, err)
	player4, err := player.NewPlayer("tests/balance_of_universe/spender.py")
	require.NoError(t, err)
	k, err := game.Play("balance_of_universe", 10, player1, player2, false)
	require.Equal(t, k, 0)
	require.NoError(t, err)
	require.Equal(t, player1.GetScore(), 200)
	require.Equal(t, player2.GetScore(), 200)
	k, err = game.Play("balance_of_universe", 10, player1, player3, false)
	require.Equal(t, k, 2)
	require.ErrorContains(t, err, "invalid input")
	require.Equal(t, player1.GetScore(), 200)
	require.Equal(t, player3.GetScore(), 0)
	k, err = game.Play("balance_of_universe", 10, player1, player4, false)
	require.Equal(t, k, 0)
	require.NoError(t, err)
	require.Equal(t, player1.GetScore(), 299) // +100 за раунд, -1 в первом раунде, остальное на счёт
	require.Equal(t, player4.GetScore(), 101) // +100 за раунд, +1 за первый раунд, остальное "проиграл"
}

func TestGame_BalanceOfUniverse_IncorrectValue(t *testing.T) {
	player1, err := player.NewPlayer("tests/hello.py")
	require.NoError(t, err)
	player2, err := player.NewPlayer("tests/balance_of_universe/default.py")
	require.NoError(t, err)
	k, err := game.Play("balance_of_universe", 10, player1, player2, false)
	require.Equal(t, k, 1)
	require.ErrorContains(t, err, "invalid input")
	require.Equal(t, player1.GetScore(), 0)
	require.Equal(t, player2.GetScore(), 0)
}

func TestGame_BalanceOfUniverse_SilentPlayer(t *testing.T) {
	player1, err := player.NewPlayer("tests/silent.py")
	require.NoError(t, err)
	player2, err := player.NewPlayer("tests/balance_of_universe/default.py")
	require.NoError(t, err)
	timer := time.Now()
	k, err := game.Play("balance_of_universe", 10, player1, player2, false)
	require.Equal(t, k, 1)
	require.ErrorContains(t, err, "timeout exceeded")
	require.Equal(t, time.Now().Compare(timer.Add(time.Second)), -1)
	require.Equal(t, player1.GetScore(), 0)
	require.Equal(t, player2.GetScore(), 0)
}

func TestGame_BalanceOfUniverse_TooBigAndTooLitleInput(t *testing.T) {
	player1, err := player.NewPlayer("tests/balance_of_universe/tooBigInput.py")
	require.NoError(t, err)
	player2, err := player.NewPlayer("tests/balance_of_universe/tooLitleInput.py")
	require.NoError(t, err)
	k, err := game.Play("balance_of_universe", 10, player1, player2, false)
	require.Equal(t, k, 1)
	require.ErrorContains(t, err, "invalid input")
	require.Equal(t, player1.GetScore(), 0)
	require.Equal(t, player2.GetScore(), 0)
}
