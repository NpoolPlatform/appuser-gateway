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
	return usermwcli.DeleteUser(ctx, *h.AppID, *h.ID)
}
