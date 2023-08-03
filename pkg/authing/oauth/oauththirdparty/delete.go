package oauththirdparty

import (
	"context"
	"fmt"

	oauththirdpartymwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/authing/oauth/oauththirdparty"
	oauththirdpartymwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/authing/oauth/oauththirdparty"
)

func (h *Handler) DeleteOAuthThirdParty(ctx context.Context) (*oauththirdpartymwpb.OAuthThirdParty, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}
	return oauththirdpartymwcli.DeleteOAuthThirdParty(ctx, *h.ID)
}
