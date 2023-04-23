package admin

import (
	"context"
	"encoding/json"
	"fmt"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	servicename "github.com/NpoolPlatform/appuser-middleware/pkg/servicename"
	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	uuid1 "github.com/NpoolPlatform/go-service-framework/pkg/const/uuid"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

type createGenesisAppHandler struct {
	*Handler
}

func (h *createGenesisAppHandler) getGenesisApps(ctx context.Context) (bool, error) {
	ids := []string{}
	for _, _app := range h.GenesisApps {
		ids = append(ids, _app.ID)
	}
	infos, _, err := appmwcli.GetApps(ctx, &appmwpb.Conds{
		IDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: ids},
	}, 0, int32(len(ids)))
	if err != nil {
		return false, err
	}
	if len(infos) == 0 {
		return false, nil
	}
	h.GenesisApps = infos
	return true, nil
}

func (h *createGenesisAppHandler) createGenesisApps(ctx context.Context) error {
	createdBy := uuid1.InvalidUUIDStr
	logo := "NOT SET"
	reqs := []*appmwpb.AppReq{}

	for _, _app := range h.GenesisApps {
		reqs = append(reqs, &appmwpb.AppReq{
			ID:          &_app.ID,
			CreatedBy:   &createdBy,
			Name:        &_app.Name,
			Logo:        &logo,
			Description: &_app.Description,
		})
	}

	infos, err := appmwcli.CreateApps(ctx, reqs)
	if err != nil {
		return err
	}

	h.GenesisApps = infos
	return nil
}

func (h *Handler) CreateAdminApps(ctx context.Context) ([]*appmwpb.App, error) {
	handler := &createGenesisAppHandler{
		Handler: h,
	}
	if err := handler.GetGenesisAppConfig(); err != nil {
		return nil, err
	}
	created, err := handler.getGenesisApps(ctx)
	if err != nil {
		return nil, err
	}
	if created {
		return handler.GenesisApps, nil
	}
	if err := handler.createGenesisApps(ctx); err != nil {
		return nil, err
	}

	return handler.GenesisApps, nil
}
