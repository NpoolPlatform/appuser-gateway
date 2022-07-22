package role

import (
	"context"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw/approle"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateRole(ctx context.Context, in *approle.CreateRoleRequest) (*approle.CreateRoleResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateRole")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	err = validate(in.GetInfo(), in.GetUserID())
	if err != nil {
		return nil, err
	}

	span.AddEvent("call grpc CreateAppRoleV2")
	resp, err := grpc.CreateAppRoleV2(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("fail create app role: %v", err)
		return &approle.CreateRoleResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &approle.CreateRoleResponse{
		Info: resp,
	}, nil
}

func (s *Server) CreateAppRole(ctx context.Context, in *approle.CreateAppRoleRequest) (*approle.CreateAppRoleResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateAppRole")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	err = validate(in.GetInfo(), in.GetUserID())
	if err != nil {
		return nil, err
	}

	appID := in.GetTargetAppID()
	info := in.GetInfo()
	info.AppID = &appID

	span.AddEvent("call grpc CreateAppRoleV2")
	resp, err := grpc.CreateAppRoleV2(ctx, info)
	if err != nil {
		logger.Sugar().Errorw("fail create app role: %v", err)
		return &approle.CreateAppRoleResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &approle.CreateAppRoleResponse{
		Info: resp,
	}, nil
}
