package dstar_test

import (
	"context"
	"os"
	"testing"

	"github.com/ivanvanderbyl/graphrag-go/pkg/dstar"
	"github.com/ivanvanderbyl/graphrag-go/pkg/llm"
	"github.com/stretchr/testify/assert"
)

func TestSummarizingChunk(t *testing.T) {
	a := assert.New(t)
	lm := llm.NewOpenAI(llm.WithCache("testdata"))
	chunk, err := os.ReadFile("../../testdata/2024-07-03/adjournment-khalil-peter-mp-alp.md")
	a.NoError(err)
	summary, err := dstar.GenerateQueriesFromChunk(context.Background(), lm, string(chunk))
	a.NoError(err)
	a.NotEmpty(summary)
	t.Log(summary)
}

func TestTruncateToTokenLength(t *testing.T) {
	a := assert.New(t)
	text := `On 20 June, health minister Mark Butler announced $1.1 million to the NHMRC`
	truncated, err := dstar.TruncateToTokenLength(text, 10)
	a.NoError(err)
	a.Equal("On 20 June, health minister Mark Butler announced", truncated)
}
