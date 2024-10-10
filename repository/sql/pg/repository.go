package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/gowool/cr"

	"github.com/gowool/cms/repository"
)

const (
	selectSQL       = "SELECT %s FROM %s"
	selectOneSQL    = "SELECT %s FROM %s WHERE %s = $1 LIMIT 1"
	countSQL        = "SELECT COUNT(*) FROM %s"
	insertSQL       = "INSERT INTO %s (%s) VALUES (%s) RETURNING %s"
	updateSQL       = "UPDATE %s SET %s WHERE id = $%d RETURNING %s"
	deleteSQL       = "DELETE FROM %s WHERE id = ANY($1)"
	uniqueViolation = "23505"
)

func WithTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, ctxTxKey{}, tx)
}

type (
	ctxTxKey struct{}
	txDB     interface {
		QueryContext(context.Context, string, ...any) (*sql.Rows, error)
		QueryRowContext(context.Context, string, ...any) *sql.Row
		ExecContext(context.Context, string, ...any) (sql.Result, error)
	}
)

type Repository[T interface{ GetID() ID }, ID any] struct {
	DB            *sql.DB
	Table         string
	SelectColumns []string
	RowScan       func(interface{ Scan(dest ...any) error }, *T) error
	InsertValues  func(*T) map[string]any
	UpdateValues  func(*T) map[string]any
	OnError       func(error) error
}

func (r Repository[T, ID]) FindAndCount(ctx context.Context, criteria *cr.Criteria) ([]T, int, error) {
	if criteria == nil {
		criteria = cr.New()
	}

	where, args := r.where(criteria)

	var total int
	if err := r.db(ctx).QueryRowContext(ctx, fmt.Sprintf(countSQL, r.Table)+where, args...).Scan(&total); err != nil {
		return nil, 0, r.error(err)
	}
	if total == 0 {
		return nil, 0, nil
	}

	data, err := r.find(ctx, criteria, where, args)
	if err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

func (r Repository[T, ID]) Find(ctx context.Context, criteria *cr.Criteria) ([]T, error) {
	if criteria == nil {
		criteria = cr.New()
	}

	where, args := r.where(criteria)

	return r.find(ctx, criteria, where, args)
}

func (r Repository[T, ID]) find(ctx context.Context, criteria *cr.Criteria, where string, args []any) ([]T, error) {
	var (
		query strings.Builder
		data  []T
		size  int
	)

	index := len(args)
	query.WriteString(where)

	if len(criteria.SortBy) > 0 {
		query.WriteString(" ORDER BY ")
		query.WriteString(criteria.SortBy.String())
	}

	if criteria.Size != nil && *criteria.Size > 0 {
		query.WriteString(" LIMIT $")
		query.WriteString(strconv.Itoa(index + 1))
		query.WriteString(" OFFSET $")
		query.WriteString(strconv.Itoa(index + 2))

		size = *criteria.Size
		args = append(args, size, criteria.GetOffset())
	}

	rows, err := r.db(ctx).QueryContext(ctx, fmt.Sprintf(selectSQL, r.columns(), r.Table)+query.String(), args...)
	if err != nil {
		return nil, r.error(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	data = make([]T, 0, size)
	for rows.Next() {
		var item T
		if err = r.RowScan(rows, &item); err != nil {
			return nil, r.error(err)
		}
		data = append(data, item)
	}
	return slices.Clip(data), nil
}

func (r Repository[T, ID]) where(criteria *cr.Criteria) (string, []any) {
	s, args := criteria.Filter.ToSQL()
	if s == "" {
		return "", nil
	}

	var (
		where strings.Builder
		index int
	)

	where.WriteString(" WHERE ")
	for i := 0; i < len(s); i++ {
		if i > 0 && (i+2) < len(s) && s[i-1] == ' ' &&
			(s[i] == 'I' || s[i] == 'i') &&
			(s[i+1] == 'N' || s[i+1] == 'n') &&
			(s[i+2] == '?' || s[i+2] == ' ' || s[i+2] == '(') {
			where.WriteString("= ANY")
			i++
			continue
		}

		if s[i] == '?' {
			index++
			where.WriteString(fmt.Sprintf("$%d", index))
			continue
		}

		where.WriteByte(s[i])
	}

	return where.String(), args
}

func (r Repository[T, ID]) FindByID(ctx context.Context, id ID) (T, error) {
	m, err := r.FindBy(ctx, "id", id)
	return m, r.error(err)
}

func (r Repository[T, ID]) FindBy(ctx context.Context, column string, value any) (m T, err error) {
	query := fmt.Sprintf(selectOneSQL, r.columns(), r.Table, column)
	row := r.db(ctx).QueryRowContext(ctx, query, value)
	err = r.error(r.RowScan(row, &m))
	return
}

func (r Repository[T, ID]) Delete(ctx context.Context, ids ...ID) error {
	_, err := r.db(ctx).ExecContext(ctx, fmt.Sprintf(deleteSQL, r.Table), ids)
	return r.error(err)
}

func (r Repository[T, ID]) Create(ctx context.Context, m *T) error {
	if m == nil {
		panic("sql: Create called with nil pointer")
	}

	data := r.InsertValues(m)
	columns := make([]string, 0, len(data))
	values := make([]string, 0, len(data))
	args := make([]any, 0, len(data))

	for column, value := range data {
		columns = append(columns, column)
		args = append(args, value)
		values = append(values, fmt.Sprintf("$%d", len(args)))
	}

	query := fmt.Sprintf(insertSQL, r.Table, strings.Join(columns, ","), strings.Join(values, ","), r.columns())

	row := r.db(ctx).QueryRowContext(ctx, query, args...)
	return r.error(r.RowScan(row, m))
}

func (r Repository[T, ID]) Update(ctx context.Context, m *T) error {
	if m == nil {
		panic("sql: Update called with nil pointer")
	}

	data := r.UpdateValues(m)
	columns := make([]string, 0, len(data))
	args := make([]any, 0, len(data)+1)

	for column, value := range data {
		args = append(args, value)
		columns = append(columns, fmt.Sprintf("%s = $%d", column, len(args)))
	}

	args = append(args, (*m).GetID())
	query := fmt.Sprintf(updateSQL, r.Table, strings.Join(columns, ","), len(args), r.columns())

	row := r.db(ctx).QueryRowContext(ctx, query, args...)
	return r.error(r.RowScan(row, m))
}

func (r Repository[T, ID]) columns() string {
	if len(r.SelectColumns) == 0 {
		return "*"
	}
	return strings.Join(r.SelectColumns, ",")
}

func (r Repository[T, ID]) error(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		err = errors.Join(err, repository.ErrNotFound)
	} else if strings.Contains(err.Error(), uniqueViolation) {
		err = errors.Join(err, repository.ErrUniqueViolation)
	}
	if r.OnError != nil {
		return r.OnError(err)
	}
	return err
}

func (r Repository[T, ID]) db(ctx context.Context) txDB {
	var db txDB = r.DB
	if tx, ok := ctx.Value(ctxTxKey{}).(*sql.Tx); ok {
		db = tx
	}
	return db
}
