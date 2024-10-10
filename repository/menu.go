package repository

import (
	"context"

	"github.com/gowool/cms/model"
)

type Menu interface {
	repository[model.Menu, int64]
	FindByHandle(ctx context.Context, handle string) (model.Menu, error)
}
