package user

import (
	"context"
	"fmt"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/google/uuid"
)

type Handler struct {
	AppID            string
	App              *appmwpb.App
	Account          string
	PasswordHash     string
	AccountType      basetypes.SignMethod
	VerificationCode string
	InvitationCode   *string
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
		app, err := appmwcli.GetApp(ctx, appID)
		if err != nil {
			return err
		}
		if app == nil {
			return fmt.Errorf("invalid app")
		}
		h.AppID = appID
		h.App = app
		return nil
	}
}

func WithPasswordHash(pwdHash string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if pwdHash == "" {
			return fmt.Errorf("invalid password")
		}
		h.PasswordHash = pwdHash
		return nil
	}
}

func WithAccount(account string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if account == "" {
			return fmt.Errorf("invalid account")
		}
		h.Account = account
		return nil
	}
}

func WithAccountType(accountType basetypes.SignMethod) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		switch accountType {
		case basetypes.SignMethod_Mobile:
		case basetypes.SignMethod_Email:
		default:
			return fmt.Errorf("invalid account type")
		}
		h.AccountType = accountType
		return nil
	}
}

func WithVerificationCode(code string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if code == "" {
			return fmt.Errorf("invalid verification code")
		}
		h.VerificationCode = code
		return nil
	}
}

func WithInvitationCode(code *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if code == nil {
			return nil
		}
		if *code == "" {
			return fmt.Errorf("invalid invitation code")
		}
		h.InvitationCode = code
		return nil
	}
}
