package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
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

	return lint(ctx)
}

func lint(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "golangci-lint", "run")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if ctx.Err() != nil {
		err = ctx.Err()
	}
	return err
}
