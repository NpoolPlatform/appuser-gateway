package appoauththirdparty

import (
	"context"

	oauththirdpartymwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/authing/oauth/appoauththirdparty"
	oauththirdpartymwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/authing/oauth/appoauththirdparty"
)

func (h *Handler) CreateOAuthThirdParty(ctx context.Context) (*oauththirdpartymwpb.OAuthThirdParty, error) {
	return oauththirdpartymwcli.CreateOAuthThirdParty(
		ctx,
		&oauththirdpartymwpb.OAuthThirdPartyReq{
			AppID:        &h.AppID,
			ThirdPartyID: h.ThirdPartyID,
			ClientID:     h.ClientID,
			ClientSecret: h.ClientSecret,
			CallbackURL:  h.CallbackURL,
		},
	)
}
