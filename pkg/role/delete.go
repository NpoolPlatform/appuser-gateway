package role

import (
	"context"

	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role"
	rolemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role"
)

func (h *Handler) DeleteRole(ctx context.Context) (*rolemwpb.Role, error) {
	return rolemwcli.DeleteRole(ctx, *h.ID)
}
