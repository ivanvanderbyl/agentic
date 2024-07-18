package entity

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

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

	Record interface {
		isNode()
		Type() string
	}

	Entity struct {
		Name         string
		Description  string
		internalType string
	}

	Relationship struct {
		Entity1  string
		Relation string
		Entity2  string
	}
)

func (Entity) isNode() {}
func (e *Entity) Type() string {
	return e.internalType
}
func (e *Entity) String() string {
	return fmt.Sprintf("Entity{Type: %q, Name: %q, Description: %q}", e.internalType, e.Name, e.Description)
}

func (Relationship) isNode() {}
func (r *Relationship) Type() string {
	return "relationship"
}
func (r *Relationship) String() string {
	return fmt.Sprintf("Relationship{From: %q, Relation: %q, To: %q}", r.Entity1, r.Relation, r.Entity2)
}

func NewEntityExtractor(llm llm.LLM, opts ...Option) *EntityExtractor {
	e := &EntityExtractor{
		llm: llm,
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (ee *EntityExtractor) Extract(ctx context.Context, text string) ([]Record, error) {
	prompt, err := prompts.RenderTemplate(prompts.EntitiesTemplate, Data{
		EntityTypes: []string{"organization", "person", "policy", "bill", "geo", "event", "role", "electorate"},
		PromptData:  prompts.DefaultPromptData,
		InputText:   text,
	})
	if err != nil {
		return nil, err
	}

	resp, err := ee.llm.Generate(ctx, prompt)
	if err != nil {
		return nil, err
	}

	records := ee.processResults(resp)

	return records, nil
}

func (ee *EntityExtractor) processResults(response string) []Record {
	parts := strings.Split(response, prompts.DefaultRecordDelimiter)
	records := make([]string, 0)
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			records = append(records, trimmed)
		}
	}

	parsedRecords := make([]Record, 0)
	for _, raw := range records {
		record := parseRecord(raw)
		parsedRecords = append(parsedRecords, record)
	}

	return parsedRecords
}

func parseRecord(record string) Record {
	record = strings.Trim(record, "\n")
	record = strings.TrimPrefix(record, "(")
	record = strings.TrimSuffix(record, ")")

	attrs := strings.Split(record, prompts.DefaultTupleDelimiter)

	for i, v := range attrs {
		str := cleanString(v)
		attrs[i] = str
	}

	recordType := attrs[0]
	switch recordType {
	case "entity":
		return &Entity{
			Name:         attrs[1],
			internalType: attrs[2],
			Description:  attrs[3],
		}
	case "relationship":
		return &Relationship{
			Entity1:  attrs[1],
			Entity2:  attrs[2],
			Relation: attrs[3],
		}
	default:
		return nil
	}
}

var re = regexp.MustCompile(`[\x00-\x1f\x7f-\x9f]`)

func cleanString(str string) string {
	cleaned := strings.TrimSpace(str)
	cleaned, _ = strconv.Unquote(cleaned)
	return re.ReplaceAllString(cleaned, "")
}
