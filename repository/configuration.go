package repository

import (
	"context"

	"github.com/gowool/cms/model"
)

type Configuration interface {
	Load(ctx context.Context) (model.Configuration, error)
	Save(ctx context.Context, m *model.Configuration) error
}
