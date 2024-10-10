package fx

import (
	"net/http"
	"strings"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gowool/cms"
	cmsmiddleware "github.com/gowool/cms/middleware"
	"github.com/gowool/cms/repository"
)

type Middleware struct {
	Name       string
	Middleware echo.MiddlewareFunc
}

func NewMiddleware(name string, middleware echo.MiddlewareFunc) Middleware {
	return Middleware{
		Name:       name,
		Middleware: middleware,
	}
}

func AsMiddleware(middleware any) any {
	return fx.Annotate(
		middleware,
		fx.ResultTags(`group:"middleware"`),
	)
}

func RecoverMiddleware(cfg RecoverConfig, logger *zap.Logger) Middleware {
	return NewMiddleware("recover", middleware.RecoverWithConfig(middleware.RecoverConfig{
		Skipper:             cfg.Skipper,
		StackSize:           cfg.StackSize,
		DisableStackAll:     cfg.DisableStackAll,
		DisableErrorHandler: true,
		LogErrorFunc: func(_ echo.Context, err error, stack []byte) error {
			logger.Error("recover middleware", zap.Error(err), zap.String("stack", string(stack)))
			return err
		},
	}))
}

func BodyLimitMiddleware(cfg BodyLimitConfig) Middleware {
	return NewMiddleware("body_limit", middleware.BodyLimitWithConfig(middleware.BodyLimitConfig{
		Skipper: cfg.Skipper,
		Limit:   cfg.Limit,
	}))
}

func CompressMiddleware(cfg GzipConfig) Middleware {
	return NewMiddleware("compress", middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper:   cfg.Skipper,
		Level:     cfg.Level,
		MinLength: cfg.MinLength,
	}))
}

func DecompressMiddleware() Middleware {
	return NewMiddleware("decompress", middleware.Decompress())
}

func RequestIDMiddleware() Middleware {
	return NewMiddleware("request_id", middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: uuid.NewString,
	}))
}

func LoggerMiddleware(logger *zap.Logger) Middleware {
	return NewMiddleware("logger", middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		Skipper:          cmsmiddleware.SuffixPathSkipper(cms.NoLogExt...),
		HandleError:      true,
		LogLatency:       true,
		LogProtocol:      true,
		LogRemoteIP:      true,
		LogHost:          true,
		LogMethod:        true,
		LogURI:           true,
		LogURIPath:       true,
		LogRoutePath:     true,
		LogRequestID:     true,
		LogReferer:       true,
		LogUserAgent:     true,
		LogStatus:        true,
		LogError:         true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			attributes := []zap.Field{
				zap.Time("start-time", v.StartTime),
				zap.Duration("latency", v.Latency),
				zap.String("protocol", v.Protocol),
				zap.String("ip", v.RemoteIP),
				zap.String("host", v.Host),
				zap.String("method", v.Method),
				zap.String("uri", v.URI),
				zap.String("path", v.URIPath),
				zap.String("route", v.RoutePath),
				zap.String("request-id", v.RequestID),
				zap.String("referer", v.Referer),
				zap.String("user-agent", v.UserAgent),
				zap.Int("status", v.Status),
				zap.String("content-length", v.ContentLength),
				zap.Int64("response-size", v.ResponseSize),
			}

			if site := cms.CtxSite(c.Request().Context()); site != nil {
				attributes = append(attributes, zap.Dict("site",
					zap.Int64("id", site.ID),
					zap.String("name", site.Name),
					zap.String("host", site.Host),
					zap.String("locale", site.Locale),
					zap.String("relative_path", site.RelativePath),
				))
			}

			if page := cms.CtxPage(c.Request().Context()); page != nil {
				attributes = append(attributes, zap.Dict("page",
					zap.Int64("id", page.ID),
					zap.Int64("site_id", page.SiteID),
					zap.Int64p("parent_id", page.ParentID),
					zap.String("title", page.Title),
					zap.String("pattern", page.Pattern),
					zap.String("url", page.URL),
				))
			}

			if v.Error != nil {
				attributes = append(attributes, zap.Error(v.Error))
			}

			switch {
			case v.Status >= http.StatusBadRequest && v.Status < http.StatusInternalServerError:
				logger.Warn("incoming request", attributes...)
			case v.Status >= http.StatusInternalServerError:
				logger.Error("incoming request", attributes...)
			default:
				logger.Info("incoming request", attributes...)
			}
			return nil
		},
	}))
}

