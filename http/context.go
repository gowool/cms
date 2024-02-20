package http

import (
	"context"
	"errors"

	"github.com/gowool/cms"
)

type (
	pageSEOKey      struct{}
	siteKey         struct{}
	pageKey         struct{}
	decorableKey    struct{}
	notDecorableKey struct{}
)

func WithSiteSEO(ctx context.Context, site cms.Site) context.Context {
	seo := SEOFrom(ctx)
	seo.AddFromSite(site)

	ctx = WithSite(ctx, site)
	ctx = WithSEO(ctx, seo)

	return ctx
}

func WithPageSEO(ctx context.Context, page cms.Page) context.Context {
	seo := SEOFrom(ctx)
	seo.AddFromPage(page)

	ctx = WithPage(ctx, page)
	ctx = WithSEO(ctx, seo)

	return ctx
}

func WithSEO(ctx context.Context, seo PageSEO) context.Context {
	return context.WithValue(ctx, pageSEOKey{}, seo)
}

func SEOFrom(ctx context.Context) PageSEO {
	if seo, ok := ctx.Value(pageSEOKey{}).(PageSEO); ok {
		return seo
	}
	return NewPageSEO()
}

func WithSite(ctx context.Context, site cms.Site) context.Context {
	if site.ID == 0 {
		panic(errors.New("site is not saved"))
	}
	return context.WithValue(ctx, siteKey{}, site)
}

func SiteFrom(ctx context.Context) (cms.Site, bool) {
	site, ok := ctx.Value(siteKey{}).(cms.Site)
	return site, ok
}

func WithPage(ctx context.Context, page cms.Page) context.Context {
	if page.ID == 0 {
		panic(errors.New("page is not saved"))
	}
	return context.WithValue(ctx, pageKey{}, page)
}

func PageFrom(ctx context.Context) (cms.Page, bool) {
	page, ok := ctx.Value(pageKey{}).(cms.Page)
	return page, ok
}

func WithDecorable(ctx context.Context, decorable bool) context.Context {
	return context.WithValue(ctx, decorableKey{}, decorable)
}

func DecorableFrom(ctx context.Context) bool {
	decorable, _ := ctx.Value(decorableKey{}).(bool)
	return decorable
}

func WithNotDecorable(ctx context.Context, notDecorable bool) context.Context {
	return context.WithValue(ctx, notDecorableKey{}, notDecorable)
}

func NotDecorableFrom(ctx context.Context) bool {
	notDecorable, _ := ctx.Value(notDecorableKey{}).(bool)
	return notDecorable
}
