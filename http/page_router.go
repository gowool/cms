package http

import (
	"context"
	"errors"
	"slices"

	"github.com/gowool/cms"
)

var (
	ErrPageNotCMS        = errors.New("page is not cms")
	ErrPageRequestMethod = errors.New("page request method is not allowed")
)

type PageRouter interface {
	Match(ctx context.Context, site cms.Site, method string, path string) (cms.Page, error)
}

type pageRouter struct {
	service cms.PageService
}

func NewPageRouter(service cms.PageService) PageRouter {
	return pageRouter{service: service}
}

func (r pageRouter) Match(ctx context.Context, site cms.Site, method string, path string) (cms.Page, error) {
	page, err := r.service.GetByURL(ctx, site.ID, path)
	if err != nil {
		return cms.Page{}, err
	}

	if !page.IsCMS() {
		return cms.Page{}, ErrPageNotCMS
	}

	if len(page.RequestMethods) > 0 && !slices.Contains(page.RequestMethods, method) {
		return cms.Page{}, ErrPageRequestMethod
	}

	page.Site = &site

	return page, nil
}
