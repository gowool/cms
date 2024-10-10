package fx

import (
	"go.uber.org/zap"

	"github.com/gowool/cms"
	"github.com/gowool/cms/api"
	"github.com/gowool/cms/repository"
)

func NewAuthAPI(cache cms.Cache, r repository.Admin, cfg JWTConfig, logger *zap.Logger) api.Auth {
	return api.NewAuth(r, cache, cfg.Secret, cfg.AccessTokenDuration, logger)
}

func NewAdminAPI(r repository.Admin) api.Admin {
	return api.NewAdmin(r, api.ErrorTransformer)
}

func NewConfigurationAPI(r repository.Configuration) api.Configuration {
	return api.NewConfiguration(r, api.ErrorTransformer)
}

func NewPageAPI(r repository.Page, cfg repository.Configuration) api.Page {
	return api.NewPage(r, cfg, api.ErrorTransformer)
}

func NewSiteAPI(r repository.Site) api.Site {
	return api.NewSite(r, api.ErrorTransformer)
}

func NewTemplateAPI(r repository.Template) api.Template {
	return api.NewTemplate(r, api.ErrorTransformer)
}

func NewMenuAPI(r repository.Menu) api.Menu {
	return api.NewMenu(r, api.ErrorTransformer)
}

func NewNodeAPI(r repository.Node) api.Node {
	return api.NewNode(r, api.ErrorTransformer)
}
