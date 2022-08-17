package role

import (
	"context"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	apiapproleuser "github.com/NpoolPlatform/appuser-manager/api/v2/approleuser"
	approlemgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approle"
	approleusermgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approleuser"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/role"
	approlepb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approle"
	approleuserpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approleuser"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validate(ctx context.Context, info *role.CreateRoleRequest) error {
	if info == nil {
		logger.Sugar().Errorw("validate", "err", "params is empty")
		return status.Error(codes.InvalidArgument, "params is empty")
	}

	if _, err := uuid.Parse(info.GetUserID()); err != nil {
		logger.Sugar().Errorw("validate", "userId", info.GetUserID(), "err", err)
		return status.Error(codes.InvalidArgument, "UserID is invalid")
	}

	if info.GetRoleName() == constant.GenesisRole {
		logger.Sugar().Errorw("validate", "RoleName", info.GetRoleName())
		return status.Error(codes.PermissionDenied, "permission denied")
	}

	exist, err := approlemgrcli.ExistAppRoleConds(ctx, &approlepb.Conds{
		AppID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: info.GetAppID(),
		},
		Role: &npool.StringVal{
			Op:    cruder.EQ,
			Value: info.GetRoleName(),
		},
	})
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	if exist {
		return status.Error(codes.AlreadyExists, "role name already exists")
	}
	return nil
}

func validateRoleUser(ctx context.Context, info *approleuserpb.AppRoleUserReq) error {
	err := apiapproleuser.Validate(info)
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return status.Error(codes.InvalidArgument, err.Error())
	}

	appRole, err := approlemgrcli.GetAppRole(ctx, info.GetRoleID())
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if appRole.GetRole() == constant.GenesisAppID {
		logger.Sugar().Errorw("validate", "Role", appRole.GetRole())
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
