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
	prompt := `Return relationship_keyword: a single word in UPPERCASE to describe the relationship between the source entity and target entity, e.g. "FRIENDSHIP", "RIVALRY", "COLLABORATION", "SUPPORTS", "OPPOSES", "WORKS_IN", "MEMBER_OF", "AFFECTED_BY"
	Input text: The team is directly involved in Operation: Dulce, executing its evolved objectives and activities.`
	completion, err := l.Generate(ctx, prompt, llm.WithCache(".cache"))
	if err != nil {
		return err
	}

	slog.Info("Generated completion", "completion", completion)

	return nil
}
