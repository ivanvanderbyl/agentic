package model_test

import (
	"testing"

	"github.com/ivanvanderbyl/graphrag-go/pkg/model"
	"github.com/stretchr/testify/require"
)

func TestReadingFromFile(t *testing.T) {
	a := require.New(t)

	tu, err := model.FromFile("../../testdata/2024-07-03/adjournment-ananda-rajah-michelle-mp-alp.md")
	a.NoError(err)

	a.Equal(969, tu.NTokens)
}
