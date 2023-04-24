package user

import (
	"context"
	"fmt"

	roleusermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role/user"
	roleusermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role/user"
)

func (h *Handler) DeleteUser(ctx context.Context) (*roleusermwpb.User, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}
	return roleusermwcli.DeleteUser(ctx, *h.ID)
}
