package oauththirdparty

import (
	"context"

	oauththirdpartymwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/oauth/oauththirdparty"
	oauththirdpartymwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/oauth/oauththirdparty"
)

func (h *Handler) DeleteOAuthThirdParty(ctx context.Context) (*oauththirdpartymwpb.OAuthThirdParty, error) {
	if err := h.ExistOAuthThirdParty(ctx); err != nil {
		return nil, err
	}
	return oauththirdpartymwcli.DeleteOAuthThirdParty(ctx, *h.ID)
}
