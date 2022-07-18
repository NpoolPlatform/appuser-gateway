// Code generated by entc, DO NOT EDIT.

package appuser

import (
	"entgo.io/ent/dialect/sql"
	"github.com/NpoolPlatform/appuser-gateway/pkg/db/ent/predicate"
	"github.com/google/uuid"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v uint32) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v uint32) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdatedAt), v))
	})
}

// DeletedAt applies equality check predicate on the "deleted_at" field. It's identical to DeletedAtEQ.
func DeletedAt(v uint32) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldDeletedAt), v))
	})
}

// AppID applies equality check predicate on the "app_id" field. It's identical to AppIDEQ.
func AppID(v uuid.UUID) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldAppID), v))
	})
}

// EmailAddress applies equality check predicate on the "email_address" field. It's identical to EmailAddressEQ.
func EmailAddress(v string) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldEmailAddress), v))
	})
}

// PhoneNo applies equality check predicate on the "phone_no" field. It's identical to PhoneNoEQ.
func PhoneNo(v string) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPhoneNo), v))
	})
}

// ImportFromApp applies equality check predicate on the "import_from_app" field. It's identical to ImportFromAppEQ.
func ImportFromApp(v uuid.UUID) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldImportFromApp), v))
	})
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v uint32) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v uint32) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...uint32) predicate.AppUser {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AppUser(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCreatedAt), v...))
	})
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...uint32) predicate.AppUser {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AppUser(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCreatedAt), v...))
	})
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v uint32) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v uint32) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v uint32) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v uint32) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreatedAt), v))
	})
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v uint32) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v uint32) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...uint32) predicate.AppUser {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AppUser(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldUpdatedAt), v...))
	})
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...uint32) predicate.AppUser {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AppUser(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldUpdatedAt), v...))
	})
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v uint32) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v uint32) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v uint32) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v uint32) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdatedAt), v))
	})
}

// DeletedAtEQ applies the EQ predicate on the "deleted_at" field.
func DeletedAtEQ(v uint32) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldDeletedAt), v))
	})
}

// DeletedAtNEQ applies the NEQ predicate on the "deleted_at" field.
func DeletedAtNEQ(v uint32) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldDeletedAt), v))
	})
}

// DeletedAtIn applies the In predicate on the "deleted_at" field.
func DeletedAtIn(vs ...uint32) predicate.AppUser {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AppUser(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldDeletedAt), v...))
	})
}

// DeletedAtNotIn applies the NotIn predicate on the "deleted_at" field.
func DeletedAtNotIn(vs ...uint32) predicate.AppUser {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AppUser(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldDeletedAt), v...))
	})
}

// DeletedAtGT applies the GT predicate on the "deleted_at" field.
func DeletedAtGT(v uint32) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldDeletedAt), v))
	})
}

// DeletedAtGTE applies the GTE predicate on the "deleted_at" field.
func DeletedAtGTE(v uint32) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldDeletedAt), v))
	})
}

// DeletedAtLT applies the LT predicate on the "deleted_at" field.
func DeletedAtLT(v uint32) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldDeletedAt), v))
	})
}

// DeletedAtLTE applies the LTE predicate on the "deleted_at" field.
func DeletedAtLTE(v uint32) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldDeletedAt), v))
	})
}

// AppIDEQ applies the EQ predicate on the "app_id" field.
func AppIDEQ(v uuid.UUID) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldAppID), v))
	})
}

// AppIDNEQ applies the NEQ predicate on the "app_id" field.
func AppIDNEQ(v uuid.UUID) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldAppID), v))
	})
}

// AppIDIn applies the In predicate on the "app_id" field.
func AppIDIn(vs ...uuid.UUID) predicate.AppUser {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AppUser(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldAppID), v...))
	})
}

// AppIDNotIn applies the NotIn predicate on the "app_id" field.
func AppIDNotIn(vs ...uuid.UUID) predicate.AppUser {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AppUser(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldAppID), v...))
	})
}

// AppIDGT applies the GT predicate on the "app_id" field.
func AppIDGT(v uuid.UUID) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldAppID), v))
	})
}

// AppIDGTE applies the GTE predicate on the "app_id" field.
func AppIDGTE(v uuid.UUID) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldAppID), v))
	})
}

// AppIDLT applies the LT predicate on the "app_id" field.
func AppIDLT(v uuid.UUID) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldAppID), v))
	})
}

// AppIDLTE applies the LTE predicate on the "app_id" field.
func AppIDLTE(v uuid.UUID) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldAppID), v))
	})
}

// EmailAddressEQ applies the EQ predicate on the "email_address" field.
func EmailAddressEQ(v string) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldEmailAddress), v))
	})
}

// EmailAddressNEQ applies the NEQ predicate on the "email_address" field.
func EmailAddressNEQ(v string) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldEmailAddress), v))
	})
}

// EmailAddressIn applies the In predicate on the "email_address" field.
func EmailAddressIn(vs ...string) predicate.AppUser {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AppUser(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldEmailAddress), v...))
	})
}

// EmailAddressNotIn applies the NotIn predicate on the "email_address" field.
func EmailAddressNotIn(vs ...string) predicate.AppUser {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AppUser(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldEmailAddress), v...))
	})
}

