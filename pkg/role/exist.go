package role

import (
	"context"
	"fmt"

	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	rolemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) ExistRole(ctx context.Context) error {
	exist, err := rolemwcli.ExistRoleConds(ctx, &rolemwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
	})
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid role")
	}
	return nil
}
