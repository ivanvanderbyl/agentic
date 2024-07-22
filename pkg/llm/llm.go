package llm

import (
	"context"
	"fmt"
)

type (
	LLM interface {
		Generate(ctx context.Context, prompt string, opts ...Option) (string, error)
		Embedding(ctx context.Context, input string, opts ...Option) ([]float32, error)
	}

	LLMHistory interface {
		History() []string
	}

	Option func(*Options)

	Options struct {
		APIKey         string  // API Key for underlying service
		MaxTokens      int     // Max tokens to generate when generating text
		Dimensions     int     // Embedding dimensions to generate when embedding
		Model          string  // Model to use
		SystemPrompt   string  // System prompt for completion
		Temperature    float64 // Temperature for sampling
		UseCache       bool    // Enable HTTP Request caching
		CacheDirectory string  // Directory to store cache
	}
)

const (
	DefaultDimensions = 1536 // Default dimensions for embeddings based on OpenAI's embedding small model
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

// WithDimensions sets the dimensions
func WithDimensions(dimensions int) Option {
	return func(o *Options) {
		o.Dimensions = dimensions
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
func WithCache(dir string) Option {
	return func(o *Options) {
		o.UseCache = true
		o.CacheDirectory = dir
	}
}
