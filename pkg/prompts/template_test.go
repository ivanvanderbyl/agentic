package prompts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestData struct {
	PromptData
	EntityTypes string
	InputText   string
}

func TestLoadingAndRendering(t *testing.T) {
	a := assert.New(t)
	result, err := RenderTemplate("entities", TestData{EntityTypes: "person,place,thing", PromptData: DefaultPromptData})
	a.NoError(err)
	a.Contains(result, "person,place,thing")
}
