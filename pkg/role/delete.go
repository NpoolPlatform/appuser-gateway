package role

import (
	"context"

	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role"
	rolemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role"
)

func (h *Handler) DeleteRole(ctx context.Context) (*rolemwpb.Role, error) {
	if err := h.ExistRole(ctx); err != nil {
		return nil, err
	}
	return rolemwcli.DeleteRole(ctx, *h.ID)
}
