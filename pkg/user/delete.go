package user

import (
	"context"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
)

func (h *Handler) DeleteUser(ctx context.Context) (*usermwpb.User, error) {
	if err := h.ExistUser(ctx); err != nil {
		return nil, err
	}
	info, err := usermwcli.DeleteUser(ctx, *h.AppID, *h.ID)
	if err != nil {
		return nil, err
	}
	meta, err := h.QueryCache(ctx)
	if err != nil {
		return nil, err
	}

	if meta == nil {
		return info, nil
	}
	h.Metadata = meta
	if err := h.DeleteCache(ctx); err != nil {
		return nil, err
	}
	return info, nil
}
