package auth

import (
	"context"

	authmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/authing/auth"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	authmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/authing/auth"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) GetAuths(ctx context.Context) ([]*authmwpb.Auth, uint32, error) {
	return authmwcli.GetAuths(ctx, &authmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}, h.Offset, h.Limit)
}
