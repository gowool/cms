package api

import (
	"github.com/gowool/cms/model"
	"github.com/gowool/cms/repository"
)

type MenuBody struct {
	NodeID  *int64 `json:"node_id,omitempty" yaml:"node_id,omitempty" required:"false"`
	Name    string `json:"name,omitempty" yaml:"name,omitempty" required:"true"`
	Handle  string `json:"handle,omitempty" yaml:"handle,omitempty" required:"false"`
	Enabled bool   `json:"enabled,omitempty" yaml:"enabled,omitempty" required:"false"`
}

func (dto MenuBody) Decode(m *model.Menu) {
	m.NodeID = dto.NodeID
	m.Name = dto.Name
	m.Handle = dto.Handle
	m.Enabled = dto.Enabled
}

type Menu struct {
	CRUD[MenuBody, model.Menu, int64]
}

func NewMenu(repo repository.Menu, errorTransformer ErrorTransformerFunc) Menu {
	return Menu{
		CRUD: NewCRUD[MenuBody](repo, errorTransformer, "/menus", "Menu", "Menus", "Menu"),
	}
}
