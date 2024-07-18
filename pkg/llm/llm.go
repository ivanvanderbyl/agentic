package llm

import (
	"context"
	"fmt"
)

type (
	LLM interface {
		Generate(ctx context.Context, prompt string, opts ...Option) (string, error)
	}

	Option func(*Options)

	Options struct {
		APIKey       string
		MaxTokens    int
		Model        string
		SystemPrompt string
		Temperature  float64
		UseCache     bool
	}
)

var ErrNoAPIKey = fmt.Errorf("API key is required")
var ErrNoModel = fmt.Errorf("Model is required")

func WithAPIKey(apiKey string) Option {
	return func(o *Options) {
		o.APIKey = apiKey
	}
}

// WithMaxTokens sets the max tokens
func WithMaxTokens(maxTokens int) Option {
	return func(o *Options) {
		o.MaxTokens = maxTokens
	}
}

// WithModel sets the model
func WithModel(model string) Option {
	return func(o *Options) {
		o.Model = model
	}
}

// WithSystemPrompt sets the system prompt
func WithSystemPrompt(systemPrompt string) Option {
	return func(o *Options) {
		o.SystemPrompt = systemPrompt
	}
}

// WithTemperature sets the temperature
func WithTemperature(temperature float64) Option {
	return func(o *Options) {
		o.Temperature = temperature
	}
}

// WithCache sets the cache
func WithCache() Option {
	return func(o *Options) {
		o.UseCache = true
	}
}
