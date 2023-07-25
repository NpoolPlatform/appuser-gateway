package user

import (
	"context"
	"fmt"

	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
)

func (h *Handler) BindUser(ctx context.Context) (*usermwpb.User, error) {
	handler := &updateHandler{
		Handler: h,
	}

	if h.UserID == nil {
		return nil, fmt.Errorf("invalid userid")
	}

	if err := handler.CheckNewAccount(ctx); err != nil {
		return nil, err
	}
	if err := handler.getUser(ctx); err != nil {
		return nil, err
	}
	if err := handler.verifyNewAccountCode(ctx); err != nil {
		return nil, err
	}
	if err := handler.updateUser(ctx); err != nil {
		return nil, err
	}

	if !h.ShouldUpdateCache {
		return h.User, nil
	}

	if err := h.UpdateCache(ctx); err != nil {
		return nil, err
	}
	meta, err := h.QueryCache(ctx)
	if err != nil {
		return nil, err
	}
	h.Metadata = meta

	return h.Metadata.User, nil
}
