package entity

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/ivanvanderbyl/graphrag-go/pkg/llm"
	"github.com/ivanvanderbyl/graphrag-go/pkg/prompts"
)

type (
	EntityExtractor struct {
		llm              llm.LLM
		ExtractionPrompt string
		EntityTypes      []string
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
		NodeID() string
		Type() string
	}

	Entity struct {
		Name         string
		Description  string
		Embedding    []float32
		id           string
		internalType string
	}

	Relationship struct {
		Entity1  string
		Relation string
		Entity2  string
		Keyword  string
		Weight   int
		id       string
	}
)

func (Entity) isNode() {}
func (e *Entity) Type() string {
	return e.internalType
}
func (e *Entity) String() string {
	return fmt.Sprintf("Entity{Type: %q, Name: %q, Description: %q}", e.internalType, e.Name, e.Description)
}
func (e *Entity) NodeID() string {
	return e.id
}

func (Relationship) isNode() {}
func (r *Relationship) Type() string {
	return "relationship"
}
func (e *Relationship) NodeID() string {
	return e.id
}
func (r *Relationship) String() string {
	buf := new(strings.Builder)
	buf.WriteString("Relationship{")
	buf.WriteString(fmt.Sprintf("From: %q", r.Entity1))
	buf.WriteString(", ")
	buf.WriteString(fmt.Sprintf("Relation: %q", r.Relation))
	buf.WriteString(", ")
	buf.WriteString(fmt.Sprintf("To: %q", r.Entity2))
	buf.WriteString(", ")
	buf.WriteString(fmt.Sprintf("Keyword: %q", r.Keyword))
	buf.WriteString(", ")
	buf.WriteString(fmt.Sprintf("Weight: %d", r.Weight))
	buf.WriteString("}")
	return buf.String()
}

var DefaultEntityTypes = []string{"organization", "person", "policy", "bill", "geo", "event", "role", "electorate"}

func NewEntityExtractor(llm llm.LLM, opts ...Option) *EntityExtractor {
	e := &EntityExtractor{
		llm:         llm,
		EntityTypes: DefaultEntityTypes,
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func WithEntityTypes(types []string) Option {
	return func(e *EntityExtractor) {
		e.EntityTypes = types
	}
}

func (ee *EntityExtractor) Extract(ctx context.Context, text string) ([]Record, error) {
	prompt, err := prompts.RenderTemplate(prompts.EntitiesTemplate, Data{
		EntityTypes: ee.EntityTypes,
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

	if err := ee.createEmbeddings(ctx, records); err != nil {
		return nil, err
	}

	return records, nil
}

func (ee *EntityExtractor) createEmbeddings(ctx context.Context, records []Record) error {
	for _, record := range records {
		switch r := record.(type) {
		case *Entity:
			embedding, err := ee.llm.Embedding(ctx, r.Description)
			if err != nil {
				return err
			}
			r.Embedding = embedding
		}
	}
	return nil
}

func (ee *EntityExtractor) processResults(response string) []Record {
	before, ok := strings.CutSuffix(response, prompts.DefaultCompletionDelimiter)
	if ok {
		response = before
	}

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
		id := uuid.NewSHA1(uuid.NameSpaceOID, []byte(fmt.Sprintf("%s%s", attrs[1], attrs[3])))

		return &Entity{
			Name:         attrs[1],
			internalType: attrs[2],
			Description:  attrs[3],
			id:           id.String(),
		}
	case "relationship":
		weight, err := strconv.Atoi(attrs[4])
		if err != nil {
			weight = 0
		}

		id := uuid.NewSHA1(uuid.NameSpaceOID, []byte(fmt.Sprintf("%s%s%s", attrs[1], attrs[2], attrs[3])))

		return &Relationship{
			Entity1:  attrs[1],
			Entity2:  attrs[2],
			Relation: attrs[3],
			Keyword:  attrs[5],
			Weight:   weight,
			id:       id.String(),
		}
	default:
		return nil
	}
}

var re = regexp.MustCompile(`[\x00-\x1f\x7f-\x9f]`)

func cleanString(str string) string {
	str = strings.TrimSpace(str)
	cleaned, err := strconv.Unquote(str)
	if err != nil {
		cleaned = str
	}
	return re.ReplaceAllString(cleaned, "")
}
