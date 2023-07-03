package appsubscribe

import (
	"context"
	"fmt"

	appsubscribemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/subscriber/app/subscribe"
	appsubscribemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/subscriber/app/subscribe"
)

func (h *Handler) DeleteAppSubscribe(ctx context.Context) (*appsubscribemwpb.AppSubscribe, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid appsubscribe id")
	}
	return appsubscribemwcli.DeleteAppSubscribe(ctx, *h.ID)
}
