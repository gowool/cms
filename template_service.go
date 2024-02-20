package cms

import (
	"context"
	"time"
)

type TemplateService interface {
	GetByID(ctx context.Context, id int64) (Template, error)
	GetByCode(ctx context.Context, code string) (Template, error)
	IsFresh(ctx context.Context, code string, t int64) (bool, error)
	Exists(ctx context.Context, code string) (bool, error)
}

type templateService struct {
	storage TemplateStorage
}

func NewTemplateService(storage TemplateStorage) TemplateService {
	return templateService{storage: storage}
}

func (s templateService) GetByID(ctx context.Context, id int64) (Template, error) {
	return s.storage.FindByID(ctx, id)
}

func (s templateService) GetByCode(ctx context.Context, name string) (Template, error) {
	return s.storage.FindByCode(ctx, name, time.Now())
}

func (s templateService) IsFresh(ctx context.Context, code string, t int64) (bool, error) {
	item, err := s.GetByCode(ctx, code)
	if err != nil {
		return false, err
	}

	return item.Updated.Unix() < t, nil
}

func (s templateService) Exists(ctx context.Context, code string) (bool, error) {
	_, err := s.GetByCode(ctx, code)
	return err == nil, err
}
