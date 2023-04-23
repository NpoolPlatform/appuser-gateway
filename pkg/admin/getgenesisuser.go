package admin

import (
	"context"

	roleusermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role/user"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	roleusermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role/user"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) GetGenesisUsers(ctx context.Context) error {
	const maxGenesisUsers = int32(20)

	infos, _, err := roleusermwcli.GetUsers(ctx, &roleusermwpb.Conds{
		Genesis: &basetypes.BoolVal{Op: cruder.EQ, Value: true},
	}, 0, maxGenesisUsers)
	if err != nil {
		return err
	}
	h.GenesisRoleUsers = infos

	ids := []string{}
	for _, info := range infos {
		ids = append(ids, info.UserID)
	}

	users, _, err := usermwcli.GetUsers(ctx, &usermwpb.Conds{
		IDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: ids},
	}, 0, int32(len(ids)))
	if err != nil {
		return err
	}
	h.GenesisUsers = users

	return nil
}
