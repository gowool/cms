package repository

import (
	"context"

	"github.com/gowool/cms/model"
)

type Admin interface {
	repository[model.Admin, int64]
	FindByEmail(ctx context.Context, email string) (model.Admin, error)
}
