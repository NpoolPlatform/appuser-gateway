package secret

import (
	"context"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	"github.com/NpoolPlatform/message/npool/appusergw/appusersecret"
	appusersecretcrud "github.com/NpoolPlatform/message/npool/appusermgrv2/appusersecret"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetSecret(ctx context.Context, in *appusersecret.GetSecretRequest) (*appusersecret.GetSecretResponse, error) {
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
		return &appusersecret.GetSecretResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgIDInvalid)
	}

	span.AddEvent("call grpc GetBanAppV2")
	resp, err := grpc.GetAppUserSecretV2(ctx, in.GetID())
	if err != nil {
		return &appusersecret.GetSecretResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &appusersecret.GetSecretResponse{
		Info: resp,
	}, nil
}

func (s *Server) GetAppUserSecret(ctx context.Context, in *appusersecret.GetAppUserSecretRequest) (*appusersecret.GetAppUserSecretResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetBanApp")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Error("UserID is invalid")
		return &appusersecret.GetAppUserSecretResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgUserIDInvalid)
	}

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Error("AppID is invalid")
		return &appusersecret.GetAppUserSecretResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgAppIDInvalid)
	}

	span.AddEvent("call grpc GetAppUserSecretOnlyV2")
	resp, err := grpc.GetAppUserSecretOnlyV2(ctx, &appusersecretcrud.Conds{
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
		return &appusersecret.GetAppUserSecretResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &appusersecret.GetAppUserSecretResponse{
		Info: resp,
	}, nil
}
