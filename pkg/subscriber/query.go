package subscriber

import (
	"context"

	subscribermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/subscriber"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	subscribermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/subscriber"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) GetSubscriberes(ctx context.Context) ([]*subscribermwpb.Subscriber, uint32, error) {
	conds := &subscribermwpb.Conds{}
	if h.AppID != nil {
		conds.AppID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID}
	}
	return subscribermwcli.GetSubscriberes(ctx, conds, h.Offset, h.Limit)
}
