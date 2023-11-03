package subscriber

import (
	"context"
	"fmt"

	subscribermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/subscriber"
	appsubscribemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/subscriber/app/subscribe"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	subscribermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/subscriber"
	appsubscribemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/subscriber/app/subscribe"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) CreateSubscriber(ctx context.Context) (*subscribermwpb.Subscriber, error) {
	appID := h.AppID

	if h.SubscribeAppID != nil {
		exist, err := appsubscribemwcli.ExistAppSubscribeConds(ctx, &appsubscribemwpb.Conds{
			AppID:          &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			SubscribeAppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.SubscribeAppID},
		})
		if err != nil {
			return nil, err
		}
		if !exist {
			return nil, fmt.Errorf("permission denied")
		}
		appID = h.SubscribeAppID
	}

	return subscribermwcli.CreateSubscriber(ctx, &subscribermwpb.SubscriberReq{
		AppID:        appID,
		EmailAddress: h.EmailAddress,
	})
}
