package admin

import (
	"context"

	"github.com/NpoolPlatform/appuser-manager/pkg/client/approle"
	"github.com/NpoolPlatform/appuser-manager/pkg/client/approleuser"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/admin"
	approlepb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approle"
	approleuserpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approleuser"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validate(ctx context.Context, info *admin.CreateGenesisUserRequest) error {
	if info.GetTargetAppID() == "" {
		logger.Sugar().Errorw("validate", "TargetAppID", info.GetTargetAppID())
		return status.Error(codes.InvalidArgument, "AppID is empty")
	}

	if _, err := uuid.Parse(info.GetTargetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "GetTargetAppID", info.GetTargetAppID(), "error", err)
		return status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	if info.GetEmailAddress() == "" {
		logger.Sugar().Errorw("validate", "GetEmailAddress", info.GetEmailAddress())
		return status.Error(codes.InvalidArgument, "EmailAddress is empty")
	}

	if info.GetPasswordHash() == "" {
		logger.Sugar().Errorw("validate", "GetPasswordHash", info.GetPasswordHash())
		return status.Error(codes.InvalidArgument, "PasswordHash is empty")
	}

	resp, err := approle.GetAppRoleOnly(ctx, &approlepb.Conds{
		AppID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: info.GetTargetAppID(),
		},
	})
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	if resp == nil {
		return status.Error(codes.Internal, "fail get app role")
	}

	exist, err := approleuser.ExistAppRoleUserConds(ctx, &approleuserpb.Conds{
		AppID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: info.GetTargetAppID(),
		},
		RoleID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: resp.ID,
		},
	})
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err.Error())
		return status.Error(codes.Internal, err.Error())
	}
	if exist {
		logger.Sugar().Errorw("validate", "err", "genesis user already exists")
		return status.Error(codes.AlreadyExists, "genesis user already exists")
	}

	return nil
}
