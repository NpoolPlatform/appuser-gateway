package subscriber

import (
	"context"

	subscribermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/subscriber"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	subscribermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/subscriber"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) GetSubscriberes(ctx context.Context) ([]*subscribermwpb.Subscriber, uint32, error) {
	return subscribermwcli.GetSubscriberes(
		ctx,
		&subscribermwpb.Conds{
			AppID: &basetypes.StringVal{Op: cruder.EQ, Value: h.AppID},
		},
		h.Offset,
		h.Limit,
	)
}
