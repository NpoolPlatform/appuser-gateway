//nolint:nolintlint,dupl
package role

import (
	"context"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	"github.com/NpoolPlatform/message/npool/appusergw/approle"
	approlecrud "github.com/NpoolPlatform/message/npool/appusermgrv2/approle"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetRole(ctx context.Context, in *approle.GetRoleRequest) (*approle.GetRoleResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetRole")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Error("ID is invalid")
		return &approle.GetRoleResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgIDInvalid)
	}

	span.AddEvent("call grpc GetBanAppV2")
	resp, err := grpc.GetAppRoleV2(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Error("fail get app role:%v", err)
		return &approle.GetRoleResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &approle.GetRoleResponse{
		Info: resp,
	}, nil
}

func (s *Server) GetAppRoleByRole(ctx context.Context, in *approle.GetAppRoleByRoleRequest) (*approle.GetAppRoleByRoleResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppRoleByRole")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if in.GetRole() == "" {
		logger.Sugar().Error("Role empty")
		return &approle.GetAppRoleByRoleResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgRoleEmpty)
	}

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Error("AppID is invalid")
		return &approle.GetAppRoleByRoleResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgAppIDInvalid)
	}

	span.AddEvent("call grpc GetAppRoleOnlyV2")
	resp, err := grpc.GetAppRoleOnlyV2(ctx, &approlecrud.Conds{
		AppID: &npool.StringVal{
			Value: in.GetTargetAppID(),
			Op:    cruder.EQ,
		},
		Role: &npool.StringVal{
			Value: in.GetRole(),
			Op:    cruder.EQ,
		},
	})
	if err != nil {
		logger.Sugar().Error("fail get app role:%v", err)
		return &approle.GetAppRoleByRoleResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &approle.GetAppRoleByRoleResponse{
		Info: resp,
	}, nil
}

func (s *Server) GetRoles(ctx context.Context, in *approle.GetRolesRequest) (*approle.GetRolesResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetRoles")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Error("AppID is invalid")
		return &approle.GetRolesResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgAppIDInvalid)
	}

	span.AddEvent("call grpc GetAppRolesV2")
	resp, _, err := grpc.GetAppRolesV2(ctx, &approlecrud.Conds{
		AppID: &npool.StringVal{
			Value: in.GetAppID(),
			Op:    cruder.EQ,
		},
	}, in.GetLimit(), in.GetOffset())
	if err != nil {
		logger.Sugar().Error("fail get app roles:%v", err)
		return &approle.GetRolesResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &approle.GetRolesResponse{
		Infos: resp,
	}, nil
}

func (s *Server) GetAppRoles(ctx context.Context, in *approle.GetAppRolesRequest) (*approle.GetAppRolesResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppRoles")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Error("AppID is invalid")
		return &approle.GetAppRolesResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgAppIDInvalid)
	}

	span.AddEvent("call grpc GetAppRolesV2")
	resp, _, err := grpc.GetAppRolesV2(ctx, &approlecrud.Conds{
		AppID: &npool.StringVal{
			Value: in.GetTargetAppID(),
			Op:    cruder.EQ,
		},
	}, in.GetLimit(), in.GetOffset())
	if err != nil {
		logger.Sugar().Error("fail get app roles:%v", err)
		return &approle.GetAppRolesResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &approle.GetAppRolesResponse{
		Infos: resp,
	}, nil
}
