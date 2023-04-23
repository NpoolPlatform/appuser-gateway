package auth

import (
	"context"
	"fmt"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role"

	"github.com/google/uuid"
)

type Handler struct {
	ID       *string
	AppID    string
	UserID   *string
	Token    *string
	RoleID   *string
	Resource string
	Method   string
	Offset   int32
	Limit    int32
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

func WithUserID(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		// Check app/user exist at lower layer
		h.UserID = id
		return nil
	}
}

func WithToken(token *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Token = token
		return nil
	}
}

func WithRoleID(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		exist, err := rolemwcli.ExistRole(ctx, *id)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("role not exist")
		}
		h.RoleID = id
		return nil
	}
}

func WithMethod(method string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		switch method {
		case "POST":
		case "GET":
		default:
			return fmt.Errorf("method %v invalid", method)
		}
		h.Method = method
		return nil
	}
}

func WithResource(resource string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		const leastResourceLen = 3
		if len(resource) < leastResourceLen {
			return fmt.Errorf("resource %v invalid", resource)
		}
		h.Resource = resource
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
