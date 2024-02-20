package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gowool/cms"
	cmshttp "github.com/gowool/cms/http"
)

type SelectHybridPageConfig struct {
	Skipper           func(r *http.Request) bool
	DecoratorStrategy cmshttp.DecoratorStrategy
	PageService       cms.PageService
	RouteName         func(r *http.Request) string
	PageNotFound      func(next http.Handler) http.Handler
	Logger            *slog.Logger
}

func SelectHybridPage(config SelectHybridPageConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.Skipper != nil && config.Skipper(r) {
				next.ServeHTTP(w, r)
				return
			}

			routeName := config.RouteName(r)
			if routeName == cms.PageRouteCMS {
				next.ServeHTTP(w, r)
				return
			}

			if !config.DecoratorStrategy.IsDecorable(r, routeName) {
				next.ServeHTTP(w, r)
				return
			}

			site, _ := cmshttp.SiteFrom(r.Context())
			if site.ID == 0 {
				config.Logger.Error("site not found in request context", "path", r.URL.Path)
				http.Error(w, "site not found in request context", http.StatusInternalServerError)
				return
			}

			page, err := config.PageService.GetByRouteName(r.Context(), site.ID, routeName)
			if err != nil {
				config.Logger.Warn("page not found by route name", "site_id", site.ID, "route_name", routeName)

				if config.PageNotFound != nil {
					config.PageNotFound(next).ServeHTTP(w, r)
					return
				}

				next.ServeHTTP(w, r)
				return
			}

			page.Site = &site

			r = r.WithContext(cmshttp.WithPageSEO(r.Context(), page))

			next.ServeHTTP(w, r)
		})
	}
}
