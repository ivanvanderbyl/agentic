package model

import (
	"fmt"
	"io"
	"os"

	"github.com/pkoukk/tiktoken-go"
)

const encoding = "cl100k_base"

// Identified represents a base struct with identification fields.
type Identified struct {
	ID      string `json:"id"`
	ShortID string `json:"short_id,omitempty"`
}

// TextUnit represents a protocol for a TextUnit item in a Document database.
type TextUnit struct {
	Identified
	Text            string              `json:"text"`
	TextEmbedding   []float64           `json:"text_embedding,omitempty"`
	EntityIDs       []string            `json:"entity_ids,omitempty"`
	RelationshipIDs []string            `json:"relationship_ids,omitempty"`
	CovariateIDs    map[string][]string `json:"covariate_ids,omitempty"`
	NTokens         int                 `json:"n_tokens,omitempty"`
	DocumentIDs     []string            `json:"document_ids,omitempty"`
	Attributes      map[string]any      `json:"attributes,omitempty"`
}

func FromFile(file string) (*TextUnit, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var tu *TextUnit

	body, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	tke, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		err = fmt.Errorf("getEncoding: %v", err)
		return nil, err
	}

	token := tke.Encode(string(body), nil, nil)
	tu = &TextUnit{
		Text:    string(body),
		NTokens: len(token),
	}

	for _, t := range token {
		fmt.Println(t)
		// tu.TextEmbedding = append(tu.TextEmbedding, t.Embedding)
	}

	return tu, nil
}
