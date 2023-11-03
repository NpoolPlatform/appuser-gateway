package subscriber

import (
	"context"

	subscribermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/subscriber"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	subscribermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/subscriber"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) DeleteSubscriber(ctx context.Context) (*subscribermwpb.Subscriber, error) {
	info, err := subscribermwcli.GetSubscriberOnly(ctx, &subscribermwpb.Conds{
		AppID:        &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EmailAddress: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EmailAddress},
	})
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}
	return subscribermwcli.DeleteSubscriber(ctx, info.ID)
}
