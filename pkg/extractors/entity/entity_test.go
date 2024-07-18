package entity

import (
	"testing"
)

var testDocument = "##(\"entity\"<|>\"Byrnes, Alison MP\"<|>\"person\"<|>\"Byrnes, Alison MP is a Member of Parliament representing the electorate of Cunningham from the Australian Labor Party (ALP).\")##\n##(\"entity\"<|>\"Cunningham\"<|>\"place\"<|>\"Cunningham is an electorate represented by Alison Byrnes MP.\")##\n##(\"entity\"<|>\"Minister for Industry and Science\"<|>\"person\"<|>\"The Minister for Industry and Science is a government official responsible for overseeing industrial and scientific development.\")##\n##(\"entity\"<|>\"Albanese Labor Government\"<|>\"person\"<|>\"The Albanese Labor Government, led by Anthony Albanese, is the current ruling government in Australia.\")##\n##(\"entity\"<|>\"Future Made in Australia plan\"<|>\"thing\"<|>\"The Future Made in Australia plan is a government initiative aimed at securing jobs and investing in local manufacturing.\")##\n##(\"relationship\"<|>\"Byrnes, Alison MP\"<|>\"Cunningham\"<|>\"Byrnes, Alison MP represents the Cunningham electorate in Parliament.\"<|>10)##\n##(\"relationship\"<|>\"Byrnes, Alison MP\"<|>\"Minister for Industry and Science\"<|>\"Byrnes, Alison MP addresses a question to the Minister for Industry and Science during a parliamentary session.\"<|>8)##\n##(\"relationship\"<|>\"Byrnes, Alison MP\"<|>\"Albanese Labor Government\"<|>\"Byrnes, Alison MP, as a member of the ALP, supports and asks about the policies of the Albanese Labor Government.\"<|>9)##\n##(\"relationship\"<|>\"Albanese Labor Government\"<|>\"Future Made in Australia plan\"<|>\"The Albanese Labor Government is implementing the Future Made in Australia plan to secure jobs and boost local manufacturing.\"<|>10)<|COMPLETE|>"

func TestEntityParsing(t *testing.T) {
	// a := assert.New(t)

	extractor := NewEntityExtractor(nil)
	extractor.processResults(testDocument)
}
