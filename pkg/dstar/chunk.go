package dstar

import (
	"context"
	"fmt"

	"github.com/ivanvanderbyl/graphrag-go/pkg/llm"
	"github.com/tiktoken-go/tokenizer"
)

func SummarizeChunk(ctx context.Context, lm llm.LLM, chunk string, opts ...llm.Option) (string, error) {
	prompt := createSummaryPrompt(chunk)
	return lm.Generate(ctx, prompt, append([]llm.Option{llm.WithModel("gpt-4o-mini")}, opts...)...)
}

func createSummaryPrompt(chunk string) string {
	return `INSTRUCTIONS
What is the following section about?

Your response should be a single sentence, and it shouldn't be an excessively long sentence. DO NOT respond with anything else.

Your response should take the form of "This section is about: X". For example, if the section is a balance sheet from a financial report about Apple, your response might be "This section is about: the financial position of Apple as of the end of the fiscal year." If the section is a chapter from a book on the history of the United States, and this chapter covers the Civil War, your response might be "This section is about: the causes and consequences of the American Civil War."

SECTION
` + chunk
}

func GenerateQueriesFromChunk(ctx context.Context, lm llm.LLM, chunk string) (string, error) {
	decoded, err := TruncateToTokenLength(chunk, 8000)
	if err != nil {
		return "", err
	}

	title, err := lm.Generate(ctx, createDocumentTitlePrompt(decoded), llm.WithModel("gpt-4o-mini"))
	if err != nil {
		return "", err
	}

	docSummary, err := lm.Generate(ctx, createDocumentSummarizationPrompt(title, decoded), llm.WithModel("gpt-4o-mini"))
	if err != nil {
		return "", err
	}

	return docSummary, nil
}

func TruncateToTokenLength(text string, maxTokens int) (string, error) {
	enc, err := tokenizer.Get(tokenizer.Cl100kBase)
	if err != nil {
		return "", err
	}

	ids, _, err := enc.Encode(text)
	if err != nil {
		return "", err
	}

	if len(ids) > maxTokens {
		ids = ids[:maxTokens]
	}
	decoded, err := enc.Decode(ids)
	if err != nil {
		return "", err
	}

	return decoded, nil
}

func createQueryPrompt(chunk string) string {
	return systemAutoQuertyPrompt + chunk
}

var systemAutoQuertyPrompt = `You are a query generation system.
Please generate one or more search queries (up to a maximum of 5) based on the provided user input. DO NOT generate the answer, just queries.

Each of the queries you generate will be used to search a knowledge base for information that can be used to respond to the user input. Make sure each query is specific enough to return relevant information. If multiple pieces of information would be useful, you should generate multiple queries, one for each specific piece of information needed.

User input:
`

var documentTitlePrompt = `
INSTRUCTIONS
What is the title of the following document?

Your response MUST be the title of the document, and nothing else. DO NOT respond with anything else.

DOCUMENT
%s
`

func createDocumentTitlePrompt(document string) string {
	return fmt.Sprintf(documentTitlePrompt, document)
}

var documentSummarizationPrompt = `
INSTRUCTIONS
What is the following document, and what is it about?

Your response should be a single sentence, and it shouldn't be an excessively long sentence. DO NOT respond with anything else.

Your response should take the form of "This document is about: X". For example, if the document is a book about the history of the United States called A People's History of the United States, your response might be "This document is about: the history of the United States, covering the period from 1776 to the present day." If the document is the 2023 Form 10-K for Apple Inc., your response might be "This document is about: the financial performance and operations of Apple Inc. during the fiscal year 2023."

%s

DOCUMENT
Document name: %s

%s
`

var truncationMessage = `Also note that the document text provided below is just the first ~500 words of the document. That should be plenty for this task. Your response should still pertain to the entire document, not just the text provided below.`

func createDocumentSummarizationPrompt(title, document string) string {
	return fmt.Sprintf(documentSummarizationPrompt, truncationMessage, title, document)
}
