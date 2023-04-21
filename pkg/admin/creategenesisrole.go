package admin

import (
	"context"
	"encoding/json"
	"fmt"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	servicename "github.com/NpoolPlatform/appuser-manager/pkg/servicename"
	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role"
	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	uuid1 "github.com/NpoolPlatform/go-service-framework/pkg/const/uuid"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	rolemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

type createGenesisRoleHandler struct {
	*Handler
	roles []*rolemwpb.Role
}

func (h *createGenesisRoleHandler) getGenesisRoleConfig() error {
	str := config.GetStringValueWithNameSpace(
		servicename.ServiceDomain,
		constant.KeyGenesisRole,
	)
	if err := json.Unmarshal([]byte(str), &h.roles); err != nil {
		return err
	}
	if len(h.roles) == 0 {
		return fmt.Errorf("invalid genesis roles")
	}
	return nil
}

func (h *createGenesisRoleHandler) getGenesisRoles(ctx context.Context) (bool, error) {
	ids := []string{}
	for _, _role := range h.roles {
		ids = append(ids, _role.ID)
	}
	infos, _, err := rolemwcli.GetRoles(ctx, &rolemwpb.Conds{
		IDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: ids},
	}, 0, int32(len(ids)))
	if err != nil {
		return false, err
	}
	if len(infos) == 0 {
		return false, nil
	}
	h.roles = infos
	return true, nil
}

func (h *createGenesisRoleHandler) createGenesisRoles(ctx context.Context) error {
	reqs := []*rolemwpb.RoleReq{}
	createdBy := uuid1.InvalidUUIDStr
	defautl := false
	genesis := true

	for _, _role := range h.roles {
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

	h.roles = infos
	return nil
}

func (h *Handler) CreateGenesisRoles(ctx context.Context) ([]*rolemwpb.Role, error) {
	handler := &createGenesisRoleHandler{
		Handler: h,
	}
	if err := handler.getGenesisRoleConfig(); err != nil {
		return nil, err
	}
	created, err := handler.getGenesisRoles(ctx)
	if err != nil {
		return nil, err
	}
	if created {
		return handler.roles, nil
	}
	if err := handler.createGenesisRoles(ctx); err != nil {
		return nil, err
	}
	return handler.roles, nil
}
