package oauththirdparty

import (
	"context"

	oauththirdpartymwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/oauth/oauththirdparty"
	oauththirdpartymwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/oauth/oauththirdparty"
)

func (h *Handler) GetOAuthThirdParties(ctx context.Context) ([]*oauththirdpartymwpb.OAuthThirdParty, uint32, error) {
	return oauththirdpartymwcli.GetOAuthThirdParties(
		ctx,
		&oauththirdpartymwpb.Conds{},
		h.Offset,
		h.Limit,
	)
}
