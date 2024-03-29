package user

import (
	"context"

	roleusermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role/user"
	roleusermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role/user"
)

func (h *Handler) DeleteUser(ctx context.Context) (*roleusermwpb.User, error) {
	if err := h.ExistUser(ctx); err != nil {
		return nil, err
	}
	return roleusermwcli.DeleteUser(ctx, *h.ID)
}
