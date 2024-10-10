package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gowool/cr"
	"github.com/labstack/echo/v4"

	"github.com/gowool/cms"
	"github.com/gowool/cms/repository"
)

var Security = []map[string][]string{
	{"BasicAuth": {}},
	{"AuthJWT": {}},
}

type ErrorTransformerFunc func(_ context.Context, err error) error

type Response[B any] struct {
	Body B
}

type CreateResponse struct {
	Location string `header:"Content-Location"`
}

func Location[ID any](path string, id ID) *CreateResponse {
	return &CreateResponse{Location: strings.ReplaceAll(path, "{id}", fmt.Sprintf("%v", id))}
}

type CreateInput[B any] struct {
	Body B
}

type Create[B interface{ Decode(*M) }, M interface{ GetID() ID }, ID any] struct {
	Saver            func(context.Context, *M) error
	ErrorTransformer ErrorTransformerFunc
	PathID           string
}

func NewCreate[B interface{ Decode(*M) }, M interface{ GetID() ID }, ID any](
	saver func(context.Context, *M) error,
	errorTransformer ErrorTransformerFunc,
	pathID string,
) Create[B, M, ID] {
	return Create[B, M, ID]{Saver: saver, ErrorTransformer: errorTransformer, PathID: pathID}
}

func (h Create[B, M, ID]) Handler(ctx context.Context, in *CreateInput[B]) (*CreateResponse, error) {
	var m M
	in.Body.Decode(&m)
	if err := h.Saver(ctx, &m); err != nil {
		return nil, h.ErrorTransformer(ctx, err)
	}
	return Location(h.PathID, m.GetID()), nil
}

type UpdateInput[B any, ID any] struct {
	ID   ID `path:"id"`
	Body B
}

type Update[B interface{ Decode(*M) }, M any, ID any] struct {
	Finder           func(context.Context, ID) (M, error)
	Saver            func(context.Context, *M) error
	ErrorTransformer ErrorTransformerFunc
}

func NewUpdate[B interface{ Decode(*M) }, M any, ID any](
	finder func(context.Context, ID) (M, error),
	saver func(context.Context, *M) error,
	errorTransformer ErrorTransformerFunc,
) Update[B, M, ID] {
	return Update[B, M, ID]{Finder: finder, Saver: saver, ErrorTransformer: errorTransformer}
}

func (h Update[B, M, ID]) Handler(ctx context.Context, in *UpdateInput[B, ID]) (*struct{}, error) {
	m, err := h.Finder(ctx, in.ID)
	if err != nil {
		return nil, h.ErrorTransformer(ctx, err)
	}

	in.Body.Decode(&m)
	if err = h.Saver(ctx, &m); err != nil {
		return nil, h.ErrorTransformer(ctx, err)
	}
	return nil, nil
}

type ListInput struct {
	Page   int    `query:"page" json:"page,omitempty" yaml:"page,omitempty" required:"false"`
	Limit  int    `query:"limit" json:"limit,omitempty" yaml:"limit,omitempty" required:"false"`
	Sort   string `query:"sort" json:"sort,omitempty" yaml:"sort,omitempty" required:"false"`
	Filter string `query:"filter" json:"filter,omitempty" yaml:"filter,omitempty" required:"false"`
}

func (in *ListInput) Resolve(huma.Context) []error {
	if in.Page < 1 {
		in.Page = 1
	}
	if in.Limit < 1 || in.Limit > 100 {
		in.Limit = 100
	}
	if in.Sort == "" {
		in.Sort = "-id"
	}
	return nil
}

func (in *ListInput) criteria() *cr.Criteria {
	return cr.New(in.Filter, in.Sort).SetOffset((in.Page - 1) * in.Limit).SetSize(in.Limit)
}

type ListOutput[E any] struct {
	ListInput
	Items []E `json:"items,omitempty" yaml:"items,omitempty" required:"false"`
	Total int `json:"total,omitempty" yaml:"total,omitempty" required:"false"`
}

