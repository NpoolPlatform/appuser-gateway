//nolint:nolintlint,dupl
package admin

import (
	"context"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	tracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer/admin"
	"google.golang.org/grpc/codes"

	madmin "github.com/NpoolPlatform/appuser-gateway/pkg/admin"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/admin"
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

	span = commontracer.TraceInvoker(span, "admin", "pkg", "CreateAdminApps")

	info, err := madmin.CreateAdminApps(ctx)
	if err != nil {
		logger.Sugar().Errorw("CreateAdminApps", "err", err)
		return &admin.CreateAdminAppsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &admin.CreateAdminAppsResponse{
		Infos: info,
	}, nil
}

func (s *Server) CreateGenesisRoles(ctx context.Context, in *admin.CreateGenesisRolesRequest) (*admin.CreateGenesisRolesResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateGenesisRoles")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "admin", "pkg", "CreateGenesisRoles")

	infos, err := madmin.CreateGenesisRoles(ctx)
	if err != nil {
		logger.Sugar().Errorw("CreateGenesisRoles", "err", err)
		return &admin.CreateGenesisRolesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &admin.CreateGenesisRolesResponse{
		Infos: infos,
	}, nil
}

func (s *Server) CreateGenesisUser(ctx context.Context,
	in *admin.CreateGenesisUserRequest) (*admin.CreateGenesisUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateGenesisUser")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = tracer.Trace(span, in)

	err = validate(ctx, in)
	if err != nil {
		logger.Sugar().Errorw("CreateGenesisUser", "err", err)
		return nil, err
	}

	span = commontracer.TraceInvoker(span, "admin", "middleware", "CreateGenesisUser")

	info, err := madmin.CreateGenesisUser(
		ctx,
		in.GetTargetAppID(),
		in.GetEmailAddress(),
		in.GetPasswordHash(),
	)
	if err != nil {
		logger.Sugar().Errorw("CreateGenesisUser", "err", err)
		return &admin.CreateGenesisUserResponse{}, status.Error(codes.Internal, "fail create genesis user")
	}

	return &admin.CreateGenesisUserResponse{
		Info: info,
	}, nil
}

func (s *Server) AuthorizeGenesis(ctx context.Context, in *admin.AuthorizeGenesisRequest) (*admin.AuthorizeGenesisResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "AuthorizeGenesis")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	infos, total, err := madmin.AuthorizeGenesis(ctx)
	if err != nil {
		logger.Sugar().Errorw("AuthorizeGenesis", "err", err)
		return &admin.AuthorizeGenesisResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &admin.AuthorizeGenesisResponse{
		Infos: infos,
		Total: total,
	}, nil
}