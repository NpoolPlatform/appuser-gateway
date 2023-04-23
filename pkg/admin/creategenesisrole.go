package admin

import (
	"context"

	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role"
	uuid1 "github.com/NpoolPlatform/go-service-framework/pkg/const/uuid"
	rolemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role"
)

type createGenesisRoleHandler struct {
	*Handler
}

func (h *createGenesisRoleHandler) createGenesisRoles(ctx context.Context) error {
	reqs := []*rolemwpb.RoleReq{}
	createdBy := uuid1.InvalidUUIDStr
	defautl := false
	genesis := true

	for _, _role := range h.GenesisRoles {
		reqs = append(reqs, &rolemwpb.RoleReq{
			AppID:       &_role.AppID,
			CreatedBy:   &createdBy,
			Role:        &_role.Role,
			Description: &_role.Description,
			Default:     &defautl,
			Genesis:     &genesis,
		})
	}
	infos, err := rolemwcli.CreateRoles(ctx, reqs)
	if err != nil {
		return err
	}

	h.GenesisRoles = infos
	return nil
}

func (h *Handler) CreateGenesisRoles(ctx context.Context) ([]*rolemwpb.Role, error) {
	handler := &createGenesisRoleHandler{
		Handler: h,
	}
	if err := h.GetGenesisRoleConfig(); err != nil {
		return nil, err
	}
	created, err := h.GetGenesisRoles(ctx)
	if err != nil {
		return nil, err
	}
	if created {
		return handler.GenesisRoles, nil
	}
	if err := handler.createGenesisRoles(ctx); err != nil {
		return nil, err
	}
	return handler.GenesisRoles, nil
}
