package repository

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"runtime/debug"

	"github.com/rafirandoni/hrengine/db"
	"github.com/rafirandoni/hrengine/utility"
	"github.com/uptrace/bun"
)

type BaseRepository[ENTITY any] struct {
	Table         string
	UseSoftDelete *bool
	DBClient      db.DBClient
}

func scanError(err error) error {
	if err == nil {
		return nil
	}

	if err == sql.ErrNoRows {
		return nil
	}

	return &utility.ResError{
		StatusCode: utility.ToPointer(http.StatusInternalServerError),
		Err:        err,
		StackTrace: utility.ToPointer(string(debug.Stack())),
	}
}

func (slf *BaseRepository[ENTITY]) FindAll(condition func(bun.QueryBuilder) bun.QueryBuilder, options ...map[string]any) ([]*ENTITY, error) {
	result := make([]*ENTITY, 0)
	query := slf.DBClient.ConnReader().NewSelect().
		Table(slf.Table)

	if condition != nil {
		query = query.ApplyQueryBuilder(condition)
	}

	withTrashed := false
	if len(options) > 0 {
		option := options[0]
		if orderBy, ok := option["order_by"]; ok {
			query = query.OrderExpr(orderBy.(string))
		}

		if limit, ok := option["limit"]; ok {
			query = query.Limit(limit.(int))
		}

		if offset, ok := option["offset"]; ok {
			query = query.Offset(offset.(int))
		}

		if group_by, ok := option["group_by"]; ok {
			query = query.GroupExpr(group_by.(string))
		}

		if _, ok := option["with_trashed"]; ok {
			withTrashed = true
		}
	}

	useSoftDelete := true
	if slf.UseSoftDelete != nil {
		useSoftDelete = *slf.UseSoftDelete
	}

	if !withTrashed && useSoftDelete {
		query = query.Where("deleted_at IS NULL")
	}

	err := query.Scan(context.Background(), &result)
	if err != nil {
		return nil, scanError(err)
	}

	return result, nil
}

func (slf *BaseRepository[ENTITY]) FindOne(condition func(bun.QueryBuilder) bun.QueryBuilder, options ...map[string]any) (*ENTITY, error) {
	if condition == nil {
		return nil, nil
	}

	result := new(ENTITY)
	query := slf.DBClient.ConnReader().NewSelect().
		Model(result).
		Limit(1).
		ApplyQueryBuilder(condition)

	withTrashed := false
	if len(options) > 0 {
		option := options[0]
		if orderBy, ok := option["order_by"]; ok {
			query = query.OrderExpr(orderBy.(string))
		}

		if limit, ok := option["limit"]; ok {
			query = query.Limit(limit.(int))
		}

		if offset, ok := option["offset"]; ok {
			query = query.Offset(offset.(int))
		}

		if group_by, ok := option["group_by"]; ok {
			query = query.GroupExpr(group_by.(string))
		}

		if _, ok := option["with_trashed"]; ok {
			withTrashed = true
		}
	}

	useSoftDelete := true
	if slf.UseSoftDelete != nil {
		useSoftDelete = *slf.UseSoftDelete
	}

	if !withTrashed && useSoftDelete {
		query = query.Where("deleted_at IS NULL")
	}

	err := query.Scan(context.Background())
	if err != nil {
		return nil, scanError(err)
	}

	return result, nil
}

func (slf *BaseRepository[ENTITY]) Count(condition func(bun.QueryBuilder) bun.QueryBuilder) (int, error) {
	query := slf.DBClient.ConnReader().NewSelect().
		Table(slf.Table)

	if condition != nil {
		query = query.ApplyQueryBuilder(condition)
	}

	count, err := query.Count(context.Background())
	if err != nil {
		return 0, scanError(err)
	}

	return count, nil
}

func (slf *BaseRepository[ENTITY]) Store(e *ENTITY) error {
	_, err := slf.DBClient.ConnWriter().NewInsert().
		Model(e).
		Exec(context.Background())
	if err != nil {
		return scanError(err)
	}

	return nil
}

func (slf *BaseRepository[ENTITY]) StoreBulk(e []*ENTITY) error {
	_, err := slf.DBClient.ConnWriter().NewInsert().
		Model(&e).
		Exec(context.Background())
	if err != nil {
		return scanError(err)
	}

	return nil
}

func (slf *BaseRepository[ENTITY]) StoreReturn(e *ENTITY, returnField string) (any, error) {
	var returnValue any

	_, err := slf.DBClient.ConnWriter().NewInsert().
		Model(e).
		Returning(returnField).
		Exec(context.Background(), &returnValue)
	if err != nil {
		return nil, scanError(err)
	}

	return returnValue, nil
}

func (slf *BaseRepository[ENTITY]) Update(condition func(bun.QueryBuilder) bun.QueryBuilder, param map[string]any) error {
	if condition == nil {
		return errors.New("cannot update without condition")
	}

	_, err := slf.DBClient.ConnWriter().NewUpdate().
		Table(slf.Table).
		Model(&param).
		ApplyQueryBuilder(condition).
		Exec(context.Background())

	return err
}

func (slf *BaseRepository[ENTITY]) Delete(condition func(bun.QueryBuilder) bun.QueryBuilder) error {
	useSoftDelete := true
	if slf.UseSoftDelete != nil {
		useSoftDelete = *slf.UseSoftDelete
	}

	if useSoftDelete {
		return slf.Update(condition, map[string]any{"deleted_at": utility.TimeNow()})
	}

	_, err := slf.DBClient.ConnWriter().NewDelete().
		Table(slf.Table).
		ApplyQueryBuilder(condition).
		Exec(context.Background())

	return err
}
