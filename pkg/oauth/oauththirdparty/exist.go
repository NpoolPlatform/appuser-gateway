package oauththirdparty

import (
	"context"
	"fmt"

	oauththirdpartymwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/oauth/oauththirdparty"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	oauththirdpartymwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/oauth/oauththirdparty"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) ExistOAuthThirdParty(ctx context.Context) error {
	exist, err := oauththirdpartymwcli.ExistOAuthThirdPartyConds(ctx, &oauththirdpartymwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
	})
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid thirdparty")
	}
	return nil
}
