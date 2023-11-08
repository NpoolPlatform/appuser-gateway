package user

import (
	"context"
	"fmt"

	roleusermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role/user"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	roleusermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) ExistUser(ctx context.Context) error {
	exist, err := roleusermwcli.ExistUserConds(ctx, &roleusermwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
	})
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid roleuser")
	}
	return nil
}
