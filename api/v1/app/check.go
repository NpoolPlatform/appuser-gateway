package app

import (
	"context"

	appuserapp "github.com/NpoolPlatform/appuser-manager/pkg/client/app"
	"github.com/NpoolPlatform/appuser-middleware/api/v1/app"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	appmgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/app"
	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validate(ctx context.Context, info *appmwpb.AppReq) error {
	if info == nil {
		logger.Sugar().Errorw("validate", "err", "params is empty")
		return status.Error(codes.InvalidArgument, "params is empty")
	}

	err := app.Validate(ctx, info)
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if info.GetID() != "" {
		exist, err := appuserapp.ExistApp(ctx, info.GetID())
		if err != nil {
			logger.Sugar().Errorw("validate", "err", err)
			return status.Error(codes.Internal, err.Error())
		}
		if exist {
			logger.Sugar().Errorw("validate", "err", "app already exists")
			return status.Error(codes.AlreadyExists, "app already exists")
		}
	}

	exist, err := appuserapp.ExistAppConds(ctx, &appmgrpb.Conds{Name: &npool.StringVal{
		Value: info.GetName(),
		Op:    cruder.EQ,
	}})
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return status.Error(codes.AlreadyExists, "Logo is empty")
	}
	if exist {
		logger.Sugar().Errorw("validate", "err", "app name already exists")
		return status.Error(codes.AlreadyExists, "app name already exists")
	}

	return nil
}
