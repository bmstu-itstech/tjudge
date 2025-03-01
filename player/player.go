package player

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	ErrPlayerNotRunning = errors.New("player process is not running")
	ErrTimeout          = errors.New("timeout exceeded")
)

type Player struct {
	path      string
	cmd       *exec.Cmd
	stdin     *bufio.Writer
	stdout    *bufio.Reader
	isRunning bool
	score     int
}

func checkFileExists(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", path)
		}
		return fmt.Errorf("failed to check file: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a file: %s", path)
	}

	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("file is not accessible: %w", err)
	}
	file.Close()

	return nil
}

func NewPlayer(path string) (*Player, error) {
	if path == "" {
		return nil, errors.New("path to executable is required")
	}
	if err := checkFileExists(path); err != nil {
		return nil, err
	}
	return &Player{path: path}, nil
}

func (p *Player) StartGame() error {
	if p.isRunning {
		return errors.New("player is already running")
	}

	cmd := exec.Command(p.path)

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("stdin pipe: %w", err)
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %w", err)
	}

	p.cmd = cmd
	p.stdin = bufio.NewWriter(stdinPipe)
	p.stdout = bufio.NewReader(stdoutPipe)
	p.isRunning = true

	if err := cmd.Start(); err != nil {
		p.isRunning = false
		return fmt.Errorf("process start: %w", err)
	}

	return nil
}

func (p *Player) Send(data string) error {
	if !p.isRunning {
		return ErrPlayerNotRunning
	}

	if _, err := p.stdin.WriteString(data + "\n"); err != nil {
		return fmt.Errorf("stdin write: %w", err)
	}
	return p.stdin.Flush()
}

func (p *Player) Receive(timeout time.Duration) (string, error) {
	if !p.isRunning {
		return "", ErrPlayerNotRunning
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result := make(chan string)
	errChan := make(chan error)

	go func() {
		select {
		case <-ctx.Done():
			return
		default:
			response, err := p.stdout.ReadString('\n')
			if err != nil {
				errChan <- err
				return
			}
			result <- strings.TrimSpace(response)
		}
	}()

	select {
	case res := <-result:
		return res, nil
	case err := <-errChan:
		return "", fmt.Errorf("stdout read: %w", err)
	case <-ctx.Done():
		return "", ErrTimeout
	}
}

func (p *Player) StopGame() error {
	if !p.isRunning {
		return ErrPlayerNotRunning
	}

	if err := p.cmd.Process.Kill(); err != nil {
		return fmt.Errorf("process kill: %w", err)
	}

	p.isRunning = false
	return nil
}

func (p *Player) AddScore(i int) {
	p.score += i
}

func (p *Player) GetScore() int {
	return p.score
}
