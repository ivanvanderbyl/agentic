package entity

import (
	"context"
	"fmt"

	"github.com/ivanvanderbyl/graphrag-go/pkg/llm"
	"github.com/ivanvanderbyl/graphrag-go/pkg/prompts"
)

type (
	EntityExtractor struct {
		llm              llm.LLM
		ExtractionPrompt string
		MaxGleanings     int
	}

	Option func(*EntityExtractor)

	Data struct {
		prompts.PromptData
		EntityTypes []string
		InputText   string
	}
)

func NewEntityExtractor(llm llm.LLM, opts ...Option) *EntityExtractor {
	e := &EntityExtractor{
		llm: llm,
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (ee *EntityExtractor) Extract(ctx context.Context, text string) error {
	tmpl, err := prompts.RenderTemplate(prompts.EntitiesTemplate, Data{
		EntityTypes: []string{"organization", "person", "policy", "bill", "geo", "event", "role", "electorate"},
		PromptData:  prompts.DefaultPromptData,
		InputText:   text,
	})
	if err != nil {
		return err
	}

	ee.ExtractionPrompt = tmpl

	resp, err := ee.llm.Generate(ctx, ee.ExtractionPrompt)
	if err != nil {
		return err
	}

	fmt.Println(resp)

	return nil
}
