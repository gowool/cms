package api

import (
	"github.com/gowool/cms/model"
	"github.com/gowool/cms/repository"
)

type TemplateBody struct {
	Name    string `json:"name,omitempty" yaml:"name,omitempty" required:"true"`
	Content string `json:"content,omitempty" yaml:"content,omitempty" required:"false"`
	Enabled bool   `json:"enabled,omitempty" yaml:"enabled,omitempty" required:"false"`
}

func (dto TemplateBody) Decode(m *model.Template) {
	m.Name = dto.Name
	m.Content = dto.Content
	m.Enabled = dto.Enabled
}

type Template struct {
	CRUD[TemplateBody, model.Template, int64]
}

func NewTemplate(repo repository.Template, errorTransformer ErrorTransformerFunc) Template {
	return Template{
		CRUD: NewCRUD[TemplateBody](repo, errorTransformer, "/templates", "Template", "Templates", "Template"),
	}
}