func SecureMiddleware(cfg SecureConfig) Middleware {
	return NewMiddleware("secure", middleware.SecureWithConfig(middleware.SecureConfig{
		Skipper:               cfg.Skipper,
		XSSProtection:         cfg.XSSProtection,
		ContentTypeNosniff:    cfg.ContentTypeNosniff,
		XFrameOptions:         cfg.XFrameOptions,
		HSTSMaxAge:            cfg.HSTSMaxAge,
		HSTSExcludeSubdomains: cfg.HSTSExcludeSubdomains,
		ContentSecurityPolicy: cfg.ContentSecurityPolicy,
		CSPReportOnly:         cfg.CSPReportOnly,
		HSTSPreloadEnabled:    cfg.HSTSPreloadEnabled,
		ReferrerPolicy:        cfg.ReferrerPolicy,
	}))
}

func CORSMiddleware(cfg CORSConfig) Middleware {
	return NewMiddleware("cors", middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:                                  cfg.Skipper,
		AllowOrigins:                             cfg.AllowOrigins,
		AllowOriginFunc:                          cfg.AllowOriginFunc,
		AllowMethods:                             cfg.AllowMethods,
		AllowHeaders:                             cfg.AllowHeaders,
		AllowCredentials:                         cfg.AllowCredentials,
		UnsafeWildcardOriginWithAllowCredentials: cfg.UnsafeWildcardOriginWithAllowCredentials,
		ExposeHeaders:                            cfg.ExposeHeaders,
		MaxAge:                                   cfg.MaxAge,
	}))
}

func CSRFMiddleware(cfg CSRFConfig) Middleware {
	return NewMiddleware("csrf", middleware.CSRFWithConfig(middleware.CSRFConfig{
		Skipper:        cfg.Skipper,
		ErrorHandler:   cfg.ErrorHandler,
		TokenLength:    cfg.TokenLength,
		TokenLookup:    cfg.TokenLookup,
		ContextKey:     cfg.ContextKey,
		CookieName:     cfg.Cookie.Name,
		CookieDomain:   cfg.Cookie.Domain,
		CookiePath:     cfg.Cookie.Path,
		CookieMaxAge:   int(cfg.Cookie.MaxAge.Seconds()),
		CookieSecure:   cfg.Cookie.Secure,
		CookieHTTPOnly: cfg.Cookie.HTTPOnly,
		CookieSameSite: cfg.Cookie.SameSite.HTTP(),
	}))
}

func BasicAuthMiddleware(repo repository.Admin) Middleware {
	return NewMiddleware("basic_auth", middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
		Skipper: func(c echo.Context) bool {
			h := c.Request().Header.Get(echo.HeaderAuthorization)
			return !strings.HasPrefix(strings.ToLower(h), "basic ")
		},
		Validator: cms.BasicAuthValidator(repo),
	}))
}

func JWTAuthMiddleware(repo repository.Admin, cfg JWTConfig) Middleware {
	return NewMiddleware("jwt_auth", middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper: func(c echo.Context) bool {
			h := c.Request().Header.Get(echo.HeaderAuthorization)
			return !strings.HasPrefix(strings.ToLower(h), "bearer ")
		},
		Validator: cms.JWTAuthValidator(repo, cfg.Secret),
	}))
}

func SessionMiddleware(sessionManager *scs.SessionManager) Middleware {
	return NewMiddleware("session", cmsmiddleware.Session(cmsmiddleware.SessionConfig{
		SessionManager: sessionManager,
	}))
}

func SiteSelectorMiddleware(siteSelector cms.SiteSelector, cfgRepository repository.Configuration) Middleware {
	return NewMiddleware("site_selector", cmsmiddleware.SiteSelector(cmsmiddleware.SiteSelectorConfig{
		SiteSelector:  siteSelector,
		CfgRepository: cfgRepository,
	}))
}

type PageSelectorParams struct {
	fx.In
	PageHandler    cms.PageHandler
	CfgRepository  repository.Configuration
	PageRepository repository.Page
}

func PageSelectorMiddleware(params PageSelectorParams) Middleware {
	return NewMiddleware("page_selector", cmsmiddleware.PageSelector(cmsmiddleware.PageSelectorConfig{
		PageHandler:    params.PageHandler,
		CfgRepository:  params.CfgRepository,
		PageRepository: params.PageRepository,
	}))
}

func HybridPageMiddleware(pageHandler cms.PageHandler, cfgRepository repository.Configuration) Middleware {
	return NewMiddleware("hybrid_page", cmsmiddleware.HybridPage(cmsmiddleware.HybridPageConfig{
		PageHandler:   pageHandler,
		CfgRepository: cfgRepository,
	}))
}
