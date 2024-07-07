package loaders

import "github.com/ivanvanderbyl/graphrag-go/pkg/model"

type (
	Loader interface {
		Load(path string) (*model.TextUnit, error)
	}
)
