package oauththirdparty

import (
	"context"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/google/uuid"
)

type Handler struct {
	ID             *string
	ClientName     *basetypes.SignMethod
	ClientTag      *string
	ClientOAuthURL *string
	ClientLogoURL  *string
	ResponseType   *string
	Scope          *string
	Offset         int32
	Limit          int32
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

func WithClientName(clientName *basetypes.SignMethod) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.ClientName = clientName
		return nil
	}
}

func WithClientTag(code *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.ClientTag = code
		return nil
	}
}

func WithClientLogoURL(state *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.ClientLogoURL = state
		return nil
	}
}

func WithClientOAuthURL(state *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.ClientOAuthURL = state
		return nil
	}
}

func WithResponseType(state *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.ResponseType = state
		return nil
	}
}

func WithScope(state *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Scope = state
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
