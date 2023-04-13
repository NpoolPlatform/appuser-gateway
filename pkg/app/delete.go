package app

import (
	"context"
	"fmt"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
)

func (h *Handler) DeleteApp(ctx context.Context) (info *appmwpb.App, err error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}

	return appmwcli.DeleteApp(ctx, *h.ID)
}
