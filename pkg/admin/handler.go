package admin

import (
	"context"
	"fmt"
	"net/mail"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
	rolemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role"
	roleusermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role/user"

	"github.com/google/uuid"
)

type Handler struct {
	AppID        string
	EmailAddress *string
	PasswordHash *string
	Offset       int32
	Limit        int32
	GenesisApps  []*appmwpb.App
	GenesisRoles []*rolemwpb.Role
	GenesisUsers []*roleusermwpb.User
}

func NewHandler(ctx context.Context, options ...func(context.Context, *Handler) error) (*Handler, error) {
	handler := &Handler{}
	for _, opt := range options {
		if err := opt(ctx, handler); err != nil {
			return nil, err
		}
	}
	return handler, nil
}

func WithAppID(id string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if _, err := uuid.Parse(id); err != nil {
			return err
		}
		h.AppID = id
		return nil
	}
}

func WithEmailAddress(emailAddress *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if emailAddress == nil {
			return nil
		}
		if _, err := mail.ParseAddress(*emailAddress); err != nil {
			return err
		}
		h.EmailAddress = emailAddress
		return nil
	}
}

func WithPasswordHash(pwdHash *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if pwdHash == nil {
			return nil
		}
		if *pwdHash == "" {
			return fmt.Errorf("invalid password")
		}
		h.PasswordHash = pwdHash
		return nil
	}
}

func WithOffset(offset int32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Offset = offset
		return nil
	}
}

func WithLimit(limit int32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if limit <= 0 {
			limit = constant.DefaultRowLimit
		}
		h.Limit = limit
		return nil
	}
}
