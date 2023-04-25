//nolint:dupl
package admin

import (
	"context"
	"encoding/json"
	"fmt"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role"
	servicename "github.com/NpoolPlatform/appuser-middleware/pkg/servicename"
	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	rolemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

type getGenesisRoleHandler struct {
	*Handler
}

func (h *getGenesisRoleHandler) getGenesisRoleConfig() error {
	str := config.GetStringValueWithNameSpace(
		servicename.ServiceDomain,
		constant.KeyGenesisRole,
	)
	if err := json.Unmarshal([]byte(str), &h.GenesisRoles); err != nil {
		return err
	}
	if len(h.GenesisRoles) == 0 {
		return fmt.Errorf("invalid genesis roles")
	}
	return nil
}

func (h *Handler) GetGenesisRoles(ctx context.Context) ([]*rolemwpb.Role, error) {
	handler := &getGenesisRoleHandler{
		Handler: h,
	}
	if err := handler.getGenesisRoleConfig(); err != nil {
		return nil, err
	}

	appIDs := []string{}
	roles := []string{}
	for _, _role := range h.GenesisRoles {
		appIDs = append(appIDs, _role.AppID)
		roles = append(roles, _role.Role)
	}

	infos, _, err := rolemwcli.GetRoles(ctx, &rolemwpb.Conds{
		AppIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: appIDs},
		Roles:  &basetypes.StringSliceVal{Op: cruder.IN, Value: roles},
	}, 0, int32(len(appIDs)*len(roles)))
	if err != nil {
		return nil, err
	}
	return infos, nil
}
