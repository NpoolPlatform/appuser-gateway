package role

import (
	"context"

	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role"
	rolemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role"
)

func (h *Handler) UpdateRole(ctx context.Context) (*rolemwpb.Role, error) {
	if err := h.ExistRole(ctx); err != nil {
		return nil, err
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
