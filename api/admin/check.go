package admin

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	approlemgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approle"
	appusermgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/appuser"

	admingwpb "github.com/NpoolPlatform/message/npool/appuser/gw/v1/admin"

	approlemgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approle"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/appuser"

	commonpb "github.com/NpoolPlatform/message/npool"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func validate(ctx context.Context, info *admingwpb.CreateGenesisUserRequest) error {
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

	roles, _, err := approlemgrcli.GetAppRoles(ctx, &approlemgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: info.GetTargetAppID(),
		},
	}, 0, 100) // nolint
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	if len(roles) == 0 {
		return status.Error(codes.Internal, "fail get app role")
	}

	exist, err := appusermgrcli.ExistAppUserConds(ctx, &appusermgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: info.GetTargetAppID(),
		},
		EmailAddress: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: info.GetEmailAddress(),
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
