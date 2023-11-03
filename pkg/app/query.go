package app

import (
	"context"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) GetApp(ctx context.Context) (*appmwpb.App, error) {
	return appmwcli.GetApp(ctx, *h.EntID)
}

func (h *Handler) GetApps(ctx context.Context) ([]*appmwpb.App, uint32, error) {
	conds := &appmwpb.Conds{}
	if h.EntID != nil {
		conds.EntID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID}
	}
	if len(h.EntIDs) > 0 {
		conds.EntIDs = &basetypes.StringSliceVal{Op: cruder.IN, Value: h.EntIDs}
	}
	if h.CreatedBy != nil {
		conds.CreatedBy = &basetypes.StringVal{Op: cruder.EQ, Value: *h.CreatedBy}
	}
	if h.Name != nil {
		conds.Name = &basetypes.StringVal{Op: cruder.EQ, Value: *h.Name}
	}
	return appmwcli.GetApps(ctx, conds, h.Offset, h.Limit)
}
