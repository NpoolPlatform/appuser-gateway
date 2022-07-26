package banappuser

import (
	"context"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	"github.com/NpoolPlatform/message/npool/appusergw/banappuser"
	banappusercrud "github.com/NpoolPlatform/message/npool/appusermgrv2/banappuser"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetBanAppUser(ctx context.Context, in *banappuser.GetBanAppUserRequest) (*banappuser.GetBanAppUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetBanApp")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Error("ID is invalid")
		return &banappuser.GetBanAppUserResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgIDInvalid)
	}

	span.AddEvent("call grpc GetBanAppUserV2")
	resp, err := grpc.GetBanAppUserV2(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Error("fail get ban app user:%v", err)
		return &banappuser.GetBanAppUserResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &banappuser.GetBanAppUserResponse{
		Info: resp,
	}, nil
}

func (s *Server) GetAppUserBanAppUser(ctx context.Context, in *banappuser.GetAppUserBanAppUserRequest) (*banappuser.GetAppUserBanAppUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppUserBanAppUser")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Error("UserID is invalid")
		return &banappuser.GetAppUserBanAppUserResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgUserIDInvalid)
	}

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Error("AppID is invalid")
		return &banappuser.GetAppUserBanAppUserResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgAppIDInvalid)
	}

	span.AddEvent("call grpc GetBanAppUserOnlyV2")
	resp, err := grpc.GetBanAppUserOnlyV2(ctx, &banappusercrud.Conds{
		AppID: &npool.StringVal{
			Value: in.GetTargetAppID(),
			Op:    cruder.EQ,
		},
		UserID: &npool.StringVal{
			Value: in.GetUserID(),
			Op:    cruder.EQ,
		},
	})
	if err != nil {
		logger.Sugar().Error("fail get ban app user:%v", err)
		return &banappuser.GetAppUserBanAppUserResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &banappuser.GetAppUserBanAppUserResponse{
		Info: resp,
	}, nil
}
