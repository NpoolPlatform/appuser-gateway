package oauththirdparty

import (
	"context"

	oauththirdpartymwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/authing/oauth/oauththirdparty"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	oauththirdpartymwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/authing/oauth/oauththirdparty"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) GetOAuthThirdParties(ctx context.Context) ([]*oauththirdpartymwpb.OAuthThirdParty, uint32, error) {
	return oauththirdpartymwcli.GetOAuthThirdParties(
		ctx,
		&oauththirdpartymwpb.Conds{
			ClientName: &basetypes.Int32Val{Op: cruder.EQ, Value: int32(*h.ClientName)},
		},
		h.Offset,
		h.Limit,
	)
}
