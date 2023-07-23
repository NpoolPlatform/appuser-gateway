package oauththirdparty

import (
	"context"

	oauththirdpartymwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/authing/oauth/oauththirdparty"
	oauththirdpartymwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/authing/oauth/oauththirdparty"
)

func (h *Handler) CreateOAuthThirdParty(ctx context.Context) (*oauththirdpartymwpb.OAuthThirdParty, error) {
	return oauththirdpartymwcli.CreateOAuthThirdParty(
		ctx,
		&oauththirdpartymwpb.OAuthThirdPartyReq{
			ClientName:     h.ClientName,
			ClientTag:      h.ClientTag,
			ClientLogoURL:  h.ClientLogoURL,
			ClientOAuthURL: h.ClientOAuthURL,
			ResponseType:   h.ResponseType,
			Scope:          h.Scope,
		},
	)
}
