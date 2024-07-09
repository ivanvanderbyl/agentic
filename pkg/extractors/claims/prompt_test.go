package claims

import (
	"context"
	"os"
	"testing"

	"github.com/ivanvanderbyl/graphrag-go/pkg/llm"
	"github.com/stretchr/testify/assert"
)

func TestCompletionTemplate(t *testing.T) {
	a := assert.New(t)

	txt, err := GetCompletionPrompt(PromptData{
		InputText:        "Byron Bay is a renowned beachside town located in the far-northeastern corner of New South Wales, Australia, within Bundjalung Country. It is situated approximately 772 kilometers north of Sydney and 165 kilometers south of Brisbane. Cape Byron, adjacent to the town, is the easternmost point of mainland Australia.",
		ClaimDescription: "It is the most easterly point of mainland Australia.",
		EntitySpecs:      "PLACE",
	})
	a.NoError(err)
	a.NotEmpty(txt)

	oai := llm.NewOpenAI(llm.WithAPIKey(os.Getenv("OPENAI_API_KEY")))
	result, err := oai.Generate(context.TODO(), txt, llm.WithMaxTokens(4000))
	a.NoError(err)
	a.NotEmpty(result)
	t.Log(result)
}
