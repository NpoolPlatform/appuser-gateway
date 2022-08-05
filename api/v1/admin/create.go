//nolint:nolintlint,dupl
package admin

import (
	"context"
	"fmt"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	tracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer/admin"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"

	mw "github.com/NpoolPlatform/appuser-gateway/pkg/admin"
	constants "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	appusermwadmin "github.com/NpoolPlatform/appuser-middleware/pkg/client/admin"
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

	resp, err := mw.CreateAdminApps(ctx)
	if err != nil {
		logger.Sugar().Errorw("CreateAdminApps", "err", err)
		return &admin.CreateAdminAppsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &admin.CreateAdminAppsResponse{
		Infos: resp,
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

	resp, err := mw.CreateGenesisRoles(ctx)
	if err != nil {
		logger.Sugar().Errorw("CreateGenesisRoles", "err", err)
		return &admin.CreateGenesisRolesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &admin.CreateGenesisRolesResponse{
		Infos: resp,
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

	role := constants.GenesisRole
	if in.GetTargetAppID() == constants.ChurchAppID {
		role = constants.ChurchRole
	}

	span = tracer.Trace(span, in)

	err = validate(ctx, in, role)
	if err != nil {
		logger.Sugar().Errorw("CreateGenesisUser", "err", err)
		return nil, err
	}

	span = commontracer.TraceInvoker(span, "admin", "middleware", "CreateGenesisUser")

	resp, err := appusermwadmin.CreateGenesisUser(
		ctx,
		in.GetTargetAppID(),
		uuid.NewString(),
		role,
		in.GetEmailAddress(),
		in.GetPasswordHash(),
	)
	if err != nil {
		logger.Sugar().Errorw("CreateGenesisUser", "err", err)
		return &admin.CreateGenesisUserResponse{}, status.Error(codes.Internal, "fail create genesis user")
	}

	return &admin.CreateGenesisUserResponse{
		Info: resp,
	}, nil
}

func (s *Server) AuthorizeGenesis(ctx context.Context, in *admin.AuthorizeGenesisRequest) (*admin.AuthorizeGenesisResponse, error) {
	// TODO: Wait for authing-gateway refactoring to complete the API
	return &admin.AuthorizeGenesisResponse{}, status.Error(codes.Internal, fmt.Errorf("NOT IMPLEMENTED").Error())
}
