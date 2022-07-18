// Code generated by entc, DO NOT EDIT.

package appuserthirdparty

import (
	"entgo.io/ent"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the appuserthirdparty type in the database.
	Label = "app_user_third_party"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldDeletedAt holds the string denoting the deleted_at field in the database.
	FieldDeletedAt = "deleted_at"
	// FieldAppID holds the string denoting the app_id field in the database.
	FieldAppID = "app_id"
	// FieldUserID holds the string denoting the user_id field in the database.
	FieldUserID = "user_id"
	// FieldThirdPartyUserID holds the string denoting the third_party_user_id field in the database.
	FieldThirdPartyUserID = "third_party_user_id"
	// FieldThirdPartyID holds the string denoting the third_party_id field in the database.
	FieldThirdPartyID = "third_party_id"
	// FieldThirdPartyUsername holds the string denoting the third_party_username field in the database.
	FieldThirdPartyUsername = "third_party_username"
	// FieldThirdPartyUserAvatar holds the string denoting the third_party_user_avatar field in the database.
	FieldThirdPartyUserAvatar = "third_party_user_avatar"
	// Table holds the table name of the appuserthirdparty in the database.
	Table = "app_user_third_parties"
)

// Columns holds all SQL columns for appuserthirdparty fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldDeletedAt,
	FieldAppID,
	FieldUserID,
	FieldThirdPartyUserID,
	FieldThirdPartyID,
	FieldThirdPartyUsername,
	FieldThirdPartyUserAvatar,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// Note that the variables below are initialized by the runtime
// package on the initialization of the application. Therefore,
// it should be imported in the main as follows:
//
//	import _ "github.com/NpoolPlatform/appuser-gateway/pkg/db/ent/runtime"
//
var (
	Hooks  [1]ent.Hook
	Policy ent.Policy
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() uint32
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() uint32
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() uint32
	// DefaultDeletedAt holds the default value on creation for the "deleted_at" field.
	DefaultDeletedAt func() uint32
	// DefaultThirdPartyUsername holds the default value on creation for the "third_party_username" field.
	DefaultThirdPartyUsername string
	// DefaultThirdPartyUserAvatar holds the default value on creation for the "third_party_user_avatar" field.
	DefaultThirdPartyUserAvatar string
	// ThirdPartyUserAvatarValidator is a validator for the "third_party_user_avatar" field. It is called by the builders before save.
	ThirdPartyUserAvatarValidator func(string) error
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)
