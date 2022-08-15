//nolint:nolintlint,dupl
package admin

import (
	"context"

	mw "github.com/NpoolPlatform/appuser-gateway/pkg/admin"
	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	authcli "github.com/NpoolPlatform/authing-gateway/pkg/client"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"

	bconstant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	appusermgrapp "github.com/NpoolPlatform/appuser-manager/pkg/client/app"
	appusermgrapprole "github.com/NpoolPlatform/appuser-manager/pkg/client/approle"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/admin"
	appcrud "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/app"
	approlecrud "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approle"
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

	resp, _, err := appusermgrapp.GetApps(ctx, &appcrud.Conds{
		IDs: &npool.StringSliceVal{
			Value: []string{bconstant.GenesisAppID, bconstant.ChurchAppID},
			Op:    cruder.IN,
		},
	}, 0, 2) // nolint
	if err != nil {
		logger.Sugar().Errorw("GetAdminApps", "err", err)
		return &admin.GetAdminAppsResponse{}, status.Error(codes.Internal, err.Error())
	}

	if len(resp) == 0 {
		logger.Sugar().Errorw("GetAdminApps", "err", "admin app no found")
		return nil, status.Error(codes.NotFound, "admin app no found")
	}

	return &admin.GetAdminAppsResponse{
		Infos: resp,
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

	span = commontracer.TraceInvoker(span, "admin", "manager", "GetAppRoles")

	resp, _, err := appusermgrapprole.GetAppRoles(ctx, &approlecrud.Conds{
		Roles: &npool.StringSliceVal{
			Value: []string{bconstant.GenesisRole, bconstant.ChurchRole},
			Op:    cruder.EQ,
		},
	}, 0, 2) // nolint
	if err != nil {
		logger.Sugar().Errorw("GetGenesisRole", "err", err)
		return &admin.GetGenesisRolesResponse{}, status.Error(codes.Internal, err.Error())
	}

	if len(resp) == 0 {
		logger.Sugar().Errorw("GetGenesisRole", "err", "genesis role not found")
		return &admin.GetGenesisRolesResponse{}, status.Error(codes.NotFound, "genesis role not found")
	}

	return &admin.GetGenesisRolesResponse{
		Infos: resp,
	}, nil
}

func (s *Server) GetGenesisRoleUsers(ctx context.Context,
	in *admin.GetGenesisRoleUsersRequest) (*admin.GetGenesisRoleUsersResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetGenesisRoleUsers")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "admin", "pkg", "GetGenesisRoleUsers")

	resp, err := mw.GetAppGenesisAppRoleUsers(ctx)
	if err != nil {
		logger.Sugar().Errorw("GetGenesisRoleUsers", "err", err)
		return &admin.GetGenesisRoleUsersResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &admin.GetGenesisRoleUsersResponse{
		Infos: resp,
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

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "AppID", in.GetAppID(), "error", err)
		return &admin.GetGenesisAuthsResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	infos, err := authcli.GetAuthsByOtherApp(ctx, in.GetAppID())
	if err != nil {
		logger.Sugar().Errorw("GetGenesisAuths", "err", err)
		return &admin.GetGenesisAuthsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &admin.GetGenesisAuthsResponse{
		Infos: infos,
	}, nil
}
