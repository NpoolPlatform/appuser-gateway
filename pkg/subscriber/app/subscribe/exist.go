package appsubscribe

import (
	"context"
	"fmt"

	appsubscribemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/subscriber/app/subscribe"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	appsubscribemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/subscriber/app/subscribe"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) ExistAppSubscribe(ctx context.Context) error {
	exist, err := appsubscribemwcli.ExistAppSubscribeConds(ctx, &appsubscribemwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
	})
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid appsubscribe")
	}
	return nil
}
