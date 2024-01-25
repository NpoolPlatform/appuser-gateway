package recoverycode

import (
	"context"

	recoverycodemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user/recoverycode"
	npool "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user/recoverycode"
)

func (h *Handler) GenerateRecoveryCodes(ctx context.Context) ([]*npool.RecoveryCode, error) {
	return recoverycodemwcli.GenerateRecoveryCodes(ctx, &npool.RecoveryCodeReq{
		AppID:  h.AppID,
		UserID: h.UserID,
	})
}
