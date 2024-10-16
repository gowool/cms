package cache

import (
	"context"
	"fmt"

	"github.com/gowool/cms"
	"github.com/gowool/cms/model"
	"github.com/gowool/cms/repository"
)

type NodeRepository struct {
	repository.Node
	repo[model.Node, int64]
}

func NewNodeRepository(inner repository.Node, c cms.Cache) NodeRepository {
	return NodeRepository{
		Node: inner,
		repo: repo[model.Node, int64]{inner: inner, cache: c, prefix: "cms::node"},
	}
}

func (r NodeRepository) FindByID(ctx context.Context, id int64) (model.Node, error) {
	return r.findByID(ctx, id)
}

func (r NodeRepository) Delete(ctx context.Context, ids ...int64) error {
	return r.delete(ctx, ids...)
}

func (r NodeRepository) Update(ctx context.Context, m *model.Node) error {
	defer r.del(ctx, m.ID)

	return r.Node.Update(ctx, m)
}

func (r NodeRepository) FindWithChildren(ctx context.Context, id int64) (nodes []model.Node, err error) {
	key := fmt.Sprintf("%s:with:children:%d", r.prefix, id)

	if err = r.cache.Get(ctx, key, &nodes); err == nil {
		return
	}

	if nodes, err = r.Node.FindWithChildren(ctx, id); err != nil {
		return
	}

	tags := make([]string, 0, len(nodes)+1)
	tags = append(tags, fmt.Sprintf("%s:tag:%d", r.prefix, id))

	for _, n := range nodes {
		tags = append(tags, fmt.Sprintf("%s:tag:%d", r.prefix, n.ID))
	}

	_ = r.cache.Set(ctx, key, nodes, tags...)
	return
}
