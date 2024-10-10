package fx

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gowool/cms"
	"github.com/gowool/cms/repository"
)

type SeederParams struct {
	fx.In
	Lifecycle      fx.Lifecycle
	SiteRepository repository.Site
	PageRepository repository.Page
	Logger         *zap.Logger
}

func NewSeeder(params SeederParams) cms.Seeder {
	seeder := cms.NewDefaultSeeder(params.SiteRepository, params.PageRepository, params.Logger)

	params.Lifecycle.Append(fx.StartHook(seeder.Boot))

	return seeder
}
