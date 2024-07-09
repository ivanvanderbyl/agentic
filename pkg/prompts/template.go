package prompts

import (
	"bytes"
	"embed"
	"text/template"
)

//go:embed *.tmpl
var promptFS embed.FS

// Default delimiters
const DefaultTupleDelimiter = "<|>"
const DefaultRecordDelimiter = "##"
const DefaultCompletionDelimiter = "<|COMPLETE|>"

type PromptData struct {
	RecordDelimiter     string
	TupleDelimiter      string
	CompletionDelimiter string
}

type Data interface {
	isPromptData()
}

func (PromptData) isPromptData() {}

var DefaultPromptData = PromptData{
	RecordDelimiter:     DefaultRecordDelimiter,
	TupleDelimiter:      DefaultTupleDelimiter,
	CompletionDelimiter: DefaultCompletionDelimiter,
}

func RenderTemplate[T Data](templateName string, data T) (string, error) {
	tmpl, err := loadTemplates()
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func loadTemplates() (*template.Template, error) {
	tmpl, err := template.ParseFS(promptFS, "*")
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}
