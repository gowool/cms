package fx

import (
	"context"
	"database/sql"
	"io/fs"

	"github.com/gowool/theme"

	"github.com/gowool/cms"
	"github.com/gowool/cms/model"
	"github.com/gowool/cms/repository"
	cacherepo "github.com/gowool/cms/repository/cache"
	"github.com/gowool/cms/repository/fallback"
	fsrepo "github.com/gowool/cms/repository/fs"
	"github.com/gowool/cms/repository/sql/pg"
)

func NewAdminRepository(db *sql.DB) repository.Admin {
	return pg.NewAdminRepository(db)
}

type TemplateRepositoryParams struct {
	Debug bool
	Cache cms.Cache `name:"repository-cache"`
	DB    *sql.DB
	FSS   []fs.FS
}

func NewTemplateRepository(params TemplateRepositoryParams) repository.Template {
	var r repository.Template = pg.NewTemplateRepository(params.DB)
	for _, fsys := range params.FSS {
		r = fsrepo.NewTemplateRepository(r, fsys)
	}

	if params.Debug {
		return r
	}
	return cacherepo.NewTemplateRepository(r, params.Cache)
}

func NewSiteRepository(db *sql.DB, c cms.Cache) repository.Site {
	r := pg.NewSiteRepository(db)
	return cacherepo.NewSiteRepository(r, c)
}

func NewPageRepository(db *sql.DB, c cms.Cache) repository.Page {
	r := pg.NewPageRepository(db)
	return cacherepo.NewPageRepository(r, c)
}

func NewConfigurationRepository(db *sql.DB, c cms.Cache) repository.Configuration {
	var r repository.Configuration = pg.NewConfigurationRepository(db)
	r = fallback.NewConfigurationRepository(r, model.NewConfiguration())
	return cacherepo.NewConfigurationRepository(r, c)
}

func NewMenuRepository(db *sql.DB, c cms.Cache) repository.Menu {
	r := pg.NewMenuRepository(db)
	return cacherepo.NewMenuRepository(r, c)
}

func NewNodeRepository(db *sql.DB, c cms.Cache) repository.Node {
	r := pg.NewNodeRepository(db)
	return cacherepo.NewNodeRepository(r, c)
}

type ThemeRepository struct {
	r repository.Template
}

func NewThemeRepository(r repository.Template) theme.Repository {
	return ThemeRepository{r: r}
}

func (r ThemeRepository) FindByName(ctx context.Context, name string) (theme.Template, error) {
	return r.r.FindByName(ctx, name)
}
