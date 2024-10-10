package fx

import (
	"net/http"

	"github.com/gowool/theme"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/gowool/cms"
)

var (
	OptionConfigurationRepository = fx.Provide(
		fx.Annotate(
			NewConfigurationRepository,
			fx.ParamTags("", `name:"repository-cache"`),
		),
	)
	OptionSiteRepository = fx.Provide(
		fx.Annotate(
			NewSiteRepository,
			fx.ParamTags("", `name:"repository-cache"`),
		),
	)
	OptionPageRepository = fx.Provide(
		fx.Annotate(
			NewPageRepository,
			fx.ParamTags("", `name:"repository-cache"`),
		),
	)
	OptionMenuRepository = fx.Provide(
		fx.Annotate(
			NewMenuRepository,
			fx.ParamTags("", `name:"repository-cache"`),
		),
	)
	OptionNodeRepository = fx.Provide(
		fx.Annotate(
			NewNodeRepository,
			fx.ParamTags("", `name:"repository-cache"`),
		),
	)
	OptionAdminRepository    = fx.Provide(NewAdminRepository)
	OptionTemplateRepository = fx.Provide(NewTemplateRepository)
	OptionThemeRepository    = fx.Provide(NewThemeRepository)

	OptionAuthorizer     = fx.Provide(fx.Annotate(cms.NewDefaultAuthorizer, fx.As(new(cms.Authorizer))))
	OptionSessionStore   = fx.Provide(NewSessionStore)
	OptionSessionManager = fx.Provide(NewSessionManager)
	OptionSeeder         = fx.Provide(NewSeeder)
	OptionMenu           = fx.Provide(fx.Annotate(cms.NewDefaultMenu, fx.As(new(cms.Menu))))
	OptionMatcher        = fx.Provide(
		fx.Annotate(
			cms.NewDefaultMatcher,
			fx.As(new(cms.Matcher)),
			fx.ParamTags(`group:"menu-voter"`),
		),
	)
	OptionURLVoter = fx.Provide(
		fx.Annotate(
			cms.NewURLVoter,
			fx.As(new(cms.Voter)),
			fx.ResultTags(`group:"menu-voter"`),
		),
	)

	OptionSiteSelector = fx.Provide(
		fx.Annotate(
			cms.NewDefaultSiteSelector,
			fx.As(new(cms.SiteSelector)),
		),
	)
	OptionPageHandler = fx.Provide(
		fx.Annotate(
			cms.NewDefaultPageHandler,
			fx.As(new(cms.PageHandler)),
		),
	)
	OptionPageCreateHandler = fx.Provide(cms.NewPageCreateHandler)
	OptionErrorHandler      = fx.Provide(cms.NewErrorHandler)
	OptionErrorResolver     = fx.Provide(cms.ErrorResolver)
	OptionRenderer          = fx.Provide(fx.Annotate(cms.NewRenderer, fx.As(new(echo.Renderer))))
	OptionIPExtractor       = fx.Provide(IPExtractor)
	OptionEcho              = fx.Provide(NewEcho)
	OptionHandler           = fx.Provide(func(e *echo.Echo) http.Handler { return e })

	OptionThemeFuncMap = fx.Provide(fx.Annotate(FuncMap, fx.ResultTags(`group:"theme-func-map"`)))
	OptionThemeLoader  = fx.Provide(fx.Annotate(theme.NewRepositoryLoader, fx.As(new(theme.Loader))))

	OptionRecoverMiddleware      = fx.Provide(AsMiddleware(RecoverMiddleware))
	OptionBodyLimitMiddleware    = fx.Provide(AsMiddleware(BodyLimitMiddleware))
	OptionCompressMiddleware     = fx.Provide(AsMiddleware(CompressMiddleware))
	OptionDecompressMiddleware   = fx.Provide(AsMiddleware(DecompressMiddleware))
	OptionRequestIDMiddleware    = fx.Provide(AsMiddleware(RequestIDMiddleware))
	OptionLoggerMiddleware       = fx.Provide(AsMiddleware(LoggerMiddleware))
	OptionSecureMiddleware       = fx.Provide(AsMiddleware(SecureMiddleware))
	OptionCORSMiddleware         = fx.Provide(AsMiddleware(CORSMiddleware))
	OptionCSRFMiddleware         = fx.Provide(AsMiddleware(CSRFMiddleware))
	OptionBasicAuthMiddleware    = fx.Provide(AsMiddleware(BasicAuthMiddleware))
	OptionJWTAuthMiddleware      = fx.Provide(AsMiddleware(JWTAuthMiddleware))
	OptionSessionMiddleware      = fx.Provide(AsMiddleware(SessionMiddleware))
	OptionSiteSelectorMiddleware = fx.Provide(AsMiddleware(SiteSelectorMiddleware))
	OptionPageSelectorMiddleware = fx.Provide(AsMiddleware(PageSelectorMiddleware))
	OptionHybridPageMiddleware   = fx.Provide(AsMiddleware(HybridPageMiddleware))

	OptionHumaAuthorizationMiddleware = fx.Provide(AsHumaMiddleware(HumaAuthorizationMiddleware))
	OptionHumaAdminAuthAPI            = fx.Provide(
		fx.Annotate(
			NewAuthAPI,
			fx.As(new(HumaAPI)),
			fx.ParamTags(`name:"admin-auth-cache"`),
			fx.ResultTags(`group:"huma-admin-api"`),
		),
	)
	OptionHumaAdminAdminAPI         = fx.Provide(AsHumaAdminAPI(NewAdminAPI))
	OptionHumaAdminConfigurationAPI = fx.Provide(AsHumaAdminAPI(NewConfigurationAPI))
	OptionHumaAdminSiteAPI          = fx.Provide(AsHumaAdminAPI(NewSiteAPI))
	OptionHumaAdminPageAPI          = fx.Provide(AsHumaAdminAPI(NewPageAPI))
	OptionHumaAdminTemplateAPI      = fx.Provide(AsHumaAdminAPI(NewTemplateAPI))
	OptionHumaAdminMenuAPI          = fx.Provide(AsHumaAdminAPI(NewMenuAPI))
	OptionHumaAdminNodeAPI          = fx.Provide(AsHumaAdminAPI(NewNodeAPI))
)
