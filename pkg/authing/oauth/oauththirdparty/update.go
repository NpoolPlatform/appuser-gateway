package oauththirdparty

import (
	"context"
	"fmt"

	oauththirdpartymwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/authing/oauth/oauththirdparty"
	oauththirdpartymwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/authing/oauth/oauththirdparty"
)

func (h *Handler) UpdateOAuthThirdParty(ctx context.Context) (*oauththirdpartymwpb.OAuthThirdParty, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}

	return oauththirdpartymwcli.UpdateOAuthThirdParty(
		ctx,
		&oauththirdpartymwpb.OAuthThirdPartyReq{
			ID:             h.ID,
			ClientName:     h.ClientName,
			ClientTag:      h.ClientTag,
			ClientLogoURL:  h.ClientLogoURL,
			ClientOAuthURL: h.ClientOAuthURL,
			ResponseType:   h.ResponseType,
			Scope:          h.Scope,
		},
	)
}
