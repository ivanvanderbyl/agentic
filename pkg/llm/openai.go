package llm

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	options *Options
}

var _ LLM = (*OpenAI)(nil)

func NewOpenAI(opts ...Option) LLM {
	options := &Options{
		Model:      openai.GPT4o,
		MaxTokens:  4_000,
		APIKey:     os.Getenv("OPENAI_API_KEY"),
		Dimensions: DefaultDimensions,
	}
	for _, opt := range opts {
		opt(options)
	}

	return &OpenAI{
		options: options,
	}
}

func (o *OpenAI) Generate(ctx context.Context, prompt string, opts ...Option) (string, error) {
	options := &Options{
		CacheDirectory: o.options.CacheDirectory,
		UseCache:       o.options.UseCache,
		Model:          o.options.Model,
		MaxTokens:      o.options.MaxTokens,
		Temperature:    o.options.Temperature,
		SystemPrompt:   o.options.SystemPrompt,
		APIKey:         o.options.APIKey,
	}
	for _, opt := range opts {
		opt(options)
	}

	if options.APIKey == "" {
		return "", ErrNoAPIKey
	}

	cfg := openai.DefaultConfig(options.APIKey)

	if options.UseCache {
		cfg.HTTPClient = &http.Client{
			Transport: NewCacheTransport(http.DefaultTransport, nil, options.CacheDirectory, 0),
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

func (o *OpenAI) Stream(ctx context.Context, prompt string, opts ...Option) (<-chan string, error) {
	options := &Options{
		CacheDirectory: o.options.CacheDirectory,
		UseCache:       o.options.UseCache,
		Model:          o.options.Model,
		MaxTokens:      o.options.MaxTokens,
		Temperature:    o.options.Temperature,
		SystemPrompt:   o.options.SystemPrompt,
		APIKey:         o.options.APIKey,
	}
	for _, opt := range opts {
		opt(options)
	}

	if options.APIKey == "" {
		return nil, ErrNoAPIKey
	}

	cfg := openai.DefaultConfig(options.APIKey)

	if options.UseCache {
		cfg.HTTPClient = &http.Client{
			Transport: NewCacheTransport(http.DefaultTransport, nil, options.CacheDirectory, 0),
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

	stream, err := client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:       openai.GPT4o,
		MaxTokens:   options.MaxTokens,
		Messages:    msgs,
		Temperature: float32(options.Temperature),
		Stream:      true,
	})
	if err != nil {
		return nil, err
	}
	recvChan := make(chan string)

	go func() {
		defer close(recvChan)
		defer stream.Close()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				resp, err := stream.Recv()
				if err != nil {
					return
				}

				if len(resp.Choices) == 0 {
					return
				}

				recvChan <- resp.Choices[0].Delta.Content
			}
		}
	}()

	return recvChan, nil
}

func (o *OpenAI) Embedding(ctx context.Context, input string, opts ...Option) ([]float32, error) {
	options := &Options{
		CacheDirectory: o.options.CacheDirectory,
		UseCache:       o.options.UseCache,
		Model:          o.options.Model,
		MaxTokens:      o.options.MaxTokens,
		Temperature:    o.options.Temperature,
		SystemPrompt:   o.options.SystemPrompt,
		APIKey:         o.options.APIKey,
		Dimensions:     o.options.Dimensions,
	}
	for _, opt := range opts {
		opt(options)
	}

	if options.APIKey == "" {
		return nil, ErrNoAPIKey
	}

	cfg := openai.DefaultConfig(options.APIKey)

	if options.UseCache {
		cfg.HTTPClient = &http.Client{
			Transport: NewCacheTransport(http.DefaultTransport, nil, options.CacheDirectory, 0),
		}
	}

	client := openai.NewClientWithConfig(cfg)

	req := openai.EmbeddingRequest{
		Model:          openai.SmallEmbedding3,
		EncodingFormat: openai.EmbeddingEncodingFormatFloat,
		Dimensions:     options.Dimensions,
		Input:          input,
	}

	resp, err := client.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return resp.Data[0].Embedding, nil
}
