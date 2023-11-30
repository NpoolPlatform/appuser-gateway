package role

import (
	"context"

	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	rolemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) GetRoles(ctx context.Context) ([]*rolemwpb.Role, uint32, error) {
	conds := &rolemwpb.Conds{}
	if h.AppID != nil {
		conds.AppID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID}
	}
	return rolemwcli.GetRoles(ctx, conds, h.Offset, h.Limit)
}
