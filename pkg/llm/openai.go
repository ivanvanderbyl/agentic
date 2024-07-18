package llm

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	options *Options
}

func NewOpenAI(opts ...Option) LLM {
	options := &Options{
		Model:     openai.GPT4o,
		MaxTokens: 4_000,
		APIKey:    os.Getenv("OPENAI_API_KEY"),
	}
	for _, opt := range opts {
		opt(options)
	}

	return &OpenAI{
		options: options,
	}
}

func (o *OpenAI) Generate(ctx context.Context, prompt string, opts ...Option) (string, error) {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}

	if o.options.APIKey == "" {
		return "", ErrNoAPIKey
	}

	cfg := openai.DefaultConfig(o.options.APIKey)

	slog.Info("Using cache", "useCache", options.UseCache)

	if options.UseCache {
		cfg.HTTPClient = &http.Client{
			Transport: NewCacheTransport(http.DefaultTransport, nil, "./cache", 24*time.Hour),
		}
	}

	client := openai.NewClientWithConfig(cfg)

	msgs := []openai.ChatCompletionMessage{}

	if options.SystemPrompt != "" {
		msgs = append(msgs, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: options.SystemPrompt,
		})
	}

	if prompt != "" {
		msgs = append(msgs, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		})
	}

	req := openai.ChatCompletionRequest{
		Model:       openai.GPT4o,
		MaxTokens:   options.MaxTokens,
		Messages:    msgs,
		Temperature: float32(options.Temperature),
		Stream:      false,
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}

	return resp.Choices[0].Message.Content, nil
}
