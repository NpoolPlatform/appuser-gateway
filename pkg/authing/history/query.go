package history

import (
	"context"

	historymwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/authing/history"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	historymwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/authing/history"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) GetAuthHistories(ctx context.Context) ([]*historymwpb.History, uint32, error) {
	return historymwcli.GetHistories(ctx, &historymwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}, h.Offset, h.Limit)
}
