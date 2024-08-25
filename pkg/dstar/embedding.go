package dstar

import (
	"context"

	"github.com/ivanvanderbyl/graphrag-go/pkg/llm"
)

type Dimensionality int

const (
	DimensionalityEmbedEnglishV3_0           Dimensionality = 1024
	DimensionalityEmbedMultilingualV3_0      Dimensionality = 1024
	DimensionalityEmbedEnglishLightV3_0      Dimensionality = 384
	DimensionalityEmbedMultilingualLightV3_0 Dimensionality = 384
	DimensionalityVoyageLarge_2              Dimensionality = 1536
	DimensionalityVoyageLaw_2                Dimensionality = 1024
	DimensionalityVoyageCode_2               Dimensionality = 1536
	DimensionalityLlama2                     Dimensionality = 4096
	DimensionalityLlama3                     Dimensionality = 4096
	DimensionalityAllMinilm                  Dimensionality = 384
	DimensionalityNomicEmbedText             Dimensionality = 768
)

func CreateEmbedding(ctx context.Context, lm llm.LLM, dimensions Dimensionality, input string) ([]float32, error) {
	vectors, err := lm.Embedding(ctx, input,
		llm.WithDimensions(int(dimensions)),
		llm.WithModel(llm.ModelTextEmbedding3Small),
	)
	if err != nil {
		return nil, err
	}

	return vectors, nil
}
