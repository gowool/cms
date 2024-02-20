package cms

import (
	"context"
	"time"
)

type SiteService interface {
	CreateLocalhost() Site
	Save(ctx context.Context, site *Site) error
	GetByID(ctx context.Context, id int64) (Site, error)
	GetByHost(ctx context.Context, host string) ([]Site, error)
	GetAll(ctx context.Context) ([]Site, error)
}

type siteService struct {
	storage SiteStorage
}

func NewSiteService(storage SiteStorage) SiteService {
	return siteService{storage: storage}
}

func (s siteService) CreateLocalhost() Site {
	return Site{
		Name:      "Localhost",
		Separator: " | ",
		Host:      "localhost",
		Locale:    "en",
		IsDefault: true,
	}
}

func (s siteService) Save(ctx context.Context, site *Site) error {
	if site.ID != 0 {
		return s.storage.Update(ctx, site)
	}
	return s.storage.Create(ctx, site)
}

func (s siteService) GetByID(ctx context.Context, id int64) (Site, error) {
	return s.storage.FindByID(ctx, id)
}

func (s siteService) GetByHost(ctx context.Context, host string) ([]Site, error) {
	return s.storage.FindByHosts(ctx, []string{host, "localhost"}, time.Now())
}

func (s siteService) GetAll(ctx context.Context) ([]Site, error) {
	return s.storage.FindAll(ctx)
}
