//nolint:nolintlint,dupl
package app

import (
	"context"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	"github.com/NpoolPlatform/message/npool/appusergw/app"
	appcrud "github.com/NpoolPlatform/message/npool/appusermgrv2/app"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateApp(ctx context.Context, in *app.CreateAppRequest) (*app.CreateAppResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateAppV2")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if in.GetInfo().GetID() != "" {
		span.AddEvent("call grpc ExistAppV2")
		exist, err := grpc.ExistAppV2(ctx, in.GetInfo().GetID())
		if err != nil {
			logger.Sugar().Errorw("fail check app : %v", err)
			return &app.CreateAppResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
		}
		if exist {
			logger.Sugar().Errorw("app already exists")
			return &app.CreateAppResponse{}, status.Error(npool.ErrAlreadyExists, appusergw.ErrMsgAppAlreadyExists)
		}
	}

	err = validate(in.GetInfo())
	if err != nil {
		return nil, err
	}

	span.AddEvent("call grpc ExistAppCondsV2")
	exist, err := grpc.ExistAppCondsV2(ctx, &appcrud.Conds{Name: &npool.StringVal{
		Value: in.GetInfo().GetName(),
		Op:    cruder.EQ,
	}})
	if err != nil {
		logger.Sugar().Errorw("fail check app name: %v", err)
		return &app.CreateAppResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}
	if exist {
		logger.Sugar().Errorw("app name already exists")
		return &app.CreateAppResponse{}, status.Error(npool.ErrAlreadyExists, appusergw.ErrMsgAppAlreadyExists)
	}

	span.AddEvent("call grpc CreateAppV2")
	resp, err := grpc.CreateAppV2(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("fail create app: %v", err)
		return &app.CreateAppResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &app.CreateAppResponse{
		Info: resp,
	}, nil
}
