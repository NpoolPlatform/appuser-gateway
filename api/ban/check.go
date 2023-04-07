package ban

import (
	"context"

	"github.com/NpoolPlatform/appuser-manager/api/banapp"
	"github.com/NpoolPlatform/appuser-manager/api/banappuser"
	banappmgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/banapp"
	banappusermgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/banappuser"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	banapppb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/banapp"
	banappuserpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/banappuser"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validate(ctx context.Context, info *banapppb.BanAppReq) error {
	if info == nil {
		logger.Sugar().Errorw("validate", "err", "params is empty")
		return status.Error(codes.InvalidArgument, "params is empty")
	}

	err := banapp.Validate(info)
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return status.Error(codes.InvalidArgument, err.Error())
	}

	exist, err := banappmgrcli.ExistBanAppConds(ctx, &banapppb.Conds{
		AppID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: info.GetAppID(),
		},
	})
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return status.Error(codes.Internal, err.Error())
	}

	if exist {
		logger.Sugar().Errorw("validate", "err", "already exist")
		return status.Error(codes.AlreadyExists, "already exist")
	}
	return nil
}

func validateBanUser(ctx context.Context, info *banappuserpb.BanAppUserReq) error {
	if info == nil {
		logger.Sugar().Errorw("validate", "err", "params is empty")
		return status.Error(codes.InvalidArgument, "params is empty")
	}

	err := banappuser.Validate(info)
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return status.Error(codes.InvalidArgument, err.Error())
	}

	exist, err := banappusermgrcli.ExistBanAppUserConds(ctx, &banappuserpb.Conds{
		AppID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: info.GetAppID(),
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
		logger.Sugar().Errorw("validate", "err", "already exist")
		return status.Error(codes.AlreadyExists, "already exist")
	}

	return nil
}
