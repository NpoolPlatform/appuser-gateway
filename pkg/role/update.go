package role

import (
	"context"
	"fmt"

	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	rolemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) UpdateRole(ctx context.Context) (*rolemwpb.Role, error) {
	info, err := rolemwcli.GetRoleOnly(ctx, &rolemwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
	})
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid role")
	}
	if info.AppID != *h.AppID {
		return nil, fmt.Errorf("permission denied")
	}
	return rolemwcli.UpdateRole(ctx, &rolemwpb.RoleReq{
		ID:          h.ID,
		AppID:       h.AppID,
		CreatedBy:   h.CreatedBy,
		Role:        h.Role,
		Description: h.Description,
		Default:     h.Default,
	})
}
