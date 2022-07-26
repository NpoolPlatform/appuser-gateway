package admin

import (
	"context"
	"github.com/NpoolPlatform/api-manager/pkg/db/ent"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	serconst "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	appcrud "github.com/NpoolPlatform/message/npool/appusermgrv2/app"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
)

func CreateAdminApps(ctx context.Context) ([]*appcrud.App, error) {
	var err error
	apps := []*appcrud.App{}
	createApps := []*appcrud.AppReq{}

	_, span := otel.Tracer(serconst.ServiceName).Start(ctx, "CreateExtra")
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
			logger.Sugar().Errorw("fail get admin app: %v", err)
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
			logger.Sugar().Errorw("fail get admin apps: %v", err)
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

func GetAdminApps(ctx context.Context) ([]*appcrud.App, error) {
	var err error

	_, span := otel.Tracer(serconst.ServiceName).Start(ctx, "CreateExtra")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.AddEvent("call grpc GetAppV2")
	resp, _, err := grpc.GetAppsV2(ctx, &appcrud.Conds{
		IDIn: &npool.StringSlicesVal{
			Value: []string{constant.GenesisAppID, constant.ChurchAppID},
			Op:    cruder.IN,
		},
	}, 2, 0)
	if err != nil {
		if !ent.IsNotFound(err) {
			logger.Sugar().Errorw("fail get admin apps: %v", err)
			return nil, err
		}
	}

	return resp, nil
}
