package roleuser

import (
	"context"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	apiapproleuser "github.com/NpoolPlatform/appuser-manager/api/v2/approleuser"
	approlemgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approle"
	approleusermgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approleuser"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	approleuserpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approleuser"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validate(ctx context.Context, info *approleuserpb.AppRoleUserReq) error {
	err := apiapproleuser.Validate(info)
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return status.Error(codes.InvalidArgument, err.Error())
	}

	role, err := approlemgrcli.GetAppRole(ctx, info.GetRoleID())
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if role.GetRole() == constant.GenesisAppID {
		logger.Sugar().Errorw("validate", "Role", role.GetRole())
		return status.Error(codes.PermissionDenied, "permission denied")
	}

	exist, err := approleusermgrcli.ExistAppRoleUserConds(ctx, &approleuserpb.Conds{
		AppID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: info.GetAppID(),
		},
		RoleID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: info.GetRoleID(),
		},
		UserID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: info.GetUserID(),
		},
	})
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return status.Error(codes.Internal, err.Error())
	}

	if exist {
		logger.Sugar().Errorw("validate", "exist", exist)
		return status.Error(codes.AlreadyExists, "app role user already exists")
	}

	return err
}
