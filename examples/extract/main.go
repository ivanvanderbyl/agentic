package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/ivanvanderbyl/graphrag-go/pkg/extractors/entity"
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
	l := llm.NewOpenAI(llm.WithCache(".cache"))

	doc, err := os.ReadFile(os.Args[1])
	if err != nil {
		return err
	}

	extractor := entity.NewEntityExtractor(l)
	records, err := extractor.Extract(ctx, string(doc))
	if err != nil {
		return err
	}

	for _, record := range records {
		fmt.Println(record)
	}

	return nil
}
