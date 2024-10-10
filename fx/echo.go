package fx

import (
	"io/fs"
	"net"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/gowool/cms"
)

type EchoParams struct {
	fx.In

	ErrorHandler      *cms.ErrorHandler
	Renderer          echo.Renderer
	Validator         echo.Validator
	IPExtractor       echo.IPExtractor
	GlobalMiddlewares GlobalConfig
	Areas             AreasConfig
	Middlewares       []Middleware `group:"middleware"`
	WebHandlers       []Handler    `group:"web"`
	APIHandlers       []Handler    `group:"api"`
	AdminAPIHandlers  []Handler    `group:"admin-api"`
	StaticHandlers    []Handler    `group:"static"`
	Filesystem        fs.FS        `name:"fs-static"`
}

func NewEcho(params EchoParams) *echo.Echo {
	e := echo.New()
	e.Debug = false
	e.HideBanner = true
	e.HidePort = true
	e.Server.Handler = nil
	e.Server = nil
	e.TLSServer.Handler = nil
	e.TLSServer = nil
	e.StdLogger = nil
	e.Logger = nil
	e.Renderer = params.Renderer
	e.Validator = params.Validator
	e.IPExtractor = params.IPExtractor
	e.Filesystem = params.Filesystem
	e.HTTPErrorHandler = params.ErrorHandler.Handle

	middlewares := make(map[string]echo.MiddlewareFunc)
	for _, middleware := range params.Middlewares {
		middlewares[middleware.Name] = middleware.Middleware
	}

	for _, name := range params.GlobalMiddlewares.BeforeRouter {
		if middleware, ok := middlewares[name]; ok {
			e.Pre(middleware)
		}
	}

	for _, name := range params.GlobalMiddlewares.AfterRouter {
		if middleware, ok := middlewares[name]; ok {
			e.Use(middleware)
		}
	}

	if params.Areas.Static.Enabled {
		g := e.Group(params.Areas.Static.BasePath)
		for _, name := range params.Areas.Static.Middlewares {
			if middleware, ok := middlewares[name]; ok {
				g.Use(middleware)
			}
		}
		for _, h := range params.StaticHandlers {
			h.Register(e, g)
		}
	}

	if params.Areas.API.Enabled {
		g := e.Group(params.Areas.API.BasePath)
		for _, name := range params.Areas.API.Middlewares {
			if middleware, ok := middlewares[name]; ok {
				g.Use(middleware)
			}
		}
		for _, h := range params.APIHandlers {
			h.Register(e, g)
		}
	}

	if params.Areas.AdminAPI.Enabled {
		g := e.Group(params.Areas.AdminAPI.BasePath)
		for _, name := range params.Areas.AdminAPI.Middlewares {
			if middleware, ok := middlewares[name]; ok {
				g.Use(middleware)
			}
		}
		for _, h := range params.AdminAPIHandlers {
			h.Register(e, g)
		}
	}

	if params.Areas.Web.Enabled {
		g := e.Group(params.Areas.Web.BasePath)
		for _, name := range params.Areas.Web.Middlewares {
			if middleware, ok := middlewares[name]; ok {
				g.Use(middleware)
			}
		}
		for _, h := range params.WebHandlers {
			h.Register(e, g)
		}
	}

	return e
}

func IPExtractor() echo.IPExtractor {
	return func(r *http.Request) string {
		if ip := r.Header.Get(echo.HeaderXForwardedFor); ip != "" {
			i := strings.IndexAny(ip, ",")
			if i > 0 {
				xffip := strings.TrimSpace(ip[:i])
				xffip = strings.TrimPrefix(xffip, "[")
				xffip = strings.TrimSuffix(xffip, "]")
				return xffip
			}
			return ip
		}
		if ip := r.Header.Get(echo.HeaderXRealIP); ip != "" {
			ip = strings.TrimPrefix(ip, "[")
			ip = strings.TrimSuffix(ip, "]")
			return ip
		}
		ra, _, _ := net.SplitHostPort(r.RemoteAddr)
		return ra
	}
}
