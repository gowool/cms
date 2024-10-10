package api

import (
	"github.com/gowool/cms/model"
	"github.com/gowool/cms/repository"
)

type NodeBody struct {
	ParentID           int64             `json:"parent_id,omitempty" yaml:"parent_id,omitempty" required:"false"`
	Name               string            `json:"name,omitempty" yaml:"name,omitempty" required:"true"`
	Label              string            `json:"label,omitempty" yaml:"label,omitempty" required:"false"`
	URI                string            `json:"uri,omitempty" yaml:"uri,omitempty" required:"false"`
	Position           int               `json:"position,omitempty" yaml:"position,omitempty" required:"false"`
	DisplayChildren    bool              `json:"display_children,omitempty" yaml:"display_children,omitempty" required:"false"`
	Display            bool              `json:"display,omitempty" yaml:"display,omitempty" required:"false"`
	Attributes         map[string]string `json:"attributes,omitempty" yaml:"attributes,omitempty" required:"false"`
	LinkAttributes     map[string]string `json:"link_attributes,omitempty" yaml:"link_attributes,omitempty" required:"false"`
	ChildrenAttributes map[string]string `json:"children_attributes,omitempty" yaml:"children_attributes,omitempty" required:"false"`
	LabelAttributes    map[string]string `json:"label_attributes,omitempty" yaml:"label_attributes,omitempty" required:"false"`
	Metadata           map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty" required:"false"`
}

func (dto NodeBody) Decode(m *model.Node) {
	m.ParentID = dto.ParentID
	m.Name = dto.Name
	m.Label = dto.Label
	m.URI = dto.URI
	m.Position = dto.Position
	m.DisplayChildren = dto.DisplayChildren
	m.Display = dto.Display
	m.Attributes = dto.Attributes
	m.LinkAttributes = dto.LinkAttributes
	m.ChildrenAttributes = dto.ChildrenAttributes
	m.LabelAttributes = dto.LabelAttributes
	m.Metadata = dto.Metadata
}

type Node struct {
	CRUD[NodeBody, model.Node, int64]
}

func NewNode(repo repository.Node, errorTransformer ErrorTransformerFunc) Node {
	return Node{
		CRUD: NewCRUD[NodeBody](repo, errorTransformer, "/nodes", "Node", "Nodes", "Node"),
	}
}
