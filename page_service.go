package cms

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type PageService interface {
	CreateWithDefaults(ctx context.Context, defaults map[string]interface{}) (Page, error)
	Save(ctx context.Context, page *Page) error
	GetByID(ctx context.Context, id int64) (Page, error)
	GetByURL(ctx context.Context, siteID int64, url string) (Page, error)
	GetByRouteName(ctx context.Context, siteID int64, routeName string) (Page, error)
	GetByPageAlias(ctx context.Context, siteID int64, pageAlias string) (Page, error)
	GetChildren(ctx context.Context, parentID int64) ([]Page, error)
}

type pageService struct {
	storage  PageStorage
	defaults PageDefaults
}

func NewPageService(storage PageStorage, defaults PageDefaults) PageService {
	return pageService{
		storage:  storage,
		defaults: defaults,
	}
}

func (s pageService) CreateWithDefaults(ctx context.Context, defaults map[string]interface{}) (Page, error) {
	var page Page

	if data, err := s.defaults.GetDefaults(ctx); data != nil && err == nil {
		if err = defaultsUnmarshal(data, &page); err != nil {
			return page, err
		}
	}

	if routeName, ok := defaults["route_name"]; ok {
		if data, err := s.defaults.GetRouteDefaults(ctx, fmt.Sprintf("%v", routeName)); data != nil && err == nil {
			if err = defaultsUnmarshal(data, &page); err != nil {
				return page, err
			}
		}
	}

	err := defaultsUnmarshal(defaults, &page)

	return page, err
}

func (s pageService) Save(ctx context.Context, page *Page) error {
	if page.ID != 0 {
		return s.storage.Update(ctx, page)
	}
	return s.storage.Create(ctx, page)
}

func (s pageService) GetByID(ctx context.Context, id int64) (Page, error) {
	return s.storage.FindByID(ctx, id)
}

func (s pageService) GetByURL(ctx context.Context, siteID int64, url string) (Page, error) {
	return s.storage.FindByURL(ctx, siteID, url, time.Now())
}

func (s pageService) GetByRouteName(ctx context.Context, siteID int64, routeName string) (Page, error) {
	return s.storage.FindByRouteName(ctx, siteID, routeName, time.Now())
}

func (s pageService) GetByPageAlias(ctx context.Context, siteID int64, pageAlias string) (Page, error) {
	return s.storage.FindByPageAlias(ctx, siteID, pageAlias, time.Now())
}

func (s pageService) GetChildren(ctx context.Context, parentID int64) ([]Page, error) {
	return s.storage.FindByParentID(ctx, parentID, time.Now())
}

func defaultsUnmarshal(defaults map[string]interface{}, page *Page) error {
	raw, err := json.Marshal(defaults)
	if err != nil {
		return err
	}

	return json.Unmarshal(raw, page)
}
