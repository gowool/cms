package middleware

import (
	"errors"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/gowool/cms"
	"github.com/gowool/cms/model"
	"github.com/gowool/cms/repository"
)

type PageSelectorConfig struct {
	Skipper        middleware.Skipper
	PageHandler    cms.PageHandler
	CfgRepository  repository.Configuration
	PageRepository repository.Page
}

func PageSelector(cfg PageSelectorConfig) echo.MiddlewareFunc {
	if cfg.PageHandler == nil {
		panic("page handler is not specified")
	}
	if cfg.CfgRepository == nil {
		panic("configuration repository is not specified")
	}
	if cfg.PageRepository == nil {
		panic("page repository is not specified")
	}
	if cfg.Skipper == nil {
		cfg.Skipper = middleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := c.Request()

			if cfg.Skipper(c) || cms.SkipSelectSite(r.Context()) || cms.SkipSelectPage(r.Context()) {
				return next(c)
			}

			configuration, err := cfg.CfgRepository.Load(r.Context())
			if err != nil {
				return err
			}

			if configuration.IgnoreURI(r.URL.Path) {
				return next(c)
			}

			site := cms.CtxSite(r.Context())
			if site == nil {
				return errors.Join(repository.ErrSiteNotFound, cms.ErrInternal)
			}

			var now time.Time
			if !cms.CtxEditor(r.Context()) {
				now = time.Now()
			}

			page, err := cfg.PageRepository.FindByURL(r.Context(), site.ID, r.URL.Path, now)
			if err != nil {
				if errors.Is(err, repository.ErrPageNotFound) {
					goto PATTERN
				}
				return err
			}

			if page.IsCMS() {
				return withPage(c, cfg.PageHandler.Handle, page)
			}

		PATTERN:
			if configuration.IgnorePattern(r.Pattern) {
				return next(c)
			}

			page, err = cfg.PageRepository.FindByPattern(r.Context(), site.ID, r.Pattern, now)
			if err != nil {
				return err
			}

			return withPage(c, next, page)
		}
	}
}

func withPage(c echo.Context, next echo.HandlerFunc, page model.Page) error {
	r := c.Request()
	ctx := cms.WithPage(r.Context(), &page)
	c.SetRequest(r.WithContext(ctx))
	return next(c)
}
