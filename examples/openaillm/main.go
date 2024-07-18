package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ivanvanderbyl/graphrag-go/pkg/llm"
)

func main() {
	err := realMain(context.Background())
	if err != nil {
		slog.Error("Command Failed", "error", err)
		os.Exit(1)
	}
}

func realMain(ctx context.Context) error {
	l := llm.NewOpenAI()
	prompt := "A claim that is not supported by evidence"
	completion, err := l.Generate(ctx, prompt, llm.WithCache(".cache"))
	if err != nil {
		return err
	}

	slog.Info("Generated completion", "completion", completion)

	return nil
}
