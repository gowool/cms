package api

import (
	"time"

	"github.com/gowool/cms/model"
	"github.com/gowool/cms/repository"
)

type SiteBody struct {
	Name         string            `json:"name,omitempty" yaml:"name,omitempty" required:"true"`
	Title        string            `json:"title,omitempty" yaml:"title,omitempty" required:"false"`
	Separator    string            `json:"separator,omitempty" yaml:"separator,omitempty" required:"true"`
	Host         string            `json:"host,omitempty" yaml:"host,omitempty" required:"true"`
	Locale       string            `json:"locale,omitempty" yaml:"locale,omitempty" required:"false"`
	RelativePath string            `json:"relative_path,omitempty" yaml:"relative_path,omitempty" required:"false"`
	IsDefault    bool              `json:"is_default,omitempty" yaml:"is_default,omitempty" required:"false"`
	Javascript   string            `json:"javascript,omitempty" yaml:"javascript,omitempty" required:"false"`
	Stylesheet   string            `json:"stylesheet,omitempty" yaml:"stylesheet,omitempty" required:"false"`
	Metas        []model.Meta      `json:"metas,omitempty" yaml:"metas,omitempty" required:"false"`
	Metadata     map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty" required:"false"`
	Published    *time.Time        `json:"published,omitempty" yaml:"published,omitempty" required:"false"`
	Expired      *time.Time        `json:"expired,omitempty" yaml:"expired,omitempty" required:"false"`
}

func (dto SiteBody) Decode(m *model.Site) {
	m.Name = dto.Name
	m.Title = dto.Title
	m.Separator = dto.Separator
	m.Host = dto.Host
	m.Locale = dto.Locale
	m.RelativePath = dto.RelativePath
	m.IsDefault = dto.IsDefault
	m.Javascript = dto.Javascript
	m.Stylesheet = dto.Stylesheet
	m.Metas = dto.Metas
	m.Metadata = dto.Metadata
	m.Published = dto.Published
	m.Expired = dto.Expired
}

type Site struct {
	CRUD[SiteBody, model.Site, int64]
}

func NewSite(repo repository.Site, errorTransformer ErrorTransformerFunc) Site {
	return Site{
		CRUD: NewCRUD[SiteBody](repo, errorTransformer, "/sites", "Site", "Sites", "Site"),
	}
}
