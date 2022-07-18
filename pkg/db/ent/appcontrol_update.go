// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/NpoolPlatform/appuser-gateway/pkg/db/ent/appcontrol"
	"github.com/NpoolPlatform/appuser-gateway/pkg/db/ent/predicate"
	"github.com/google/uuid"
)

// AppControlUpdate is the builder for updating AppControl entities.
type AppControlUpdate struct {
	config
	hooks    []Hook
	mutation *AppControlMutation
}

// Where appends a list predicates to the AppControlUpdate builder.
func (acu *AppControlUpdate) Where(ps ...predicate.AppControl) *AppControlUpdate {
	acu.mutation.Where(ps...)
	return acu
}

// SetCreatedAt sets the "created_at" field.
func (acu *AppControlUpdate) SetCreatedAt(u uint32) *AppControlUpdate {
	acu.mutation.ResetCreatedAt()
	acu.mutation.SetCreatedAt(u)
	return acu
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (acu *AppControlUpdate) SetNillableCreatedAt(u *uint32) *AppControlUpdate {
	if u != nil {
		acu.SetCreatedAt(*u)
	}
	return acu
}

// AddCreatedAt adds u to the "created_at" field.
func (acu *AppControlUpdate) AddCreatedAt(u int32) *AppControlUpdate {
	acu.mutation.AddCreatedAt(u)
	return acu
}

// SetUpdatedAt sets the "updated_at" field.
func (acu *AppControlUpdate) SetUpdatedAt(u uint32) *AppControlUpdate {
	acu.mutation.ResetUpdatedAt()
	acu.mutation.SetUpdatedAt(u)
	return acu
}

// AddUpdatedAt adds u to the "updated_at" field.
func (acu *AppControlUpdate) AddUpdatedAt(u int32) *AppControlUpdate {
	acu.mutation.AddUpdatedAt(u)
	return acu
}

// SetDeletedAt sets the "deleted_at" field.
func (acu *AppControlUpdate) SetDeletedAt(u uint32) *AppControlUpdate {
	acu.mutation.ResetDeletedAt()
	acu.mutation.SetDeletedAt(u)
	return acu
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (acu *AppControlUpdate) SetNillableDeletedAt(u *uint32) *AppControlUpdate {
	if u != nil {
		acu.SetDeletedAt(*u)
	}
	return acu
}

// AddDeletedAt adds u to the "deleted_at" field.
func (acu *AppControlUpdate) AddDeletedAt(u int32) *AppControlUpdate {
	acu.mutation.AddDeletedAt(u)
	return acu
}

// SetAppID sets the "app_id" field.
func (acu *AppControlUpdate) SetAppID(u uuid.UUID) *AppControlUpdate {
	acu.mutation.SetAppID(u)
	return acu
}

// SetSignupMethods sets the "signup_methods" field.
func (acu *AppControlUpdate) SetSignupMethods(s []string) *AppControlUpdate {
	acu.mutation.SetSignupMethods(s)
	return acu
}

// SetExternSigninMethods sets the "extern_signin_methods" field.
func (acu *AppControlUpdate) SetExternSigninMethods(s []string) *AppControlUpdate {
	acu.mutation.SetExternSigninMethods(s)
	return acu
}

// SetRecaptchaMethod sets the "recaptcha_method" field.
func (acu *AppControlUpdate) SetRecaptchaMethod(s string) *AppControlUpdate {
	acu.mutation.SetRecaptchaMethod(s)
	return acu
}

// SetKycEnable sets the "kyc_enable" field.
func (acu *AppControlUpdate) SetKycEnable(b bool) *AppControlUpdate {
	acu.mutation.SetKycEnable(b)
	return acu
}

// SetSigninVerifyEnable sets the "signin_verify_enable" field.
func (acu *AppControlUpdate) SetSigninVerifyEnable(b bool) *AppControlUpdate {
	acu.mutation.SetSigninVerifyEnable(b)
	return acu
}

// SetInvitationCodeMust sets the "invitation_code_must" field.
func (acu *AppControlUpdate) SetInvitationCodeMust(b bool) *AppControlUpdate {
	acu.mutation.SetInvitationCodeMust(b)
	return acu
}

// Mutation returns the AppControlMutation object of the builder.
func (acu *AppControlUpdate) Mutation() *AppControlMutation {
	return acu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (acu *AppControlUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if err := acu.defaults(); err != nil {
		return 0, err
	}
	if len(acu.hooks) == 0 {
		affected, err = acu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*AppControlMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			acu.mutation = mutation
			affected, err = acu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(acu.hooks) - 1; i >= 0; i-- {
			if acu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = acu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, acu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (acu *AppControlUpdate) SaveX(ctx context.Context) int {
	affected, err := acu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (acu *AppControlUpdate) Exec(ctx context.Context) error {
	_, err := acu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (acu *AppControlUpdate) ExecX(ctx context.Context) {
	if err := acu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (acu *AppControlUpdate) defaults() error {
	if _, ok := acu.mutation.UpdatedAt(); !ok {
		if appcontrol.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized appcontrol.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := appcontrol.UpdateDefaultUpdatedAt()
		acu.mutation.SetUpdatedAt(v)
	}
	return nil
}

func (acu *AppControlUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   appcontrol.Table,
			Columns: appcontrol.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: appcontrol.FieldID,
			},
		},
	}
	if ps := acu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := acu.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: appcontrol.FieldCreatedAt,
		})
	}
	if value, ok := acu.mutation.AddedCreatedAt(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: appcontrol.FieldCreatedAt,
		})
	}
	if value, ok := acu.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: appcontrol.FieldUpdatedAt,
		})
	}
	if value, ok := acu.mutation.AddedUpdatedAt(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: appcontrol.FieldUpdatedAt,
		})
	}
	if value, ok := acu.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: appcontrol.FieldDeletedAt,
		})
	}
	if value, ok := acu.mutation.AddedDeletedAt(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: appcontrol.FieldDeletedAt,
		})
	}
	if value, ok := acu.mutation.AppID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUUID,
			Value:  value,
			Column: appcontrol.FieldAppID,
		})
	}
	if value, ok := acu.mutation.SignupMethods(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: appcontrol.FieldSignupMethods,
		})
	}
	if value, ok := acu.mutation.ExternSigninMethods(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: appcontrol.FieldExternSigninMethods,
		})
	}
	if value, ok := acu.mutation.RecaptchaMethod(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: appcontrol.FieldRecaptchaMethod,
		})
	}
	if value, ok := acu.mutation.KycEnable(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: appcontrol.FieldKycEnable,
		})
	}
	if value, ok := acu.mutation.SigninVerifyEnable(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: appcontrol.FieldSigninVerifyEnable,
		})
	}
	if value, ok := acu.mutation.InvitationCodeMust(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: appcontrol.FieldInvitationCodeMust,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, acu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{appcontrol.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// AppControlUpdateOne is the builder for updating a single AppControl entity.
type AppControlUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *AppControlMutation
}

// SetCreatedAt sets the "created_at" field.
func (acuo *AppControlUpdateOne) SetCreatedAt(u uint32) *AppControlUpdateOne {
	acuo.mutation.ResetCreatedAt()
	acuo.mutation.SetCreatedAt(u)
	return acuo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (acuo *AppControlUpdateOne) SetNillableCreatedAt(u *uint32) *AppControlUpdateOne {
	if u != nil {
		acuo.SetCreatedAt(*u)
	}
	return acuo
}

// AddCreatedAt adds u to the "created_at" field.
func (acuo *AppControlUpdateOne) AddCreatedAt(u int32) *AppControlUpdateOne {
	acuo.mutation.AddCreatedAt(u)
	return acuo
}

// SetUpdatedAt sets the "updated_at" field.
func (acuo *AppControlUpdateOne) SetUpdatedAt(u uint32) *AppControlUpdateOne {
	acuo.mutation.ResetUpdatedAt()
	acuo.mutation.SetUpdatedAt(u)
	return acuo
}

// AddUpdatedAt adds u to the "updated_at" field.
func (acuo *AppControlUpdateOne) AddUpdatedAt(u int32) *AppControlUpdateOne {
	acuo.mutation.AddUpdatedAt(u)
	return acuo
}

// SetDeletedAt sets the "deleted_at" field.
func (acuo *AppControlUpdateOne) SetDeletedAt(u uint32) *AppControlUpdateOne {
	acuo.mutation.ResetDeletedAt()
	acuo.mutation.SetDeletedAt(u)
	return acuo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (acuo *AppControlUpdateOne) SetNillableDeletedAt(u *uint32) *AppControlUpdateOne {
	if u != nil {
		acuo.SetDeletedAt(*u)
	}
	return acuo
}

// AddDeletedAt adds u to the "deleted_at" field.
func (acuo *AppControlUpdateOne) AddDeletedAt(u int32) *AppControlUpdateOne {
	acuo.mutation.AddDeletedAt(u)
	return acuo
}

// SetAppID sets the "app_id" field.
func (acuo *AppControlUpdateOne) SetAppID(u uuid.UUID) *AppControlUpdateOne {
	acuo.mutation.SetAppID(u)
	return acuo
}

// SetSignupMethods sets the "signup_methods" field.
func (acuo *AppControlUpdateOne) SetSignupMethods(s []string) *AppControlUpdateOne {
	acuo.mutation.SetSignupMethods(s)
	return acuo
}

// SetExternSigninMethods sets the "extern_signin_methods" field.
func (acuo *AppControlUpdateOne) SetExternSigninMethods(s []string) *AppControlUpdateOne {
	acuo.mutation.SetExternSigninMethods(s)
	return acuo
}

// SetRecaptchaMethod sets the "recaptcha_method" field.
func (acuo *AppControlUpdateOne) SetRecaptchaMethod(s string) *AppControlUpdateOne {
	acuo.mutation.SetRecaptchaMethod(s)
	return acuo
}

// SetKycEnable sets the "kyc_enable" field.
func (acuo *AppControlUpdateOne) SetKycEnable(b bool) *AppControlUpdateOne {
	acuo.mutation.SetKycEnable(b)
	return acuo
}

// SetSigninVerifyEnable sets the "signin_verify_enable" field.
func (acuo *AppControlUpdateOne) SetSigninVerifyEnable(b bool) *AppControlUpdateOne {
	acuo.mutation.SetSigninVerifyEnable(b)
	return acuo
}

// SetInvitationCodeMust sets the "invitation_code_must" field.
func (acuo *AppControlUpdateOne) SetInvitationCodeMust(b bool) *AppControlUpdateOne {
	acuo.mutation.SetInvitationCodeMust(b)
	return acuo
}

// Mutation returns the AppControlMutation object of the builder.
func (acuo *AppControlUpdateOne) Mutation() *AppControlMutation {
	return acuo.mutation
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (acuo *AppControlUpdateOne) Select(field string, fields ...string) *AppControlUpdateOne {
	acuo.fields = append([]string{field}, fields...)
	return acuo
}

// Save executes the query and returns the updated AppControl entity.
func (acuo *AppControlUpdateOne) Save(ctx context.Context) (*AppControl, error) {
	var (
		err  error
		node *AppControl
	)
	if err := acuo.defaults(); err != nil {
		return nil, err
	}
	if len(acuo.hooks) == 0 {
		node, err = acuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*AppControlMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			acuo.mutation = mutation
			node, err = acuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(acuo.hooks) - 1; i >= 0; i-- {
			if acuo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = acuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, acuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (acuo *AppControlUpdateOne) SaveX(ctx context.Context) *AppControl {
	node, err := acuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (acuo *AppControlUpdateOne) Exec(ctx context.Context) error {
	_, err := acuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (acuo *AppControlUpdateOne) ExecX(ctx context.Context) {
	if err := acuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (acuo *AppControlUpdateOne) defaults() error {
	if _, ok := acuo.mutation.UpdatedAt(); !ok {
		if appcontrol.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized appcontrol.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := appcontrol.UpdateDefaultUpdatedAt()
		acuo.mutation.SetUpdatedAt(v)
	}
	return nil
}

func (acuo *AppControlUpdateOne) sqlSave(ctx context.Context) (_node *AppControl, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   appcontrol.Table,
			Columns: appcontrol.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: appcontrol.FieldID,
			},
		},
	}
	id, ok := acuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "AppControl.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := acuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, appcontrol.FieldID)
		for _, f := range fields {
			if !appcontrol.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != appcontrol.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := acuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := acuo.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: appcontrol.FieldCreatedAt,
		})
	}
	if value, ok := acuo.mutation.AddedCreatedAt(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: appcontrol.FieldCreatedAt,
		})
	}
	if value, ok := acuo.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: appcontrol.FieldUpdatedAt,
		})
	}
	if value, ok := acuo.mutation.AddedUpdatedAt(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: appcontrol.FieldUpdatedAt,
		})
	}
	if value, ok := acuo.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: appcontrol.FieldDeletedAt,
		})
	}
	if value, ok := acuo.mutation.AddedDeletedAt(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: appcontrol.FieldDeletedAt,
		})
	}
	if value, ok := acuo.mutation.AppID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUUID,
			Value:  value,
			Column: appcontrol.FieldAppID,
		})
	}
	if value, ok := acuo.mutation.SignupMethods(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: appcontrol.FieldSignupMethods,
		})
	}
	if value, ok := acuo.mutation.ExternSigninMethods(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: appcontrol.FieldExternSigninMethods,
		})
	}
	if value, ok := acuo.mutation.RecaptchaMethod(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: appcontrol.FieldRecaptchaMethod,
		})
	}
	if value, ok := acuo.mutation.KycEnable(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: appcontrol.FieldKycEnable,
		})
	}
	if value, ok := acuo.mutation.SigninVerifyEnable(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: appcontrol.FieldSigninVerifyEnable,
		})
	}
	if value, ok := acuo.mutation.InvitationCodeMust(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: appcontrol.FieldInvitationCodeMust,
		})
	}
	_node = &AppControl{config: acuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, acuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{appcontrol.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
