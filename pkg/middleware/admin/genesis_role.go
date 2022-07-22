package admin

import (
	"context"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	serconst "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	approlecrud "github.com/NpoolPlatform/message/npool/appusermgrv2/approle"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
)

func CreateGenesisRole(ctx context.Context) (*approlecrud.AppRole, error) {
	var err error

	_, span := otel.Tracer(serconst.ServiceName).Start(ctx, "CreateGenesisRole")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	genesisRole := approlecrud.AppRole{
		AppID:       uuid.UUID{}.String(),
		CreatedBy:   uuid.UUID{}.String(),
		Role:        constant.GenesisRole,
		Description: "NOT SET",
		Default:     false,
	}

	span.AddEvent("call grpc CreateAppRoleV2")
	resp, err := grpc.CreateAppRoleV2(ctx, &approlecrud.AppRoleReq{
		AppID:       &genesisRole.AppID,
		CreatedBy:   &genesisRole.CreatedBy,
		Role:        &genesisRole.Role,
		Description: &genesisRole.Description,
		Default:     &genesisRole.Default,
	})
	if err != nil {
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

	span.AddEvent("call grpc CreateAppRoleV2")
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
		return nil, err
	}

	return resp, nil
}
