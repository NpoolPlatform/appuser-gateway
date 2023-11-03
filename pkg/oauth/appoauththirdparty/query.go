package appoauththirdparty

import (
	"context"

	oauththirdpartymwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/oauth/appoauththirdparty"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	oauththirdpartymwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/oauth/appoauththirdparty"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) GetOAuthThirdParties(ctx context.Context) ([]*oauththirdpartymwpb.OAuthThirdParty, uint32, error) {
	conds := &oauththirdpartymwpb.Conds{}
	if h.AppID != nil {
		conds.AppID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID}
	}
	return oauththirdpartymwcli.GetOAuthThirdParties(ctx, conds, h.Offset, h.Limit)
}
