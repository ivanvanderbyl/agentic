package claims

import (
	"github.com/ivanvanderbyl/graphrag-go/pkg/prompts"
)

type PromptData struct {
	prompts.PromptData
	InputText        string
	ClaimDescription string
	EntitySpecs      string
}

func GetCompletionPrompt(data PromptData) (string, error) {
	return prompts.RenderTemplate("claims", data)
}
