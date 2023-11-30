package admin

import (
	"context"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"

	"github.com/google/uuid"
)

type createGenesisAppHandler struct {
	*Handler
}

func (h *createGenesisAppHandler) createGenesisApps(ctx context.Context) error {
	createdBy := uuid.Nil.String()
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
	_apps, err := h.GetGenesisApps(ctx)
	if err != nil {
		return nil, err
	}
	if len(_apps) > 0 {
		return _apps, nil
	}
	if err := handler.createGenesisApps(ctx); err != nil {
		return nil, err
	}
	return h.GenesisApps, nil
}
