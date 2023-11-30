package auth

import (
	"context"

	authmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/authing/auth"
	authmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/authing/auth"
)

func (h *Handler) DeleteAuth(ctx context.Context) (*authmwpb.Auth, error) {
	if err := h.ExistAuth(ctx); err != nil {
		return nil, err
	}
	return authmwcli.DeleteAuth(ctx, *h.ID)
}
