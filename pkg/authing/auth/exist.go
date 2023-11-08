package auth

import (
	"context"
	"fmt"

	authmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/authing/auth"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	authmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/authing/auth"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) ExistAuth(ctx context.Context) error {
	exist, err := authmwcli.ExistAuthConds(ctx, &authmwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
	})
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid auth")
	}
	return nil
}
