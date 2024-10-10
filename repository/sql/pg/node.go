package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gowool/cr"

	"github.com/gowool/cms/model"
	"github.com/gowool/cms/repository"
)

var _ repository.Node = (*NodeRepository)(nil)

type NodeRepository struct {
	Repository[model.Node, int64]
	tableSequence string
}

func NewNodeRepository(db *sql.DB) *NodeRepository {
	return &NodeRepository{
		tableSequence: "sequence_nodes",
		Repository: Repository[model.Node, int64]{
			DB:    db,
			Table: "nodes",
			SelectColumns: []string{
				"id", "parent_id", "name", "label", "uri", "path", "level", "position", "display_children",
				"display", "attributes", "link_attributes", "children_attributes", "label_attributes", "metadata",
				"created", "updated",
			},
			RowScan: func(row interface{ Scan(...any) error }, m *model.Node) error {
				var (
					label              sql.NullString
					uri                sql.NullString
					attributes         StrMap
					linkAttributes     StrMap
					childrenAttributes StrMap
					labelAttributes    StrMap
					metadata           StrMap
				)

				if err := row.Scan(&m.ID, &m.ParentID, &m.Name, &label, &uri, &m.Path, &m.Level, &m.Position,
					&m.DisplayChildren, &m.Display, &attributes, &linkAttributes, &childrenAttributes, &labelAttributes,
					&metadata, &m.Created, &m.Updated); err != nil {
					return err
				}
				m.Label = label.String
				m.URI = uri.String
				m.Attributes = attributes
				m.LinkAttributes = linkAttributes
				m.ChildrenAttributes = childrenAttributes
				m.LabelAttributes = labelAttributes
				m.Metadata = metadata
				return nil
			},
			InsertValues: func(m *model.Node) map[string]any {
				now := time.Now()
				return map[string]any{
					"id":                  m.ID,
					"parent_id":           m.ParentID,
					"name":                m.Name,
					"label":               sql.NullString{String: m.Label, Valid: m.Label != ""},
					"uri":                 sql.NullString{String: m.URI, Valid: m.URI != ""},
					"path":                m.Path,
					"level":               m.Level,
					"position":            m.Position,
					"display_children":    m.DisplayChildren,
					"display":             m.Display,
					"attributes":          StrMap(m.Attributes),
					"link_attributes":     StrMap(m.LinkAttributes),
					"children_attributes": StrMap(m.ChildrenAttributes),
					"label_attributes":    StrMap(m.LabelAttributes),
					"metadata":            StrMap(m.Metadata),
					"created":             now,
					"updated":             now,
				}
			},
			UpdateValues: func(m *model.Node) map[string]any {
				return map[string]any{
					"parent_id":           m.ParentID,
					"name":                m.Name,
					"label":               sql.NullString{String: m.Label, Valid: m.Label != ""},
					"uri":                 sql.NullString{String: m.URI, Valid: m.URI != ""},
					"path":                m.Path,
					"level":               m.Level,
					"position":            m.Position,
					"display_children":    m.DisplayChildren,
					"display":             m.Display,
					"attributes":          StrMap(m.Attributes),
					"link_attributes":     StrMap(m.LinkAttributes),
					"children_attributes": StrMap(m.ChildrenAttributes),
					"label_attributes":    StrMap(m.LabelAttributes),
					"metadata":            StrMap(m.Metadata),
					"updated":             time.Now(),
				}
			},
		},
	}
}

func (r *NodeRepository) FindWithChildren(ctx context.Context, id int64) ([]model.Node, error) {
	criteria := cr.New().
		SetSortBy(cr.ParseSort("path")...).
		SetFilter(cr.Filter{
			Operator: cr.OpOR,
			Conditions: []any{
				cr.Condition{Column: "path", Operator: cr.OpLIKE, Value: fmt.Sprintf("%%/%d", id)},
				cr.Condition{Column: "path", Operator: cr.OpLIKE, Value: fmt.Sprintf("%%/%d/%%", id)},
			},
		})

	return r.Find(ctx, criteria)
}

func (r *NodeRepository) Create(ctx context.Context, m *model.Node) error {
	if m == nil {
		panic("sql: Create called with nil pointer")
	}

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return r.error(err)
	}

	defer func() {
		if err == nil {
			err = r.error(tx.Commit())
		} else {
			err = errors.Join(err, r.error(tx.Rollback()))
		}
	}()

	query := fmt.Sprintf("INSERT INTO %s values (DEFAULT) RETURNING id", r.tableSequence)
	row := tx.QueryRowContext(ctx, query)
	if row.Err() != nil {
		return r.error(row.Err())
	}

	var id int64
	if err = row.Scan(&id); err != nil {
		return r.error(err)
	}
	m.ID = id

	if err = r.fixPath(ctx, m); err != nil {
		return err
	}
	return r.Repository.Create(WithTx(ctx, tx), m)
}

func (r *NodeRepository) Update(ctx context.Context, m *model.Node) error {
	if err := r.fixPath(ctx, m); err != nil {
		return err
	}
	return r.Repository.Update(ctx, m)
}

func (r *NodeRepository) Delete(ctx context.Context, ids ...int64) error {
	_, err := r.db(ctx).ExecContext(ctx, fmt.Sprintf(deleteSQL, r.tableSequence), ids)
	return r.error(err)
}

func (r *NodeRepository) fixPath(ctx context.Context, m *model.Node) (err error) {
	if m == nil {
		panic("sql: Update called with nil pointer")
	}

	if m.ParentID != 0 && m.Parent == nil {
		var parent model.Node
		if parent, err = r.FindByID(ctx, m.ParentID); err != nil {
			return err
		}
		m.Parent = &parent
	}

	*m = m.WithFixedPathAndLevel()
	return
}
