package api

import (
	"context"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	bcode "github.com/NpoolPlatform/business-status-code"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergateway/app"
	appcrud "github.com/NpoolPlatform/message/npool/appusermgrv2/app"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AppServer) UpdateApp(ctx context.Context, in *app.UpdateAppRequest) (*app.UpdateAppResponse, error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateAppV2")
	defer span.End()
	var err error
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()
	if in.GetInfo().GetID() == "" {
		return &app.UpdateAppResponse{}, bcode.ErrParam
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
		return &app.UpdateAppResponse{}, status.Error(codes.Internal, err.Error())
	}
	if exist {
		return &app.UpdateAppResponse{}, status.Error(codes.AlreadyExists, "app name already exist")
	}
	span.AddEvent("call grpc CreateAppV2")
	resp, err := grpc.CreateAppV2(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("fail create app: %v", err)
		return &app.UpdateAppResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &app.UpdateAppResponse{
		Info: resp,
	}, nil
}
