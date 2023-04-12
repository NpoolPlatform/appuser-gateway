//nolint:nolintlint,dupl
package admin

import (
	"context"

	authmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/authing"

	madmin "github.com/NpoolPlatform/appuser-gateway/pkg/admin"
	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/admin"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetAdminApps(ctx context.Context, in *admin.GetAdminAppsRequest) (*admin.GetAdminAppsResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAdminApps")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "admin", "manager", "GetApps")

	infos, err := madmin.GetAdminApps(ctx)
	if err != nil {
		logger.Sugar().Errorw("GetAdminApps", "err", err)
		return &admin.GetAdminAppsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &admin.GetAdminAppsResponse{
		Infos: infos,
	}, nil
}

func (s *Server) GetGenesisRoles(ctx context.Context, in *admin.GetGenesisRolesRequest) (*admin.GetGenesisRolesResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetGenesisRole")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "admin", "manager", "GetGenesisRoles")

	infos, total, err := madmin.GetGenesisRoles(ctx)
	if err != nil {
		logger.Sugar().Errorw("GetGenesisRole", "err", "genesis role not found")
		return &admin.GetGenesisRolesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &admin.GetGenesisRolesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetGenesisUsers(ctx context.Context,
	in *admin.GetGenesisUsersRequest) (*admin.GetGenesisUsersResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetGenesisUsers")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "admin", "pkg", "GetGenesisUsers")

	infos, total, err := madmin.GetGenesisUsers(ctx)
	if err != nil {
		logger.Sugar().Errorw("GetGenesisUsers", "err", err)
		return &admin.GetGenesisUsersResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &admin.GetGenesisUsersResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetGenesisAuths(ctx context.Context, in *admin.GetGenesisAuthsRequest) (*admin.GetGenesisAuthsResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "AuthorizeGenesis")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "AppID", in.GetTargetAppID(), "error", err)
		return &admin.GetGenesisAuthsResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	infos, total, err := authmwcli.GetAuths(ctx, in.GetTargetAppID(), 0, 0)
	if err != nil {
		return nil, err
	}

	return &admin.GetGenesisAuthsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
