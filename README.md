# `agentic`

This package serves as a workspace for exploring state of the art RAG approaches in Go. The initial focus is on the GraphRAG approach from Microsoft, but other approaches may be explored in the future.

## Areas of Exploration

- [x] GraphRAG
- [ ] dsRAG

## Findings

### GraphRAG

GraphRAG is seems to be mostly marketing spin from Microsoft. There are some legitimately cool ideas in there such as using LLMs to extract graph relationships from a given text.

The heavy use of multi-stage summarisation, however, produces results that are not accurate to the original text or useful for an executive or research audience because they too general and lacking concrete details.

### dsRAG

Source: https://github.com/D-Star-AI/dsRAG

There are three key methods used to improve performance over vanilla RAG systems:

Semantic sectioning
AutoContext
Relevant Segment Extraction (RSE)

We will implement these in Go to explore them further.
