package admin

import (
	"context"

	roleusermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role/user"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	roleusermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) GetGenesisUsers(ctx context.Context) error {
	const maxGenesisUsers = int32(20)
	infos, _, err := roleusermwcli.GetUsers(ctx, &roleusermwpb.Conds{
		AppID:   &basetypes.StringVal{Op: cruder.EQ, Value: h.AppID},
		Genesis: &basetypes.BoolVal{Op: cruder.EQ, Value: true},
	}, 0, maxGenesisUsers)
	if err != nil {
		return err
	}
	h.GenesisUsers = infos
	return nil
}
