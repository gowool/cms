package fx

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gowool/cms"
	"github.com/gowool/cms/api"
)

type HumaAPI interface {
	Register(*echo.Echo, huma.API)
}

func AsHumaAPI(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(HumaAPI)),
		fx.ResultTags(`group:"huma-api"`),
	)
}

func AsHumaAdminAPI(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(HumaAPI)),
		fx.ResultTags(`group:"huma-admin-api"`),
	)
}

type HumaMiddleware struct {
	Name       string
	Middleware func(huma.API) func(huma.Context, func(huma.Context))
}

func NewHumaMiddleware(name string, middleware func(huma.API) func(ctx huma.Context, next func(huma.Context))) HumaMiddleware {
	return HumaMiddleware{
		Name:       name,
		Middleware: middleware,
	}
}

func AsHumaMiddleware(middleware any) any {
	return fx.Annotate(
		middleware,
		fx.ResultTags(`group:"huma-middleware"`),
	)
}

func HumaAuthorizationMiddleware(authorizer cms.Authorizer, logger *zap.Logger) HumaMiddleware {
	return NewHumaMiddleware("authorization", api.Authorization(api.AuthorizationConfig{
		Authorizer: authorizer,
		Logger:     logger,
	}))
}
