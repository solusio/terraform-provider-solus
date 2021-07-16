package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	fmt.Println("Run pre-commit hook ...")
	if err := run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	ff, err := getChangedFiles(ctx)
	if err != nil {
		return fmt.Errorf("failed to get changed files: %w", err)
	}

	return lint(ctx, ff)
}

func getChangedFiles(ctx context.Context) ([]string, error) {
	cmd := exec.CommandContext(ctx, "git", "diff", "--cached", "--name-only")
	p, err := cmd.Output()
	if ctx.Err() != nil {
		err = ctx.Err()
	}
	if err != nil {
		return nil, err
	}

	ff := make([]string, 0, bytes.Count(p, []byte("\n")))
	s := bufio.NewScanner(bytes.NewReader(p))
	for s.Scan() {
		ff = append(ff, s.Text())
	}

	return ff, s.Err()
}

func lint(ctx context.Context, ff []string) error {
	// golangci-lint doesn't allow to lint files in separate directories, so we
	// have to group all files by folders.
	groups := map[string][]string{}

	for _, f := range ff {
		group := filepath.Dir(f)
		g, ok := groups[group]
		if !ok {
			g = []string{}
		}
		g = append(g, f)
		groups[group] = g
	}

	for _, ff := range groups {
		fmt.Printf("Lint files %v\n\n", ff)
		if err := runLinter(ctx, ff); err != nil {
			if isPrintableError(err) {
				fmt.Printf("Failed to lint files %v: %s\n", ff, err)
			}
		}
	}

	return nil
}

func runLinter(ctx context.Context, ff []string) error {
	args := append([]string{"run"}, ff...)
	cmd := exec.CommandContext(ctx, "golangci-lint", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if ctx.Err() != nil {
		err = ctx.Err()
	}
	return err
}

func isPrintableError(err error) bool {
	return errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)
}
