package ga

import (
	"context"
	"fmt"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"

	"github.com/google/uuid"
)

type Handler struct {
	AppID  string
	UserID string
	Code   string
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

func WithUserID(id string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if _, err := uuid.Parse(id); err != nil {
			return err
		}
		h.UserID = id
		return nil
	}
}

func WithCode(code string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.UserID = code
		return nil
	}
}
