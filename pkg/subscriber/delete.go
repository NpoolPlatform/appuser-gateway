package subscriber

import (
	"context"
	"fmt"

	subscribermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/subscriber"
	subscribermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/subscriber"
)

func (h *Handler) DeleteSubscriber(ctx context.Context) (*subscribermwpb.Subscriber, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}
	return subscribermwcli.DeleteSubscriber(ctx, *h.ID)
}
