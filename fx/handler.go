package fx

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type Handler interface {
	Register(*echo.Echo, *echo.Group)
}

func AsHandler(f any, group string) any {
	return fx.Annotate(
		f,
		fx.As(new(Handler)),
		fx.ResultTags(fmt.Sprintf(`group:"%s"`, group)),
	)
}

func AsStatic(f any) any {
	return AsHandler(f, "static")
}

func AsAPI(f any) any {
	return AsHandler(f, "api")
}

func AsAdminAPI(f any) any {
	return AsHandler(f, "admin-api")
}

func AsWeb(f any) any {
	return AsHandler(f, "web")
}
