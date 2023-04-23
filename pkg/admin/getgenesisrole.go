package admin

import (
	"context"
	"encoding/json"
	"fmt"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	servicename "github.com/NpoolPlatform/appuser-manager/pkg/servicename"
	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role"
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

	ids := []string{}
	for _, _role := range h.GenesisRoles {
		ids = append(ids, _role.ID)
	}
	infos, _, err := rolemwcli.GetRoles(ctx, &rolemwpb.Conds{
		IDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: ids},
	}, 0, int32(len(ids)))
	if err != nil {
		return nil, err
	}
	return infos, nil
}
