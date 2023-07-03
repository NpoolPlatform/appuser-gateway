package appsubscribe

import (
	"context"

	appsubscribemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/subscriber/app/subscribe"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	appsubscribemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/subscriber/app/subscribe"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) GetAppSubscribes(ctx context.Context) ([]*appsubscribemwpb.AppSubscribe, uint32, error) {
	return appsubscribemwcli.GetAppSubscribes(
		ctx,
		&appsubscribemwpb.Conds{
			AppID: &basetypes.StringVal{Op: cruder.EQ, Value: h.AppID},
		},
		h.Offset,
		h.Limit,
	)
}
