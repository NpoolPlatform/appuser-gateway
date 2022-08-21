package user

import (
	"context"

	signmethod "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
)

func Login(
	ctx context.Context,
	appID, account, passwordHash string,
	accountType signmethod.SignMethodType,
	manMachineSpec, envSpec string,
) (
	*usermwpb.User, error,
) {
	return nil, nil
}