// EmailAddressGT applies the GT predicate on the "email_address" field.
func EmailAddressGT(v string) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldEmailAddress), v))
	})
}

// EmailAddressGTE applies the GTE predicate on the "email_address" field.
func EmailAddressGTE(v string) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldEmailAddress), v))
	})
}

// EmailAddressLT applies the LT predicate on the "email_address" field.
func EmailAddressLT(v string) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldEmailAddress), v))
	})
}

// EmailAddressLTE applies the LTE predicate on the "email_address" field.
func EmailAddressLTE(v string) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldEmailAddress), v))
	})
}

// EmailAddressContains applies the Contains predicate on the "email_address" field.
func EmailAddressContains(v string) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldEmailAddress), v))
	})
}

// EmailAddressHasPrefix applies the HasPrefix predicate on the "email_address" field.
func EmailAddressHasPrefix(v string) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldEmailAddress), v))
	})
}

// EmailAddressHasSuffix applies the HasSuffix predicate on the "email_address" field.
func EmailAddressHasSuffix(v string) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldEmailAddress), v))
	})
}

// EmailAddressEqualFold applies the EqualFold predicate on the "email_address" field.
func EmailAddressEqualFold(v string) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldEmailAddress), v))
	})
}

// EmailAddressContainsFold applies the ContainsFold predicate on the "email_address" field.
func EmailAddressContainsFold(v string) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldEmailAddress), v))
	})
}

// PhoneNoEQ applies the EQ predicate on the "phone_no" field.
func PhoneNoEQ(v string) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPhoneNo), v))
	})
}

// PhoneNoNEQ applies the NEQ predicate on the "phone_no" field.
func PhoneNoNEQ(v string) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldPhoneNo), v))
	})
}

// PhoneNoIn applies the In predicate on the "phone_no" field.
func PhoneNoIn(vs ...string) predicate.AppUser {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AppUser(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldPhoneNo), v...))
	})
}

// PhoneNoNotIn applies the NotIn predicate on the "phone_no" field.
func PhoneNoNotIn(vs ...string) predicate.AppUser {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AppUser(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldPhoneNo), v...))
	})
}

// PhoneNoGT applies the GT predicate on the "phone_no" field.
func PhoneNoGT(v string) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldPhoneNo), v))
	})
}

// PhoneNoGTE applies the GTE predicate on the "phone_no" field.
func PhoneNoGTE(v string) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldPhoneNo), v))
	})
}

// PhoneNoLT applies the LT predicate on the "phone_no" field.
func PhoneNoLT(v string) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldPhoneNo), v))
	})
}

// PhoneNoLTE applies the LTE predicate on the "phone_no" field.
func PhoneNoLTE(v string) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldPhoneNo), v))
	})
}

// PhoneNoContains applies the Contains predicate on the "phone_no" field.
func PhoneNoContains(v string) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldPhoneNo), v))
	})
}

// PhoneNoHasPrefix applies the HasPrefix predicate on the "phone_no" field.
func PhoneNoHasPrefix(v string) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldPhoneNo), v))
	})
}

// PhoneNoHasSuffix applies the HasSuffix predicate on the "phone_no" field.
func PhoneNoHasSuffix(v string) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldPhoneNo), v))
	})
}

// PhoneNoEqualFold applies the EqualFold predicate on the "phone_no" field.
func PhoneNoEqualFold(v string) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldPhoneNo), v))
	})
}

// PhoneNoContainsFold applies the ContainsFold predicate on the "phone_no" field.
func PhoneNoContainsFold(v string) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldPhoneNo), v))
	})
}

// ImportFromAppEQ applies the EQ predicate on the "import_from_app" field.
func ImportFromAppEQ(v uuid.UUID) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldImportFromApp), v))
	})
}

// ImportFromAppNEQ applies the NEQ predicate on the "import_from_app" field.
func ImportFromAppNEQ(v uuid.UUID) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldImportFromApp), v))
	})
}

// ImportFromAppIn applies the In predicate on the "import_from_app" field.
func ImportFromAppIn(vs ...uuid.UUID) predicate.AppUser {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AppUser(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldImportFromApp), v...))
	})
}

// ImportFromAppNotIn applies the NotIn predicate on the "import_from_app" field.
func ImportFromAppNotIn(vs ...uuid.UUID) predicate.AppUser {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.AppUser(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldImportFromApp), v...))
	})
}

// ImportFromAppGT applies the GT predicate on the "import_from_app" field.
func ImportFromAppGT(v uuid.UUID) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldImportFromApp), v))
	})
}

// ImportFromAppGTE applies the GTE predicate on the "import_from_app" field.
func ImportFromAppGTE(v uuid.UUID) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldImportFromApp), v))
	})
}

// ImportFromAppLT applies the LT predicate on the "import_from_app" field.
func ImportFromAppLT(v uuid.UUID) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldImportFromApp), v))
	})
}

// ImportFromAppLTE applies the LTE predicate on the "import_from_app" field.
func ImportFromAppLTE(v uuid.UUID) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldImportFromApp), v))
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.AppUser) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.AppUser) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.AppUser) predicate.AppUser {
	return predicate.AppUser(func(s *sql.Selector) {
		p(s.Not())
	})
}
