package app

import (
	"context"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw/app"
	appcrud "github.com/NpoolPlatform/message/npool/appusermgrv2/app"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateApp(ctx context.Context, in *app.UpdateAppRequest) (*app.UpdateAppResponse, error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateAppV2")
	defer span.End()
	var err error
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()
	if _, err := uuid.Parse(in.GetInfo().GetID()); err != nil {
		return &app.UpdateAppResponse{}, status.Error(npool.ErrParams, app.ErrMsgAppIDInvalid)
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
		return &app.UpdateAppResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}
	if exist {
		return &app.UpdateAppResponse{}, status.Error(npool.ErrAlreadyExists, app.ErrMsgAppNameAlreadyExists)
	}

	span.AddEvent("call grpc CreateAppV2")
	resp, err := grpc.UpdateAppV2(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("fail create app: %v", err)
		return &app.UpdateAppResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &app.UpdateAppResponse{
		Info: resp,
	}, nil
}
