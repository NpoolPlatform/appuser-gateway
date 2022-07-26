package admin

import (
	"context"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	mw "github.com/NpoolPlatform/appuser-gateway/pkg/middleware/admin"
	appusermw "github.com/NpoolPlatform/appuser-middleware/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw/admin"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateAdminApps(ctx context.Context, in *admin.CreateAdminAppsRequest) (*admin.CreateAdminAppsResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateAdminApps")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.AddEvent("call middleware CreateAdminApps")
	resp, err := mw.CreateAdminApps(ctx)
	if err != nil {
		logger.Sugar().Errorw("fail create admin apps: %v", err)
		return &admin.CreateAdminAppsResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &admin.CreateAdminAppsResponse{
		Infos: resp,
	}, nil
}

func (s *Server) CreateGenesisRole(ctx context.Context, in *admin.CreateGenesisRoleRequest) (*admin.CreateGenesisRoleResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateGenesisRole")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.AddEvent("call middleware CreateGenesisRole")
	resp, err := mw.CreateGenesisRole(ctx)
	if err != nil {
		logger.Sugar().Errorw("fail create genesis role : %v", err)
		return &admin.CreateGenesisRoleResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &admin.CreateGenesisRoleResponse{
		Info: resp,
	}, nil
}

func (s *Server) CreateGenesisRoleUser(ctx context.Context, in *admin.CreateGenesisRoleUserRequest) (*admin.CreateGenesisRoleUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateGenesisRole")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.AddEvent("call middleware CreateGenesisRoleUser")
	resp, err := appusermw.CreateGenesisRoleUser(ctx, in.GetUser(), in.GetSecret())
	if err != nil {
		logger.Sugar().Errorw("fail create genesis role : %v", err)
		return &admin.CreateGenesisRoleUserResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &admin.CreateGenesisRoleUserResponse{
		User:     resp.GetUser(),
		RoleUser: resp.GetRoleUser(),
	}, nil
}
