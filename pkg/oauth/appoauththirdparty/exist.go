package appoauththirdparty

import (
	"context"
	"fmt"

	appoauththirdpartymwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/oauth/appoauththirdparty"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	appoauththirdpartymwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/oauth/appoauththirdparty"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) ExistOAuthThirdParty(ctx context.Context) error {
	exist, err := appoauththirdpartymwcli.ExistOAuthThirdPartyConds(ctx, &appoauththirdpartymwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
	})
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid appthirdparty")
	}
	return nil
}
