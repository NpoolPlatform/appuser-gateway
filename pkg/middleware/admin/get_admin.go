package admin

import (
	"context"
	"fmt"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	serconst "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	approlecrud "github.com/NpoolPlatform/message/npool/appusermgrv2/approle"
	approleusercrud "github.com/NpoolPlatform/message/npool/appusermgrv2/approleuser"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
)

func GetAppGenesisAppRoleUsers(ctx context.Context, appID string) ([]*approleusercrud.AppRoleUser, error) {
	var err error

	_, span := otel.Tracer(serconst.ServiceName).Start(ctx, "GetAppGenesisAppRoleUsers")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.AddEvent("call methods GetGenesisRole")
	role, err := GetGenesisRole(ctx)
	if err != nil {
		logger.Sugar().Errorw("fail get genesis role: %v", err)
		return nil, fmt.Errorf("fail get genesis role: %v", err)
	}

	resp, err := grpc.GetAppRoleUsersV2(ctx, &approleusercrud.Conds{
		AppID: &npool.StringVal{
			Value: appID,
			Op:    cruder.EQ,
		},
		RoleID: &npool.StringVal{
			Value: role.GetID(),
			Op:    cruder.EQ,
		},
	})
	if err != nil {
		logger.Sugar().Errorw("fail get role user: %v", err)
		return nil, err
	}
	return resp, nil
}

func GetGenesisRole(ctx context.Context) (*approlecrud.AppRole, error) {
	var err error

	_, span := otel.Tracer(serconst.ServiceName).Start(ctx, "GetGenesisRole")
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
			Value: constant.GenesisRole,
			Op:    cruder.EQ,
		},
	})
	if err != nil {
		logger.Sugar().Errorw("fail get role: %v", err)
		return nil, err
	}

	return resp, nil
}
