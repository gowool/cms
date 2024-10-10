package repository

import (
	"context"

	"github.com/gowool/cms/model"
)

type Template interface {
	repository[model.Template, int64]
	FindByName(ctx context.Context, name string) (model.Template, error)
}
