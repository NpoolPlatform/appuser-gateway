//nolint:nolintlint,dupl
package api

import (
	"context"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	ausercrudgrpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	appmwgrpc "github.com/NpoolPlatform/appuser-middleware/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergateway/app"
	appcrud "github.com/NpoolPlatform/message/npool/appusermgrv2/app"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *AppServer) GetApp(ctx context.Context, in *app.GetAppRequest) (*app.GetAppResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateAppV2")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Error("ID is invalid")
		return &app.GetAppResponse{}, status.Error(npool.ErrParams, app.ErrMsgAppIDInvalid)
	}

	span.AddEvent("call grpc GetAppV2")
	resp, err := ausercrudgrpc.GetAppV2(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("fail get app: %v", err)
		return &app.GetAppResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}
	return &app.GetAppResponse{
		Info: resp,
	}, nil
}

func (s *AppServer) GetApps(ctx context.Context, in *app.GetAppsRequest) (*app.GetAppsResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateAppV2")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.AddEvent("call grpc GetAppsV2")
	resp, total, err := ausercrudgrpc.GetAppsV2(ctx, &appcrud.Conds{}, in.GetLimit(), in.GetOffset())
	if err != nil {
		logger.Sugar().Errorw("fail get apps: %v", err)
		return &app.GetAppsResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &app.GetAppsResponse{
		Infos: resp,
		Total: total,
	}, nil
}

func (s *AppServer) GetAppsByCreator(ctx context.Context, in *app.GetAppsByCreatorRequest) (*app.GetAppsByCreatorResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateAppV2")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.AddEvent("call grpc GetAppsV2")
	resp, total, err := ausercrudgrpc.GetAppsV2(ctx, &appcrud.Conds{
		CreatedBy: &npool.StringVal{
			Value: in.GetUserID(),
			Op:    cruder.EQ,
		},
	}, in.GetLimit(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("fail get apps: %v", err)
		return &app.GetAppsByCreatorResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &app.GetAppsByCreatorResponse{
		Infos: resp,
		Total: total,
	}, nil
}

func (s *AppServer) GetAppInfo(ctx context.Context, in *app.GetAppInfoRequest) (*app.GetAppInfoResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppInfo")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Error("ID is invalid")
		return &app.GetAppInfoResponse{}, status.Error(npool.ErrParams, app.ErrMsgAppIDInvalid)
	}

	span.AddEvent("call grpc GetAppInfo")
	resp, err := appmwgrpc.GetAppInfo(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("fail get app info: %v", err)
		return &app.GetAppInfoResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &app.GetAppInfoResponse{
		Info: resp,
	}, nil
}
