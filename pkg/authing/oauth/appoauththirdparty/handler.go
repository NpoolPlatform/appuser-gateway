package appoauththirdparty

import (
	"context"
	"fmt"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"

	"github.com/google/uuid"
)

type Handler struct {
	ID           *string
	AppID        string
	ThirdPartyID *string
	ClientID     *string
	ClientSecret *string
	CallbackURL  *string
	Offset       int32
	Limit        int32
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

func WithID(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.ID = id
		return nil
	}
}

func WithAppID(id string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if _, err := uuid.Parse(id); err != nil {
			return err
		}
		exist, err := appmwcli.ExistApp(ctx, id)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("app not exist")
		}
		h.AppID = id
		return nil
	}
}

func WithThirdPartyID(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.ThirdPartyID = id
		return nil
	}
}

func WithClientID(code *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.ClientID = code
		return nil
	}
}

func WithClientSecret(state *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.ClientSecret = state
		return nil
	}
}

func WithCallbackURL(state *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.CallbackURL = state
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
		if limit == 0 {
			limit = constant.DefaultRowLimit
		}
		h.Limit = limit
		return nil
	}
}
