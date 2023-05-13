package subscriber

import (
	"context"

	subscribermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/subscriber"
	subscribermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/subscriber"
)

func (h *Handler) CreateSubscriber(ctx context.Context) (*subscribermwpb.Subscriber, error) {
	return subscribermwcli.CreateSubscriber(
		ctx,
		&subscribermwpb.SubscriberReq{
			AppID:        &h.AppID,
			EmailAddress: &h.EmailAddress,
		})
}
