package appoauththirdparty

import (
	"context"
	"fmt"

	oauththirdpartymwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/oauth/appoauththirdparty"
	oauththirdpartymwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/oauth/appoauththirdparty"
)

func (h *Handler) UpdateOAuthThirdParty(ctx context.Context) (*oauththirdpartymwpb.OAuthThirdParty, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}

	return oauththirdpartymwcli.UpdateOAuthThirdParty(
		ctx,
		&oauththirdpartymwpb.OAuthThirdPartyReq{
			ID:           h.ID,
			AppID:        &h.AppID,
			ThirdPartyID: h.ThirdPartyID,
			ClientID:     h.ClientID,
			ClientSecret: h.ClientSecret,
			CallbackURL:  h.CallbackURL,
		},
	)
}
