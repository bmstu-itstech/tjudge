package main_test

import (
	"judge/game"
	"judge/player"
	"judge/validator"
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

func TestValidator(t *testing.T) {
	err := validator.Validate("Y", "prisoners_dilemma")
	require.NoError(t, err)
	err = validator.Validate("N", "prisoners_dilemma")
	require.NoError(t, err)
	err = validator.Validate("ABC", "prisoners_dilemma")
	require.ErrorContains(t, err, "invalid choice")
	err = validator.Validate("ABC", "not_existing_game")
	require.ErrorContains(t, err, "unsupported game")
}

func TestGame_PrisonersDilemma(t *testing.T) {
	player1, err := player.NewPlayer("tests/prisoners_dilemma/kind.py")
	require.NoError(t, err)
	player2, err := player.NewPlayer("tests/prisoners_dilemma/evil.py")
	require.NoError(t, err)
	players := []*player.Player{player1, player2}
	game, err := game.NewGame("prisoners_dilemma", &players)
	require.NoError(t, err)
	err = game.Play(10, false)
	require.NoError(t, err)
	require.Equal(t, player1.GetScore(), 0)
	require.Equal(t, player2.GetScore(), 100)
}

func TestGame_PrisonersDilemma_SilentPlayer(t *testing.T) {
	player1, err := player.NewPlayer("tests/echo.py")
	require.NoError(t, err)
	player2, err := player.NewPlayer("tests/prisoners_dilemma/evil.py")
	require.NoError(t, err)
	players := []*player.Player{player1, player2}
	game, err := game.NewGame("prisoners_dilemma", &players)
	require.NoError(t, err)
	timer := time.Now()
	err = game.Play(10, false)
	require.ErrorContains(t, err, "failed to get choice from player: timeout exceeded")
	require.Equal(t, time.Now().Compare(timer.Add(time.Second)), -1)
	require.Equal(t, player1.GetScore(), 0)
	require.Equal(t, player2.GetScore(), 0)
}

func TestGame_Unsupported(t *testing.T) {
	_, err := game.NewGame("not_existing_game", &[]*player.Player{})
	require.ErrorContains(t, err, "unsupported game")
}

func TestGame_PrisonersDilemma_ManyPlayers(t *testing.T) {
	player1, err := player.NewPlayer("tests/echo.py")
	require.NoError(t, err)
	player2, err := player.NewPlayer("tests/prisoners_dilemma/evil.py")
	require.NoError(t, err)
	player3, err := player.NewPlayer("tests/echo.py")
	require.NoError(t, err)
	player4, err := player.NewPlayer("tests/prisoners_dilemma/kind.py")
	require.NoError(t, err)
	player5, err := player.NewPlayer("tests/prisoners_dilemma/kind.py")
	require.NoError(t, err)
	player6, err := player.NewPlayer("tests/prisoners_dilemma/tit_for_tat.py")
	require.NoError(t, err)
	player7, err := player.NewPlayer("tests/prisoners_dilemma/tit_for_tat.py")
	require.NoError(t, err)
	players := []*player.Player{player1, player2, player3, player4, player5, player6, player7}
	game, err := game.NewGame("prisoners_dilemma", &players)
	require.NoError(t, err)
	timer := time.Now()
	err = game.Play(10, false)
	require.ErrorContains(t, err, "failed to get choice from player: timeout exceeded")
	require.Equal(t, time.Now().Compare(timer.Add(time.Second*2)), -1)
	require.Equal(t, player1.GetScore(), 0)
	require.Equal(t, player2.GetScore(), 256)
	require.Equal(t, player3.GetScore(), 0)
	require.Equal(t, player4.GetScore(), 180)
	require.Equal(t, player5.GetScore(), 180)
	require.Equal(t, player6.GetScore(), 198)
	require.Equal(t, player7.GetScore(), 198)
}
