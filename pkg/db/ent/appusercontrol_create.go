// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/NpoolPlatform/appuser-gateway/pkg/db/ent/appusercontrol"
	"github.com/google/uuid"
)

// AppUserControlCreate is the builder for creating a AppUserControl entity.
type AppUserControlCreate struct {
	config
	mutation *AppUserControlMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreatedAt sets the "created_at" field.
func (aucc *AppUserControlCreate) SetCreatedAt(u uint32) *AppUserControlCreate {
	aucc.mutation.SetCreatedAt(u)
	return aucc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (aucc *AppUserControlCreate) SetNillableCreatedAt(u *uint32) *AppUserControlCreate {
	if u != nil {
		aucc.SetCreatedAt(*u)
	}
	return aucc
}

// SetUpdatedAt sets the "updated_at" field.
func (aucc *AppUserControlCreate) SetUpdatedAt(u uint32) *AppUserControlCreate {
	aucc.mutation.SetUpdatedAt(u)
	return aucc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (aucc *AppUserControlCreate) SetNillableUpdatedAt(u *uint32) *AppUserControlCreate {
	if u != nil {
		aucc.SetUpdatedAt(*u)
	}
	return aucc
}

// SetDeletedAt sets the "deleted_at" field.
func (aucc *AppUserControlCreate) SetDeletedAt(u uint32) *AppUserControlCreate {
	aucc.mutation.SetDeletedAt(u)
	return aucc
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (aucc *AppUserControlCreate) SetNillableDeletedAt(u *uint32) *AppUserControlCreate {
	if u != nil {
		aucc.SetDeletedAt(*u)
	}
	return aucc
}

// SetAppID sets the "app_id" field.
func (aucc *AppUserControlCreate) SetAppID(u uuid.UUID) *AppUserControlCreate {
	aucc.mutation.SetAppID(u)
	return aucc
}

// SetUserID sets the "user_id" field.
func (aucc *AppUserControlCreate) SetUserID(u uuid.UUID) *AppUserControlCreate {
	aucc.mutation.SetUserID(u)
	return aucc
}

// SetSigninVerifyByGoogleAuthentication sets the "signin_verify_by_google_authentication" field.
func (aucc *AppUserControlCreate) SetSigninVerifyByGoogleAuthentication(b bool) *AppUserControlCreate {
	aucc.mutation.SetSigninVerifyByGoogleAuthentication(b)
	return aucc
}

// SetGoogleAuthenticationVerified sets the "google_authentication_verified" field.
func (aucc *AppUserControlCreate) SetGoogleAuthenticationVerified(b bool) *AppUserControlCreate {
	aucc.mutation.SetGoogleAuthenticationVerified(b)
	return aucc
}

// SetID sets the "id" field.
func (aucc *AppUserControlCreate) SetID(u uuid.UUID) *AppUserControlCreate {
	aucc.mutation.SetID(u)
	return aucc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (aucc *AppUserControlCreate) SetNillableID(u *uuid.UUID) *AppUserControlCreate {
	if u != nil {
		aucc.SetID(*u)
	}
	return aucc
}

// Mutation returns the AppUserControlMutation object of the builder.
func (aucc *AppUserControlCreate) Mutation() *AppUserControlMutation {
	return aucc.mutation
}

// Save creates the AppUserControl in the database.
func (aucc *AppUserControlCreate) Save(ctx context.Context) (*AppUserControl, error) {
	var (
		err  error
		node *AppUserControl
	)
	if err := aucc.defaults(); err != nil {
		return nil, err
	}
	if len(aucc.hooks) == 0 {
		if err = aucc.check(); err != nil {
			return nil, err
		}
		node, err = aucc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*AppUserControlMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = aucc.check(); err != nil {
				return nil, err
			}
			aucc.mutation = mutation
			if node, err = aucc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(aucc.hooks) - 1; i >= 0; i-- {
			if aucc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = aucc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, aucc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (aucc *AppUserControlCreate) SaveX(ctx context.Context) *AppUserControl {
	v, err := aucc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (aucc *AppUserControlCreate) Exec(ctx context.Context) error {
	_, err := aucc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (aucc *AppUserControlCreate) ExecX(ctx context.Context) {
	if err := aucc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (aucc *AppUserControlCreate) defaults() error {
	if _, ok := aucc.mutation.CreatedAt(); !ok {
		if appusercontrol.DefaultCreatedAt == nil {
			return fmt.Errorf("ent: uninitialized appusercontrol.DefaultCreatedAt (forgotten import ent/runtime?)")
		}
		v := appusercontrol.DefaultCreatedAt()
		aucc.mutation.SetCreatedAt(v)
	}
	if _, ok := aucc.mutation.UpdatedAt(); !ok {
		if appusercontrol.DefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized appusercontrol.DefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := appusercontrol.DefaultUpdatedAt()
		aucc.mutation.SetUpdatedAt(v)
	}
	if _, ok := aucc.mutation.DeletedAt(); !ok {
		if appusercontrol.DefaultDeletedAt == nil {
			return fmt.Errorf("ent: uninitialized appusercontrol.DefaultDeletedAt (forgotten import ent/runtime?)")
		}
		v := appusercontrol.DefaultDeletedAt()
		aucc.mutation.SetDeletedAt(v)
	}
	if _, ok := aucc.mutation.ID(); !ok {
		if appusercontrol.DefaultID == nil {
			return fmt.Errorf("ent: uninitialized appusercontrol.DefaultID (forgotten import ent/runtime?)")
		}
		v := appusercontrol.DefaultID()
		aucc.mutation.SetID(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (aucc *AppUserControlCreate) check() error {
	if _, ok := aucc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "AppUserControl.created_at"`)}
	}
	if _, ok := aucc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "AppUserControl.updated_at"`)}
	}
	if _, ok := aucc.mutation.DeletedAt(); !ok {
		return &ValidationError{Name: "deleted_at", err: errors.New(`ent: missing required field "AppUserControl.deleted_at"`)}
	}
	if _, ok := aucc.mutation.AppID(); !ok {
		return &ValidationError{Name: "app_id", err: errors.New(`ent: missing required field "AppUserControl.app_id"`)}
	}
	if _, ok := aucc.mutation.UserID(); !ok {
		return &ValidationError{Name: "user_id", err: errors.New(`ent: missing required field "AppUserControl.user_id"`)}
	}
	if _, ok := aucc.mutation.SigninVerifyByGoogleAuthentication(); !ok {
		return &ValidationError{Name: "signin_verify_by_google_authentication", err: errors.New(`ent: missing required field "AppUserControl.signin_verify_by_google_authentication"`)}
	}
	if _, ok := aucc.mutation.GoogleAuthenticationVerified(); !ok {
		return &ValidationError{Name: "google_authentication_verified", err: errors.New(`ent: missing required field "AppUserControl.google_authentication_verified"`)}
	}
	return nil
}

func (aucc *AppUserControlCreate) sqlSave(ctx context.Context) (*AppUserControl, error) {
	_node, _spec := aucc.createSpec()
	if err := sqlgraph.CreateNode(ctx, aucc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	return _node, nil
}

func (aucc *AppUserControlCreate) createSpec() (*AppUserControl, *sqlgraph.CreateSpec) {
	var (
		_node = &AppUserControl{config: aucc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: appusercontrol.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: appusercontrol.FieldID,
			},
		}
	)
	_spec.OnConflict = aucc.conflict
	if id, ok := aucc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := aucc.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: appusercontrol.FieldCreatedAt,
		})
		_node.CreatedAt = value
	}
	if value, ok := aucc.mutation.UpdatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: appusercontrol.FieldUpdatedAt,
		})
		_node.UpdatedAt = value
	}
	if value, ok := aucc.mutation.DeletedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: appusercontrol.FieldDeletedAt,
		})
		_node.DeletedAt = value
	}
	if value, ok := aucc.mutation.AppID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeUUID,
			Value:  value,
			Column: appusercontrol.FieldAppID,
		})
		_node.AppID = value
	}
	if value, ok := aucc.mutation.UserID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeUUID,
			Value:  value,
			Column: appusercontrol.FieldUserID,
		})
		_node.UserID = value
	}
	if value, ok := aucc.mutation.SigninVerifyByGoogleAuthentication(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: appusercontrol.FieldSigninVerifyByGoogleAuthentication,
		})
		_node.SigninVerifyByGoogleAuthentication = value
	}
	if value, ok := aucc.mutation.GoogleAuthenticationVerified(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: appusercontrol.FieldGoogleAuthenticationVerified,
		})
		_node.GoogleAuthenticationVerified = value
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.AppUserControl.Create().
//		SetCreatedAt(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.AppUserControlUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
//
func (aucc *AppUserControlCreate) OnConflict(opts ...sql.ConflictOption) *AppUserControlUpsertOne {
	aucc.conflict = opts
	return &AppUserControlUpsertOne{
		create: aucc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.AppUserControl.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
//
func (aucc *AppUserControlCreate) OnConflictColumns(columns ...string) *AppUserControlUpsertOne {
	aucc.conflict = append(aucc.conflict, sql.ConflictColumns(columns...))
	return &AppUserControlUpsertOne{
		create: aucc,
	}
}

type (
	// AppUserControlUpsertOne is the builder for "upsert"-ing
	//  one AppUserControl node.
	AppUserControlUpsertOne struct {
		create *AppUserControlCreate
	}

	// AppUserControlUpsert is the "OnConflict" setter.
	AppUserControlUpsert struct {
		*sql.UpdateSet
	}
)

// SetCreatedAt sets the "created_at" field.
func (u *AppUserControlUpsert) SetCreatedAt(v uint32) *AppUserControlUpsert {
	u.Set(appusercontrol.FieldCreatedAt, v)
	return u
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *AppUserControlUpsert) UpdateCreatedAt() *AppUserControlUpsert {
	u.SetExcluded(appusercontrol.FieldCreatedAt)
	return u
}

// AddCreatedAt adds v to the "created_at" field.
func (u *AppUserControlUpsert) AddCreatedAt(v uint32) *AppUserControlUpsert {
	u.Add(appusercontrol.FieldCreatedAt, v)
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *AppUserControlUpsert) SetUpdatedAt(v uint32) *AppUserControlUpsert {
	u.Set(appusercontrol.FieldUpdatedAt, v)
	return u
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *AppUserControlUpsert) UpdateUpdatedAt() *AppUserControlUpsert {
	u.SetExcluded(appusercontrol.FieldUpdatedAt)
	return u
}

// AddUpdatedAt adds v to the "updated_at" field.
func (u *AppUserControlUpsert) AddUpdatedAt(v uint32) *AppUserControlUpsert {
	u.Add(appusercontrol.FieldUpdatedAt, v)
	return u
}

// SetDeletedAt sets the "deleted_at" field.
func (u *AppUserControlUpsert) SetDeletedAt(v uint32) *AppUserControlUpsert {
	u.Set(appusercontrol.FieldDeletedAt, v)
	return u
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *AppUserControlUpsert) UpdateDeletedAt() *AppUserControlUpsert {
	u.SetExcluded(appusercontrol.FieldDeletedAt)
	return u
}

// AddDeletedAt adds v to the "deleted_at" field.
func (u *AppUserControlUpsert) AddDeletedAt(v uint32) *AppUserControlUpsert {
	u.Add(appusercontrol.FieldDeletedAt, v)
	return u
}

// SetAppID sets the "app_id" field.
func (u *AppUserControlUpsert) SetAppID(v uuid.UUID) *AppUserControlUpsert {
	u.Set(appusercontrol.FieldAppID, v)
	return u
}

// UpdateAppID sets the "app_id" field to the value that was provided on create.
func (u *AppUserControlUpsert) UpdateAppID() *AppUserControlUpsert {
	u.SetExcluded(appusercontrol.FieldAppID)
	return u
}

// SetUserID sets the "user_id" field.
func (u *AppUserControlUpsert) SetUserID(v uuid.UUID) *AppUserControlUpsert {
	u.Set(appusercontrol.FieldUserID, v)
	return u
}

// UpdateUserID sets the "user_id" field to the value that was provided on create.
func (u *AppUserControlUpsert) UpdateUserID() *AppUserControlUpsert {
	u.SetExcluded(appusercontrol.FieldUserID)
	return u
}

// SetSigninVerifyByGoogleAuthentication sets the "signin_verify_by_google_authentication" field.
func (u *AppUserControlUpsert) SetSigninVerifyByGoogleAuthentication(v bool) *AppUserControlUpsert {
	u.Set(appusercontrol.FieldSigninVerifyByGoogleAuthentication, v)
	return u
}

// UpdateSigninVerifyByGoogleAuthentication sets the "signin_verify_by_google_authentication" field to the value that was provided on create.
func (u *AppUserControlUpsert) UpdateSigninVerifyByGoogleAuthentication() *AppUserControlUpsert {
	u.SetExcluded(appusercontrol.FieldSigninVerifyByGoogleAuthentication)
	return u
}

// SetGoogleAuthenticationVerified sets the "google_authentication_verified" field.
func (u *AppUserControlUpsert) SetGoogleAuthenticationVerified(v bool) *AppUserControlUpsert {
	u.Set(appusercontrol.FieldGoogleAuthenticationVerified, v)
	return u
}

// UpdateGoogleAuthenticationVerified sets the "google_authentication_verified" field to the value that was provided on create.
func (u *AppUserControlUpsert) UpdateGoogleAuthenticationVerified() *AppUserControlUpsert {
	u.SetExcluded(appusercontrol.FieldGoogleAuthenticationVerified)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.AppUserControl.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(appusercontrol.FieldID)
//			}),
//		).
//		Exec(ctx)
//
func (u *AppUserControlUpsertOne) UpdateNewValues() *AppUserControlUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(appusercontrol.FieldID)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//  client.AppUserControl.Create().
//      OnConflict(sql.ResolveWithIgnore()).
//      Exec(ctx)
//
func (u *AppUserControlUpsertOne) Ignore() *AppUserControlUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *AppUserControlUpsertOne) DoNothing() *AppUserControlUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the AppUserControlCreate.OnConflict
// documentation for more info.
func (u *AppUserControlUpsertOne) Update(set func(*AppUserControlUpsert)) *AppUserControlUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&AppUserControlUpsert{UpdateSet: update})
	}))
	return u
}

