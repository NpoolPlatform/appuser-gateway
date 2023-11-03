package appsubscribe

import (
	"context"

	appsubscribemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/subscriber/app/subscribe"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	appsubscribemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/subscriber/app/subscribe"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) GetAppSubscribes(ctx context.Context) ([]*appsubscribemwpb.AppSubscribe, uint32, error) {
	conds := &appsubscribemwpb.Conds{}
	if h.AppID != nil {
		conds.AppID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID}
	}
	return appsubscribemwcli.GetAppSubscribes(ctx, conds, h.Offset, h.Limit)
}
