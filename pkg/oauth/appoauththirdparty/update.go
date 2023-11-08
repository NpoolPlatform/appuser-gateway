package appoauththirdparty

import (
	"context"

	oauththirdpartymwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/oauth/appoauththirdparty"
	oauththirdpartymwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/oauth/appoauththirdparty"
)

func (h *Handler) UpdateOAuthThirdParty(ctx context.Context) (*oauththirdpartymwpb.OAuthThirdParty, error) {
	if err := h.ExistOAuthThirdParty(ctx); err != nil {
		return nil, err
	}
	return oauththirdpartymwcli.UpdateOAuthThirdParty(ctx, &oauththirdpartymwpb.OAuthThirdPartyReq{
		ID:           h.ID,
		AppID:        h.AppID,
		ThirdPartyID: h.ThirdPartyID,
		ClientID:     h.ClientID,
		ClientSecret: h.ClientSecret,
		CallbackURL:  h.CallbackURL,
	})
}
