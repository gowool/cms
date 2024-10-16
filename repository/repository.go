package repository

import (
	"context"
	"errors"

	"github.com/gowool/cr"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrUniqueViolation = errors.New("unique violation")
)

type Repository[M any, ID any] interface {
	FindAndCount(ctx context.Context, criteria *cr.Criteria) ([]M, int, error)
	FindByID(ctx context.Context, id ID) (M, error)
	Delete(ctx context.Context, ids ...ID) error
	Create(ctx context.Context, m *M) error
	Update(ctx context.Context, m *M) error
}

type repository[M any, ID any] interface {
	Repository[M, ID]
	Find(ctx context.Context, criteria *cr.Criteria) ([]M, error)
}
