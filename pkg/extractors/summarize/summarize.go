package summarize

import (
	"bytes"
	"context"
	"encoding/json"
	"slices"
	"strings"

	"github.com/ivanvanderbyl/graphrag-go/pkg/llm"
	"github.com/pkoukk/tiktoken-go"
)

// Max token size for input prompts
const DEFAULT_MAX_INPUT_TOKENS = 4_000

// Max token count for LLM answers
const DEFAULT_MAX_SUMMARY_LENGTH = 500

type (
	SummarizeExtractor struct {
		LLM                  llm.LLM
		entityNameKey        string
		inputDescriptionsKey string
		summarizationPrompt  string
		maxSummaryLength     int
		maxInputTokens       int
	}

	SummarizationResult struct {
		Items       []string
		Description string
	}

	Option func(*SummarizeExtractor)
)

// WithEntityNameKey sets the entity name key
func WithEntityNameKey(entityNameKey string) Option {
	return func(se *SummarizeExtractor) {
		se.entityNameKey = entityNameKey
	}
}

// WithInputDescriptionsKey sets the input descriptions key
func WithInputDescriptionsKey(inputDescriptionsKey string) Option {
	return func(se *SummarizeExtractor) {
		se.inputDescriptionsKey = inputDescriptionsKey
	}
}

// WithSummarizationPrompt sets the summarization prompt
func WithSummarizationPrompt(summarizationPrompt string) Option {
	return func(se *SummarizeExtractor) {
		se.summarizationPrompt = summarizationPrompt
	}
}

// WithMaxSummaryLength sets the max summary length
func WithMaxSummaryLength(maxSummaryLength int) Option {
	return func(se *SummarizeExtractor) {
		se.maxSummaryLength = maxSummaryLength
	}
}

// WithMaxInputTokens sets the max input tokens
func WithMaxInputTokens(maxInputTokens int) Option {
	return func(se *SummarizeExtractor) {
		se.maxInputTokens = maxInputTokens
	}
}

// NewSummarizeExtractor creates a new SummarizeExtractor
func NewSummarizeExtractor(llm llm.LLM, opts ...Option) *SummarizeExtractor {
	s := &SummarizeExtractor{
		LLM:                  llm,
		maxSummaryLength:     DEFAULT_MAX_SUMMARY_LENGTH,
		maxInputTokens:       DEFAULT_MAX_INPUT_TOKENS,
		summarizationPrompt:  prompt,
		entityNameKey:        "entity_name",
		inputDescriptionsKey: "description_list",
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Summarize extracts a summary from the given entity and descriptions returning a single summary
func (se *SummarizeExtractor) Summarize(ctx context.Context, entities []string, descriptions []string) (*SummarizationResult, error) {
	var err error

	result := ""
	if len(descriptions) == 0 {
		result = ""
	} else if len(descriptions) == 1 {
		result = descriptions[0]
	} else {
		result, err = se.summarizeDescriptions(ctx, entities, descriptions)
		if err != nil {
			return nil, err
		}
	}

	return &SummarizationResult{
		Items:       entities,
		Description: result,
	}, nil
}

// summarizeDescriptions Summarize descriptions into a single description
func (se *SummarizeExtractor) summarizeDescriptions(ctx context.Context, entities, descriptions []string) (string, error) {
	slices.Sort(descriptions)
	slices.Sort(entities)

	// Iterate over descriptions, adding all until the max input tokens is reached
	summarizationPromptTokens, err := se.numTokensFromString(se.summarizationPrompt)
	if err != nil {
		return "", err
	}
	usableTokens := se.maxInputTokens - summarizationPromptTokens
	descriptionsCollected := make([]string, 0, len(descriptions))
	result := ""

	for i, description := range descriptions {
		inputTokens, err := se.numTokensFromString(description)
		if err != nil {
			return "", err
		}

		usableTokens -= inputTokens
		descriptionsCollected = append(descriptionsCollected, description)

		isLastDescription := i == len(descriptions)-1
		bufferIsFull := usableTokens < 0 && len(descriptionsCollected) > 0

		if isLastDescription || bufferIsFull {
			result, err = se.summarizeWithLLM(ctx, entities, descriptionsCollected)
			if err != nil {
				return "", err
			}

			if i != len(descriptions)-1 {
				resultTokens, err := se.numTokensFromString(result)
				if err != nil {
					return "", err
				}

				descriptionsCollected = []string{result}
				usableTokens = se.maxInputTokens - summarizationPromptTokens - resultTokens
			}
		}
	}

	return result, nil
}

func (se *SummarizeExtractor) summarizeWithLLM(ctx context.Context, entities, descriptions []string) (string, error) {
	input := se.summarizationPrompt
	input = strings.ReplaceAll(input, "{entity_name}", jsonStrings(entities))
	input = strings.ReplaceAll(input, "{description_list}", jsonStrings(descriptions))
	return se.LLM.Generate(ctx, input, llm.WithMaxTokens(se.maxSummaryLength))
}

func (se *SummarizeExtractor) numTokensFromString(input string) (int, error) {
	t, err := tiktoken.GetEncoding(tiktoken.MODEL_CL100K_BASE)
	if err != nil {
		return 0, err
	}

	result := t.Encode(input, nil, nil)
	return len(result), nil
}

func jsonStrings(slice []string) string {
	result := bytes.NewBuffer(nil)
	_ = json.NewEncoder(result).Encode(slice)
	return result.String()
}
