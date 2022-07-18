// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/NpoolPlatform/appuser-gateway/pkg/db/ent/appuserextra"
	"github.com/NpoolPlatform/appuser-gateway/pkg/db/ent/predicate"
)

// AppUserExtraDelete is the builder for deleting a AppUserExtra entity.
type AppUserExtraDelete struct {
	config
	hooks    []Hook
	mutation *AppUserExtraMutation
}

// Where appends a list predicates to the AppUserExtraDelete builder.
func (aued *AppUserExtraDelete) Where(ps ...predicate.AppUserExtra) *AppUserExtraDelete {
	aued.mutation.Where(ps...)
	return aued
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (aued *AppUserExtraDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(aued.hooks) == 0 {
		affected, err = aued.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*AppUserExtraMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			aued.mutation = mutation
			affected, err = aued.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(aued.hooks) - 1; i >= 0; i-- {
			if aued.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = aued.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, aued.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (aued *AppUserExtraDelete) ExecX(ctx context.Context) int {
	n, err := aued.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (aued *AppUserExtraDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: appuserextra.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: appuserextra.FieldID,
			},
		},
	}
	if ps := aued.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, aued.driver, _spec)
}

// AppUserExtraDeleteOne is the builder for deleting a single AppUserExtra entity.
type AppUserExtraDeleteOne struct {
	aued *AppUserExtraDelete
}

// Exec executes the deletion query.
func (auedo *AppUserExtraDeleteOne) Exec(ctx context.Context) error {
	n, err := auedo.aued.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{appuserextra.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (auedo *AppUserExtraDeleteOne) ExecX(ctx context.Context) {
	auedo.aued.ExecX(ctx)
}