type List[M any] struct {
	Finder           func(context.Context, *cr.Criteria) ([]M, int, error)
	ErrorTransformer ErrorTransformerFunc
}

func NewList[T any](finder func(context.Context, *cr.Criteria) ([]T, int, error), errorTransformer ErrorTransformerFunc) List[T] {
	return List[T]{Finder: finder, ErrorTransformer: errorTransformer}
}

func (h List[M]) Handler(ctx context.Context, in *ListInput) (*Response[ListOutput[M]], error) {
	items, total, err := h.Finder(ctx, in.criteria())
	if err != nil {
		return nil, h.ErrorTransformer(ctx, err)
	}
	return &Response[ListOutput[M]]{
		Body: ListOutput[M]{
			ListInput: *in,
			Items:     items,
			Total:     total,
		},
	}, nil
}

type IDInput[ID any] struct {
	ID ID `path:"id"`
}

type Read[M any, ID any] struct {
	Finder           func(context.Context, ID) (M, error)
	ErrorTransformer ErrorTransformerFunc
}

func NewRead[T any, ID any](finder func(context.Context, ID) (T, error), errorTransformer ErrorTransformerFunc) Read[T, ID] {
	return Read[T, ID]{Finder: finder, ErrorTransformer: errorTransformer}
}

func (h Read[M, ID]) Handler(ctx context.Context, in *IDInput[ID]) (*Response[M], error) {
	item, err := h.Finder(ctx, in.ID)
	if err != nil {
		return nil, h.ErrorTransformer(ctx, err)
	}
	return &Response[M]{Body: item}, nil
}

type Delete[ID any] struct {
	Deleter          func(context.Context, ...ID) error
	ErrorTransformer ErrorTransformerFunc
}

func NewDelete[ID any](deleter func(context.Context, ...ID) error, errorTransformer ErrorTransformerFunc) Delete[ID] {
	return Delete[ID]{Deleter: deleter, ErrorTransformer: errorTransformer}
}

func (h Delete[ID]) Handler(ctx context.Context, in *IDInput[ID]) (*struct{}, error) {
	if err := h.Deleter(ctx, in.ID); err != nil {
		return nil, h.ErrorTransformer(ctx, err)
	}
	return nil, nil
}

type IDsInput[ID any] struct {
	Body struct {
		IDs []ID `json:"ids" required:"true" minItems:"1" nullable:"false"`
	}
}

type DeleteMany[ID any] struct {
	Deleter          func(context.Context, ...ID) error
	ErrorTransformer ErrorTransformerFunc
}

func NewDeleteMany[ID any](deleter func(context.Context, ...ID) error, errorTransformer ErrorTransformerFunc) DeleteMany[ID] {
	return DeleteMany[ID]{Deleter: deleter, ErrorTransformer: errorTransformer}
}

func (h DeleteMany[ID]) Handler(ctx context.Context, in *IDsInput[ID]) (*struct{}, error) {
	if err := h.Deleter(ctx, in.Body.IDs...); err != nil {
		return nil, h.ErrorTransformer(ctx, err)
	}
	return nil, nil
}

func ErrorTransformer(_ context.Context, err error) error {
	var statusErr huma.StatusError
	if errors.As(err, &statusErr) {
		return statusErr
	}

	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, repository.ErrNotFound) ||
		errors.Is(err, repository.ErrSiteNotFound) || errors.Is(err, repository.ErrPageNotFound) {
		return huma.Error404NotFound("Not Found", err)
	}

	if errors.Is(err, repository.ErrUniqueViolation) {
		return huma.Error409Conflict("Conflict", err)
	}

	return huma.Error500InternalServerError("Internal Server Error", err)
}

