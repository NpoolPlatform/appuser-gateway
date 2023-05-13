package auth

import (
	"context"
	"fmt"

	authmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/authing/auth"
	authmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/authing/auth"
)

func (h *Handler) DeleteAuth(ctx context.Context) (*authmwpb.Auth, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}
	return authmwcli.DeleteAuth(ctx, *h.ID)
}
