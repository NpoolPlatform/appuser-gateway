package roleuser

import (
	"context"
	"fmt"
	"github.com/NpoolPlatform/api-manager/pkg/db/ent"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusermgrv2/approleuser"
)

func CreateRoleUser(ctx context.Context, in *approleuser.AppRoleUserReq) (*approleuser.AppRoleUser, error) {
	role, err := grpc.GetAppRoleV2(ctx, in.GetRoleID())
	if err != nil {
		return nil, err
	}

	if role.GetRole() == constant.GenesisRole {
		return nil, fmt.Errorf("permission denied")
	}

	resp, err := grpc.GetAppRoleUserOnlyV2(ctx, &approleuser.Conds{
		UserID: &npool.StringVal{
			Value: in.GetUserID(),
			Op:    cruder.EQ,
		},
		AppID: &npool.StringVal{
			Value: in.GetAppID(),
			Op:    cruder.EQ,
		},
		RoleID: &npool.StringVal{
			Value: in.GetRoleID(),
			Op:    cruder.EQ,
		},
	})
	if err != nil {
		if ent.IsNotFound(err) {
			resp, err = grpc.CreateAppRoleUserV2(ctx, in)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("fail get app role user: %v", err)
		}
	}

	return resp, nil
}
