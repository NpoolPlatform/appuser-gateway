package user

import (
	"context"

	hismwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user/login/history"
	hismwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user/login/history"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) GetLoginHistories(ctx context.Context) ([]*hismwpb.History, uint32, error) {
	return hismwcli.GetHistories(ctx, &hismwpb.Conds{
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		UserID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID},
	}, h.Offset, h.Limit)
}
