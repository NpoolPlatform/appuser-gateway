package appsubscribe

import (
	"context"

	appsubscribemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/subscriber/app/subscribe"
	appsubscribemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/subscriber/app/subscribe"
)

func (h *Handler) DeleteAppSubscribe(ctx context.Context) (*appsubscribemwpb.AppSubscribe, error) {
	return appsubscribemwcli.DeleteAppSubscribe(ctx, *h.ID)
}
