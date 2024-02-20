package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/dlclark/regexp2"
)

const (
	headerRequestedWith    = "X-Requested-With"
	xmlHTTPRequest         = "XMLHttpRequest"
	headerPageDecorable    = "X-Page-Decorable"
	headerPageNotDecorable = "X-Page-Not-Decorable"
)

func SetDecoratorHeaders(ctx context.Context, w http.ResponseWriter) {
	if NotDecorableFrom(ctx) {
		w.Header().Set(headerPageNotDecorable, "1")
	} else if DecorableFrom(ctx) {
		w.Header().Set(headerPageDecorable, "1")
	}
}

type DecoratorStrategy interface {
	IsDecorable(r *http.Request, routeName string) bool
	IsRouteNameDecorable(ctx context.Context, routeName string) bool
	IsRouteURIDecorable(ctx context.Context, uri string) bool
}

type Ignore interface {
	GetIgnoreRoutes(ctx context.Context) []string
	GetIgnoreRoute(ctx context.Context) []*regexp2.Regexp
	GetIgnoreURI(ctx context.Context) []*regexp2.Regexp
}

type decoratorStrategy struct {
	ignore Ignore
}

func NewDecoratorStrategy(ignore Ignore) DecoratorStrategy {
	return decoratorStrategy{ignore: ignore}
}

func (ds decoratorStrategy) IsDecorable(r *http.Request, routeName string) bool {
	if NotDecorableFrom(r.Context()) {
		return false
	}

	if DecorableFrom(r.Context()) {
		return true
	}

	if r.Header.Get(headerRequestedWith) == xmlHTTPRequest {
		return false
	}

	return routeName != "" && ds.IsRouteNameDecorable(r.Context(), routeName) && ds.IsRouteURIDecorable(r.Context(), r.URL.Path)
}

func (ds decoratorStrategy) IsRouteNameDecorable(ctx context.Context, routeName string) bool {
	for _, route := range ds.ignore.GetIgnoreRoutes(ctx) {
		if strings.EqualFold(routeName, route) {
			return false
		}
	}

	for _, re := range ds.ignore.GetIgnoreRoute(ctx) {
		if ok, _ := re.MatchString(routeName); ok {
			return false
		}
	}
	return true
}

func (ds decoratorStrategy) IsRouteURIDecorable(ctx context.Context, uri string) bool {
	for _, re := range ds.ignore.GetIgnoreURI(ctx) {
		if ok, _ := re.MatchString(uri); ok {
			return false
		}
	}
	return true
}
