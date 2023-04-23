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

func (h *Handler) GetGenesisAppConfig() error {
	str := config.GetStringValueWithNameSpace(
		servicename.ServiceDomain,
		constant.KeyGenesisApp,
	)
	if err := json.Unmarshal([]byte(str), &h.GenesisApps); err != nil {
		return err
	}
	if len(h.GenesisApps) == 0 {
		return fmt.Errorf("invalid genesis app")
	}
	return nil
}

func (h *Handler) getGenesisApps(ctx context.Context) (bool, error) {
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
