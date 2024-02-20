package middleware

import (
	"net/http"

	cmshttp "github.com/gowool/cms/http"
)

func PageSEO(skipper func(r *http.Request) bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if skipper != nil && skipper(r) {
				next.ServeHTTP(w, r)
				return
			}

			ctx := cmshttp.WithSEO(r.Context(), cmshttp.NewPageSEO())

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
