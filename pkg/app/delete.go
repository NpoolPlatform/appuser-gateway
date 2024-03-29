package app

import (
	"context"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
)

func (h *Handler) DeleteApp(ctx context.Context) (info *appmwpb.App, err error) {
	if err := h.ExistApp(ctx); err != nil {
		return nil, err
	}
	return appmwcli.DeleteApp(ctx, *h.ID)
}
