package app

import (
	"context"
	"fmt"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) GetApp(ctx context.Context) (*appmwpb.App, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}
	return appmwcli.GetApp(ctx, *h.ID)
}

func (h *Handler) GetApps(ctx context.Context) ([]*appmwpb.App, uint32, error) {
	conds := &appmwpb.Conds{}
	if h.ID != nil {
		conds.ID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.ID}
	}
	if len(h.IDs) > 0 {
		conds.IDs = &basetypes.StringSliceVal{Op: cruder.IN, Value: h.IDs}
	}
	if h.CreatedBy != nil {
		conds.CreatedBy = &basetypes.StringVal{Op: cruder.EQ, Value: *h.CreatedBy}
	}
	if h.Name != nil {
		conds.Name = &basetypes.StringVal{Op: cruder.EQ, Value: *h.Name}
	}
	return appmwcli.GetApps(ctx, conds, h.Offset, h.Limit)
}