// SetCreatedAt sets the "created_at" field.
func (u *AppUserControlUpsertOne) SetCreatedAt(v uint32) *AppUserControlUpsertOne {
	return u.Update(func(s *AppUserControlUpsert) {
		s.SetCreatedAt(v)
	})
}

// AddCreatedAt adds v to the "created_at" field.
func (u *AppUserControlUpsertOne) AddCreatedAt(v uint32) *AppUserControlUpsertOne {
	return u.Update(func(s *AppUserControlUpsert) {
		s.AddCreatedAt(v)
	})
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *AppUserControlUpsertOne) UpdateCreatedAt() *AppUserControlUpsertOne {
	return u.Update(func(s *AppUserControlUpsert) {
		s.UpdateCreatedAt()
	})
}

// SetUpdatedAt sets the "updated_at" field.
func (u *AppUserControlUpsertOne) SetUpdatedAt(v uint32) *AppUserControlUpsertOne {
	return u.Update(func(s *AppUserControlUpsert) {
		s.SetUpdatedAt(v)
	})
}

// AddUpdatedAt adds v to the "updated_at" field.
func (u *AppUserControlUpsertOne) AddUpdatedAt(v uint32) *AppUserControlUpsertOne {
	return u.Update(func(s *AppUserControlUpsert) {
		s.AddUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *AppUserControlUpsertOne) UpdateUpdatedAt() *AppUserControlUpsertOne {
	return u.Update(func(s *AppUserControlUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *AppUserControlUpsertOne) SetDeletedAt(v uint32) *AppUserControlUpsertOne {
	return u.Update(func(s *AppUserControlUpsert) {
		s.SetDeletedAt(v)
	})
}

// AddDeletedAt adds v to the "deleted_at" field.
func (u *AppUserControlUpsertOne) AddDeletedAt(v uint32) *AppUserControlUpsertOne {
	return u.Update(func(s *AppUserControlUpsert) {
		s.AddDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *AppUserControlUpsertOne) UpdateDeletedAt() *AppUserControlUpsertOne {
	return u.Update(func(s *AppUserControlUpsert) {
		s.UpdateDeletedAt()
	})
}

// SetAppID sets the "app_id" field.
func (u *AppUserControlUpsertOne) SetAppID(v uuid.UUID) *AppUserControlUpsertOne {
	return u.Update(func(s *AppUserControlUpsert) {
		s.SetAppID(v)
	})
}

// UpdateAppID sets the "app_id" field to the value that was provided on create.
func (u *AppUserControlUpsertOne) UpdateAppID() *AppUserControlUpsertOne {
	return u.Update(func(s *AppUserControlUpsert) {
		s.UpdateAppID()
	})
}

// SetUserID sets the "user_id" field.
func (u *AppUserControlUpsertOne) SetUserID(v uuid.UUID) *AppUserControlUpsertOne {
	return u.Update(func(s *AppUserControlUpsert) {
		s.SetUserID(v)
	})
}

// UpdateUserID sets the "user_id" field to the value that was provided on create.
func (u *AppUserControlUpsertOne) UpdateUserID() *AppUserControlUpsertOne {
	return u.Update(func(s *AppUserControlUpsert) {
		s.UpdateUserID()
	})
}

// SetSigninVerifyByGoogleAuthentication sets the "signin_verify_by_google_authentication" field.
func (u *AppUserControlUpsertOne) SetSigninVerifyByGoogleAuthentication(v bool) *AppUserControlUpsertOne {
	return u.Update(func(s *AppUserControlUpsert) {
		s.SetSigninVerifyByGoogleAuthentication(v)
	})
}

// UpdateSigninVerifyByGoogleAuthentication sets the "signin_verify_by_google_authentication" field to the value that was provided on create.
func (u *AppUserControlUpsertOne) UpdateSigninVerifyByGoogleAuthentication() *AppUserControlUpsertOne {
	return u.Update(func(s *AppUserControlUpsert) {
		s.UpdateSigninVerifyByGoogleAuthentication()
	})
}

// SetGoogleAuthenticationVerified sets the "google_authentication_verified" field.
func (u *AppUserControlUpsertOne) SetGoogleAuthenticationVerified(v bool) *AppUserControlUpsertOne {
	return u.Update(func(s *AppUserControlUpsert) {
		s.SetGoogleAuthenticationVerified(v)
	})
}

// UpdateGoogleAuthenticationVerified sets the "google_authentication_verified" field to the value that was provided on create.
func (u *AppUserControlUpsertOne) UpdateGoogleAuthenticationVerified() *AppUserControlUpsertOne {
	return u.Update(func(s *AppUserControlUpsert) {
		s.UpdateGoogleAuthenticationVerified()
	})
}

// Exec executes the query.
func (u *AppUserControlUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for AppUserControlCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *AppUserControlUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *AppUserControlUpsertOne) ID(ctx context.Context) (id uuid.UUID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("ent: AppUserControlUpsertOne.ID is not supported by MySQL driver. Use AppUserControlUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *AppUserControlUpsertOne) IDX(ctx context.Context) uuid.UUID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// AppUserControlCreateBulk is the builder for creating many AppUserControl entities in bulk.
type AppUserControlCreateBulk struct {
	config
	builders []*AppUserControlCreate
	conflict []sql.ConflictOption
}

// Save creates the AppUserControl entities in the database.
func (auccb *AppUserControlCreateBulk) Save(ctx context.Context) ([]*AppUserControl, error) {
	specs := make([]*sqlgraph.CreateSpec, len(auccb.builders))
	nodes := make([]*AppUserControl, len(auccb.builders))
	mutators := make([]Mutator, len(auccb.builders))
	for i := range auccb.builders {
		func(i int, root context.Context) {
			builder := auccb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*AppUserControlMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, auccb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = auccb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, auccb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{err.Error(), err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, auccb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (auccb *AppUserControlCreateBulk) SaveX(ctx context.Context) []*AppUserControl {
	v, err := auccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (auccb *AppUserControlCreateBulk) Exec(ctx context.Context) error {
	_, err := auccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (auccb *AppUserControlCreateBulk) ExecX(ctx context.Context) {
	if err := auccb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.AppUserControl.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.AppUserControlUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
//
func (auccb *AppUserControlCreateBulk) OnConflict(opts ...sql.ConflictOption) *AppUserControlUpsertBulk {
	auccb.conflict = opts
	return &AppUserControlUpsertBulk{
		create: auccb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.AppUserControl.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
//
func (auccb *AppUserControlCreateBulk) OnConflictColumns(columns ...string) *AppUserControlUpsertBulk {
	auccb.conflict = append(auccb.conflict, sql.ConflictColumns(columns...))
	return &AppUserControlUpsertBulk{
		create: auccb,
	}
}

// AppUserControlUpsertBulk is the builder for "upsert"-ing
// a bulk of AppUserControl nodes.
type AppUserControlUpsertBulk struct {
	create *AppUserControlCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.AppUserControl.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(appusercontrol.FieldID)
//			}),
//		).
//		Exec(ctx)
//
func (u *AppUserControlUpsertBulk) UpdateNewValues() *AppUserControlUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(appusercontrol.FieldID)
				return
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.AppUserControl.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
//
func (u *AppUserControlUpsertBulk) Ignore() *AppUserControlUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *AppUserControlUpsertBulk) DoNothing() *AppUserControlUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the AppUserControlCreateBulk.OnConflict
// documentation for more info.
func (u *AppUserControlUpsertBulk) Update(set func(*AppUserControlUpsert)) *AppUserControlUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&AppUserControlUpsert{UpdateSet: update})
	}))
	return u
}

// SetCreatedAt sets the "created_at" field.
func (u *AppUserControlUpsertBulk) SetCreatedAt(v uint32) *AppUserControlUpsertBulk {
	return u.Update(func(s *AppUserControlUpsert) {
		s.SetCreatedAt(v)
	})
}

// AddCreatedAt adds v to the "created_at" field.
func (u *AppUserControlUpsertBulk) AddCreatedAt(v uint32) *AppUserControlUpsertBulk {
	return u.Update(func(s *AppUserControlUpsert) {
		s.AddCreatedAt(v)
	})
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *AppUserControlUpsertBulk) UpdateCreatedAt() *AppUserControlUpsertBulk {
	return u.Update(func(s *AppUserControlUpsert) {
		s.UpdateCreatedAt()
	})
}

// SetUpdatedAt sets the "updated_at" field.
func (u *AppUserControlUpsertBulk) SetUpdatedAt(v uint32) *AppUserControlUpsertBulk {
	return u.Update(func(s *AppUserControlUpsert) {
		s.SetUpdatedAt(v)
	})
}

// AddUpdatedAt adds v to the "updated_at" field.
func (u *AppUserControlUpsertBulk) AddUpdatedAt(v uint32) *AppUserControlUpsertBulk {
	return u.Update(func(s *AppUserControlUpsert) {
		s.AddUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *AppUserControlUpsertBulk) UpdateUpdatedAt() *AppUserControlUpsertBulk {
	return u.Update(func(s *AppUserControlUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeletedAt sets the "deleted_at" field.
func (u *AppUserControlUpsertBulk) SetDeletedAt(v uint32) *AppUserControlUpsertBulk {
	return u.Update(func(s *AppUserControlUpsert) {
		s.SetDeletedAt(v)
	})
}

// AddDeletedAt adds v to the "deleted_at" field.
func (u *AppUserControlUpsertBulk) AddDeletedAt(v uint32) *AppUserControlUpsertBulk {
	return u.Update(func(s *AppUserControlUpsert) {
		s.AddDeletedAt(v)
	})
}

// UpdateDeletedAt sets the "deleted_at" field to the value that was provided on create.
func (u *AppUserControlUpsertBulk) UpdateDeletedAt() *AppUserControlUpsertBulk {
	return u.Update(func(s *AppUserControlUpsert) {
		s.UpdateDeletedAt()
	})
}

// SetAppID sets the "app_id" field.
func (u *AppUserControlUpsertBulk) SetAppID(v uuid.UUID) *AppUserControlUpsertBulk {
	return u.Update(func(s *AppUserControlUpsert) {
		s.SetAppID(v)
	})
}

// UpdateAppID sets the "app_id" field to the value that was provided on create.
func (u *AppUserControlUpsertBulk) UpdateAppID() *AppUserControlUpsertBulk {
	return u.Update(func(s *AppUserControlUpsert) {
		s.UpdateAppID()
	})
}

// SetUserID sets the "user_id" field.
func (u *AppUserControlUpsertBulk) SetUserID(v uuid.UUID) *AppUserControlUpsertBulk {
	return u.Update(func(s *AppUserControlUpsert) {
		s.SetUserID(v)
	})
}

// UpdateUserID sets the "user_id" field to the value that was provided on create.
func (u *AppUserControlUpsertBulk) UpdateUserID() *AppUserControlUpsertBulk {
	return u.Update(func(s *AppUserControlUpsert) {
		s.UpdateUserID()
	})
}

// SetSigninVerifyByGoogleAuthentication sets the "signin_verify_by_google_authentication" field.
func (u *AppUserControlUpsertBulk) SetSigninVerifyByGoogleAuthentication(v bool) *AppUserControlUpsertBulk {
	return u.Update(func(s *AppUserControlUpsert) {
		s.SetSigninVerifyByGoogleAuthentication(v)
	})
}

// UpdateSigninVerifyByGoogleAuthentication sets the "signin_verify_by_google_authentication" field to the value that was provided on create.
func (u *AppUserControlUpsertBulk) UpdateSigninVerifyByGoogleAuthentication() *AppUserControlUpsertBulk {
	return u.Update(func(s *AppUserControlUpsert) {
		s.UpdateSigninVerifyByGoogleAuthentication()
	})
}

// SetGoogleAuthenticationVerified sets the "google_authentication_verified" field.
func (u *AppUserControlUpsertBulk) SetGoogleAuthenticationVerified(v bool) *AppUserControlUpsertBulk {
	return u.Update(func(s *AppUserControlUpsert) {
		s.SetGoogleAuthenticationVerified(v)
	})
}

// UpdateGoogleAuthenticationVerified sets the "google_authentication_verified" field to the value that was provided on create.
func (u *AppUserControlUpsertBulk) UpdateGoogleAuthenticationVerified() *AppUserControlUpsertBulk {
	return u.Update(func(s *AppUserControlUpsert) {
		s.UpdateGoogleAuthenticationVerified()
	})
}

// Exec executes the query.
func (u *AppUserControlUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the AppUserControlCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for AppUserControlCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *AppUserControlUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
