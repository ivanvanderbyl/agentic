package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ivanvanderbyl/graphrag-go/pkg/llm"
)

var usage = `createvector is a simple tool to return a vector representation of a given input from a language model.
It defaults to OpenAI's embedding-small-3 model with 1536 vectors.

createvector <input string>
`

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := realMain(ctx); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func realMain(ctx context.Context) error {
	if len(os.Args) < 2 {
		return fmt.Errorf(usage)
	}

	input := os.Args[1]
	model := llm.NewOpenAI(llm.WithCache(".cache"))
	embedding, err := model.Embedding(ctx, input)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", embedding)
	return nil
}