func Register[I, O any](api huma.API, handler func(context.Context, *I) (*O, error), operation huma.Operation) {
	var o *O
	if operation.OperationID == "" {
		operation.OperationID = huma.GenerateOperationID(operation.Method, operation.Path, o)
	}
	if operation.Summary == "" {
		operation.Summary = huma.GenerateSummary(operation.Method, operation.Path, o)
	}
	if operation.Security == nil {
		operation.Security = Security
	}
	if operation.Metadata == nil {
		operation.Metadata = make(map[string]any)
	}
	if target, ok := operation.Metadata["target"].(*cms.CallTarget); ok {
		if target.OperationID == "" {
			target.OperationID = operation.OperationID
		}
	}
	huma.Register(api, operation, handler)
}

type CRUD[B interface{ Decode(*M) }, M interface{ GetID() ID }, ID any] struct {
	List[M]
	Read[M, ID]
	Create[B, M, ID]
	Update[B, M, ID]
	Delete[ID]
	DeleteMany[ID]

	Path          string
	LabelSingular string
	LabelPlural   string
	Tags          []string
}

func NewCRUD[B interface{ Decode(*M) }, M interface{ GetID() ID }, ID any](
	repo repository.Repository[M, ID],
	errorTransformer ErrorTransformerFunc,
	path,
	labelSingular,
	labelPlural string,
	tags ...string,
) CRUD[B, M, ID] {
	return CRUD[B, M, ID]{
		List:          NewList(repo.FindAndCount, errorTransformer),
		Read:          NewRead(repo.FindByID, errorTransformer),
		Create:        NewCreate[B](repo.Create, errorTransformer, path+"/{id}"),
		Update:        NewUpdate[B](repo.FindByID, repo.Update, errorTransformer),
		Delete:        NewDelete(repo.Delete, errorTransformer),
		DeleteMany:    NewDeleteMany(repo.Delete, errorTransformer),
		Path:          path,
		LabelSingular: labelSingular,
		LabelPlural:   labelPlural,
		Tags:          tags,
	}
}

func (h CRUD[B, M, ID]) Register(_ *echo.Echo, api huma.API) {
	Register(api, h.List.Handler, huma.Operation{
		Summary: "Get " + h.LabelPlural,
		Method:  http.MethodGet,
		Path:    h.Path,
		Tags:    h.Tags,
		Metadata: map[string]any{
			"target": cms.NewCallTarget(cms.AccessRead),
		},
	})
	Register(api, h.Read.Handler, huma.Operation{
		Summary: "Get " + h.LabelSingular,
		Method:  http.MethodGet,
		Path:    h.PathID,
		Tags:    h.Tags,
		Metadata: map[string]any{
			"target": cms.NewCallTarget(cms.AccessRead),
		},
	})
	Register(api, h.DeleteMany.Handler, huma.Operation{
		Summary: "Delete " + h.LabelPlural,
		Method:  http.MethodDelete,
		Path:    h.Path,
		Tags:    h.Tags,
		Metadata: map[string]any{
			"target": cms.NewCallTarget(cms.AccessWrite),
		},
	})
	Register(api, h.Delete.Handler, huma.Operation{
		Summary: "Delete " + h.LabelSingular,
		Method:  http.MethodDelete,
		Path:    h.PathID,
		Tags:    h.Tags,
		Metadata: map[string]any{
			"target": cms.NewCallTarget(cms.AccessWrite),
		},
	})
	Register(api, h.Create.Handler, huma.Operation{
		Summary:       "Create " + h.LabelSingular,
		DefaultStatus: http.StatusCreated,
		Method:        http.MethodPost,
		Path:          h.Path,
		Tags:          h.Tags,
		Metadata: map[string]any{
			"target": cms.NewCallTarget(cms.AccessWrite),
		},
	})
	Register(api, h.Update.Handler, huma.Operation{
		Summary: "Update " + h.LabelSingular,
		Method:  http.MethodPut,
		Path:    h.PathID,
		Tags:    h.Tags,
		Metadata: map[string]any{
			"target": cms.NewCallTarget(cms.AccessWrite),
		},
	})
}
