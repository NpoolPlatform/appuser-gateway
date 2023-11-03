package user

import (
	"context"

	roleusermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role/user"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	roleusermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) GetUsers(ctx context.Context) ([]*roleusermwpb.User, uint32, error) {
	conds := &roleusermwpb.Conds{}
	if h.AppID != nil {
		conds.AppID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID}
	}
	if h.RoleID != nil {
		conds.RoleID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.RoleID}
	}
	return roleusermwcli.GetUsers(ctx, conds, h.Offset, h.Limit)
}
