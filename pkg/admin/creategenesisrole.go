package admin

import (
	"context"

	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role"
	rolemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role"

	"github.com/google/uuid"
)

type createGenesisRoleHandler struct {
	*Handler
}

func (h *createGenesisRoleHandler) createGenesisRoles(ctx context.Context) error {
	reqs := []*rolemwpb.RoleReq{}
	createdBy := uuid.Nil.String()
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

func (h *createGenesisRoleHandler) patchGenesisRole(ctx context.Context) error {
	for _, _role := range h.GenesisRoles {
		if _role.Genesis {
			continue
		}
		genesis := true
		_, err := rolemwcli.UpdateRole(ctx, &rolemwpb.RoleReq{
			ID:      &_role.ID,
			AppID:   &_role.AppID,
			Genesis: &genesis,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Handler) CreateGenesisRoles(ctx context.Context) ([]*rolemwpb.Role, error) {
	handler := &createGenesisRoleHandler{
		Handler: h,
	}
	_roles, err := h.GetGenesisRoles(ctx)
	if err != nil {
		return nil, err
	}
	if len(_roles) > 0 {
		h.GenesisRoles = _roles
		_ = handler.patchGenesisRole(ctx)
		return _roles, nil
	}
	if err := handler.createGenesisRoles(ctx); err != nil {
		return nil, err
	}
	return handler.GenesisRoles, nil
}
