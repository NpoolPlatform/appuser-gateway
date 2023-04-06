package user

import (
	"context"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/google/uuid"
)

type Handler struct {
	AppID            string
	Account          string
	PasswordHash     string
	AccountType      basetypes.SignMethod
	VerificationCode string
	InvitationCode   string
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

func WithAppID(appID string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if _, err := uuid.Parse(appID); err != nil {
			return err
		}
		h.AppID = appID
		return nil
	}
}
