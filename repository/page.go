package repository

import (
	"context"
	"errors"
	"time"

	"github.com/gowool/cms/model"
)

var ErrPageNotFound = errors.New("page not found")

type Page interface {
	repository[model.Page, int64]
	FindByParentID(ctx context.Context, parentID int64, now time.Time) ([]model.Page, error)
	FindByPattern(ctx context.Context, siteID int64, pattern string, now time.Time) (model.Page, error)
	FindByAlias(ctx context.Context, siteID int64, alias string, now time.Time) (model.Page, error)
	FindByURL(ctx context.Context, siteID int64, url string, now time.Time) (model.Page, error)
}
