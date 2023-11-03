package oauththirdparty

import (
	"context"

	oauththirdpartymwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/oauth/oauththirdparty"
	oauththirdpartymwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/oauth/oauththirdparty"
)

func (h *Handler) UpdateOAuthThirdParty(ctx context.Context) (*oauththirdpartymwpb.OAuthThirdParty, error) {
	return oauththirdpartymwcli.UpdateOAuthThirdParty(ctx, &oauththirdpartymwpb.OAuthThirdPartyReq{
		ID:             h.ID,
		ClientName:     h.ClientName,
		ClientTag:      h.ClientTag,
		ClientLogoURL:  h.ClientLogoURL,
		ClientOAuthURL: h.ClientOAuthURL,
		ResponseType:   h.ResponseType,
		Scope:          h.Scope,
	})
}
