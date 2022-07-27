//nolint:nolintlint,dupl
package admin

import (
	"context"

	bconstant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	mw "github.com/NpoolPlatform/appuser-gateway/pkg/middleware/admin"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/appuser-manager/pkg/db/ent"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	"github.com/NpoolPlatform/message/npool/appusergw/admin"
	appcrud "github.com/NpoolPlatform/message/npool/appusermgrv2/app"
	approlecrud "github.com/NpoolPlatform/message/npool/appusermgrv2/approle"
	"github.com/google/uuid"
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

	span.AddEvent("call grpc GetAppsV2")
	resp, _, err := grpc.GetAppsV2(ctx, &appcrud.Conds{
		IDIn: &npool.StringSlicesVal{
			Value: []string{bconstant.GenesisAppID, bconstant.ChurchAppID},
			Op:    cruder.IN,
		},
	}, 2, 0)
	if err != nil {
		if !ent.IsNotFound(err) {
			logger.Sugar().Errorw("fail get admin apps: %v", err)
			return nil, err
		}
	}

	if len(resp) == 0 {
		logger.Sugar().Errorw("admin app not found : %v", err)
		return nil, status.Error(npool.ErrNotFound, appusergw.ErrMsgAdminAppNotFound)
	}

	return &admin.GetAdminAppsResponse{
		Infos: resp,
	}, nil
}

func (s *Server) GetGenesisRole(ctx context.Context, in *admin.GetGenesisRoleRequest) (*admin.GetGenesisRoleResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetGenesisRole")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.AddEvent("call grpc GetAppRoleOnlyV2")
	resp, err := grpc.GetAppRoleOnlyV2(ctx, &approlecrud.Conds{
		AppID: &npool.StringVal{
			Value: uuid.UUID{}.String(),
			Op:    cruder.EQ,
		},
		Role: &npool.StringVal{
			Value: bconstant.GenesisRole,
			Op:    cruder.EQ,
		},
	})
	if err != nil {
		if ent.IsNotFound(err) {
			logger.Sugar().Errorw("genesis role not found: %v", err)
			return &admin.GetGenesisRoleResponse{}, status.Error(npool.ErrNotFound, appusergw.ErrMsgAdminAppNotFound)
		}
		logger.Sugar().Errorw("fail get genesis role : %v", err)
		return &admin.GetGenesisRoleResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &admin.GetGenesisRoleResponse{
		Info: resp,
	}, nil
}

func (s *Server) GetAppGenesisAppRoleUsers(ctx context.Context,
	in *admin.GetAppGenesisAppRoleUsersRequest) (*admin.GetAppGenesisAppRoleUsersResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppGenesisAppRoleUsers")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Error("TargetAppID is invalid")
		return &admin.GetAppGenesisAppRoleUsersResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgAppIDInvalid)
	}

	span.AddEvent("call middleware GetAppGenesisAppRoleUsers")
	resp, err := mw.GetAppGenesisAppRoleUsers(ctx, in.GetTargetAppID())
	if err != nil {
		logger.Sugar().Errorw("fail get genesis app role user: %v", err)
		return &admin.GetAppGenesisAppRoleUsersResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &admin.GetAppGenesisAppRoleUsersResponse{
		Infos: resp,
	}, nil
}
