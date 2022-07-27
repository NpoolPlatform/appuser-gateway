package admin

import (
	"context"

	"github.com/NpoolPlatform/api-manager/pkg/db/ent"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	serconst "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	appcrud "github.com/NpoolPlatform/message/npool/appusermgrv2/app"
	approlecrud "github.com/NpoolPlatform/message/npool/appusermgrv2/approle"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
)

func CreateAdminApps(ctx context.Context) ([]*appcrud.App, error) {
	var err error
	apps := []*appcrud.App{}
	createApps := []*appcrud.AppReq{}

	_, span := otel.Tracer(serconst.ServiceName).Start(ctx, "CreateAdminApps")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.AddEvent("call grpc GetAppV2")
	genesisApp, err := grpc.GetAppV2(ctx, constant.GenesisAppID)
	if err != nil {
		if !ent.IsNotFound(err) {
			logger.Sugar().Errorw("fail get app: %v", err)
			return nil, err
		}
	}

	if genesisApp != nil {
		apps = append(apps, genesisApp)
	} else {
		genesisApp := appcrud.App{
			Description: "NOT SET",
			ID:          constant.GenesisAppID,
			CreatedBy:   uuid.UUID{}.String(),
			Name:        constant.GenesisAppName,
			Logo:        "NOT SET",
		}
		createApps = append(createApps, &appcrud.AppReq{
			Description: &genesisApp.Description,
			ID:          &genesisApp.ID,
			CreatedBy:   &genesisApp.CreatedBy,
			Name:        &genesisApp.Name,
			Logo:        &genesisApp.Logo,
		})
	}

	span.AddEvent("call grpc GetAppV2")
	churchApp, err := grpc.GetAppV2(ctx, constant.ChurchAppID)
	if err != nil {
		if !ent.IsNotFound(err) {
			logger.Sugar().Errorw("fail get apps: %v", err)
			return nil, err
		}
	}

	if churchApp != nil {
		apps = append(apps, churchApp)
	} else {
		churchApp := appcrud.App{
			Description: "NOT SET",
			ID:          constant.ChurchAppID,
			CreatedBy:   uuid.UUID{}.String(),
			Name:        constant.ChurchAppName,
			Logo:        "NOT SET",
		}
		createApps = append(createApps, &appcrud.AppReq{
			Description: &churchApp.Description,
			ID:          &churchApp.ID,
			CreatedBy:   &churchApp.CreatedBy,
			Name:        &churchApp.Name,
			Logo:        &churchApp.Logo,
		})
	}

	span.AddEvent("call grpc CreateAppsV2")
	resp, err := grpc.CreateAppsV2(ctx, createApps)
	if err != nil {
		logger.Sugar().Errorw("fail create admin apps: %v", err)
		return nil, err
	}
	apps = append(apps, resp...)

	return apps, nil
}

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
		logger.Sugar().Errorw("fail create admin app: %v", err)
		return nil, err
	}

	return resp, nil
}
