package api

import (
	"context"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergateway/app"
	appcrud "github.com/NpoolPlatform/message/npool/appusermgrv2/app"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func checkAppInfo(info *appcrud.AppReq) error {

	if _, err := uuid.Parse(info.GetCreatedBy()); err != nil {
		logger.Sugar().Error("CreatedBy is invalid")
		return status.Error(npool.ErrParams, app.ErrMsgAppCreatedByInvalid)
	}

	if info.Name == nil {
		logger.Sugar().Error("Name is empty")
		return status.Error(npool.ErrParams, app.ErrMsgAppNameEmpty)
	}

	if info.GetLogo() == "" {
		logger.Sugar().Error("Logo is empty")
		return status.Error(npool.ErrParams, app.ErrMsgAppLogoEmpty)
	}

	return nil
}

func (s *AppServer) CreateApp(ctx context.Context, in *app.CreateAppRequest) (*app.CreateAppResponse, error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateAppV2")
	defer span.End()
	var err error
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
			return &app.CreateAppResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
		}
		if exist {
			return &app.CreateAppResponse{}, status.Error(npool.ErrAlreadyExists, app.ErrMsgAppAlreadyExists)
		}
	}
	err = checkAppInfo(in.GetInfo())
	if err != nil {
		return nil, err
	}
	span.AddEvent("call grpc ExistAppCondsV2")
	exist, err := grpc.ExistAppCondsV2(ctx, &appcrud.Conds{Name: &npool.StringVal{
		Value: in.GetInfo().GetName(),
		Op:    cruder.EQ,
	}})
	if err != nil {
		return &app.CreateAppResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}
	if exist {
		return &app.CreateAppResponse{}, status.Error(npool.ErrAlreadyExists, "app name already exist")
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
