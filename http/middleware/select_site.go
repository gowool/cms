package middleware

import (
	"context"
	"log/slog"
	"net/http"
	neturl "net/url"
	"slices"

	cmshttp "github.com/gowool/cms/http"
	"github.com/gowool/cms/internal"
)

type SelectSiteConfig struct {
	Skipper           func(r *http.Request) bool
	DecoratorStrategy cmshttp.DecoratorStrategy
	Selectors         []cmshttp.SiteSelector
	SelectorType      func(r *http.Request) cmshttp.SelectorType
	SiteNotFound      func(next http.Handler) http.Handler
	Logger            *slog.Logger
}

func SelectSite(config SelectSiteConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.Skipper != nil && config.Skipper(r) {
				next.ServeHTTP(w, r)
				return
			}

			if !config.DecoratorStrategy.IsRouteURIDecorable(r.Context(), r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			selectorType := config.SelectorType(r)
			selectorIndex := slices.IndexFunc(config.Selectors, func(selector cmshttp.SiteSelector) bool {
				return selector.Type() == selectorType
			})

			notFoundOrNext := func() {
				if config.SiteNotFound != nil {
					config.SiteNotFound(next).ServeHTTP(w, r)
					return
				}
				next.ServeHTTP(w, r)
			}

			if selectorIndex < 0 {
				config.Logger.Warn("site selector not found", "type", selectorType)

				notFoundOrNext()
				return
			}

			site, url, err := config.Selectors[selectorIndex].Select(r)
			if err != nil {
				config.Logger.Warn("site not found", "type", selectorType, "error", err, "host", r.URL.Host)

				notFoundOrNext()
				return
			}

			if url != "" {
				config.Logger.Info("redirecting", "type", selectorType, "host", r.URL.Host, "url", url)

				http.Redirect(w, r, url, http.StatusMovedPermanently)
				return
			}

			if site.ID == 0 {
				notFoundOrNext()
				return
			}

			config.Logger.Info("site found", "type", selectorType, "host", r.URL.Host, "site_id", site.ID)

			ctx := cmshttp.WithSiteSEO(r.Context(), site)
			ctx = context.WithValue(ctx, "url", internal.Must(neturl.Parse(r.URL.String())))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
