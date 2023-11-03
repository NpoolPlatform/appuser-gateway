package user

import (
	"context"

	roleusermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role/user"
	roleusermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role/user"
)

func (h *Handler) CreateUser(ctx context.Context) (*roleusermwpb.User, error) {
	return roleusermwcli.CreateUser(ctx, &roleusermwpb.UserReq{
		AppID:  h.AppID,
		RoleID: h.RoleID,
		UserID: h.UserID,
	})
}
