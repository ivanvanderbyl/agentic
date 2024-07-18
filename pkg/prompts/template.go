package prompts

import (
	"bytes"
	"embed"
	"strings"
	"text/template"
)

//go:embed *.tmpl
var promptFS embed.FS

const (
	EntitiesTemplate = "entities"
	ClaimsTemplate   = "claims"
)

// Default delimiters
const DefaultTupleDelimiter = "<|>"
const DefaultRecordDelimiter = "##"
const DefaultCompletionDelimiter = "<|COMPLETE|>"

var DefaultEntityTypes = [...]string{"organization", "person", "geo", "event", "role", "electorate"}

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

func joinStrings(strs []string) string {
	return strings.Join(strs, ", ")
}

func RenderTemplate[T Data](templateName string, data T) (string, error) {
	tmpl, err := loadTemplates()
	if err != nil {
		return "", err
	}

	funcMap := template.FuncMap{
		"joinStrings": joinStrings,
	}

	var buf bytes.Buffer
	err = tmpl.Funcs(funcMap).ExecuteTemplate(&buf, templateName, data)
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
