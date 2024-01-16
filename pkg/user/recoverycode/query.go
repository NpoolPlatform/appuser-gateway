package recoverycode

import (
	"context"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	recoverycodemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user/recoverycode"
	npool "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user/recoverycode"
)

func (h *Handler) GetRecoveryCodes(ctx context.Context) ([]*npool.RecoveryCode, uint32, error) {
	infos, total, err := recoverycodemwcli.GetRecoveryCodes(ctx, &npool.Conds{
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		UserID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID},
	}, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	return infos, total, nil
}
