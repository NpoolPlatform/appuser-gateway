package admin

import (
	"context"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	serconst "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"

	appmgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/app"
	approleapp "github.com/NpoolPlatform/appuser-manager/pkg/client/approle"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	appcrud "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/app"
	approlecrud "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approle"
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

	span = commontracer.TraceInvoker(span, "admin", "middleware", "GetApps")

	genesisApp, _, err := appmgrcli.GetApps(ctx, &appcrud.Conds{
		ID: &npool.StringVal{
			Value: constant.GenesisAppID,
			Op:    cruder.EQ,
		},
	}, 1, 0)

	if err != nil {
		logger.Sugar().Errorw("CreateAdminApps", "error", err)
		return nil, err
	}

	apps = append(apps, genesisApp...)

	if len(genesisApp) > 0 {
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

	span = commontracer.TraceInvoker(span, "admin", "middleware", "GetApps")

	churchApp, _, err := appmgrcli.GetApps(ctx, &appcrud.Conds{
		ID: &npool.StringVal{
			Value: constant.ChurchAppID,
			Op:    cruder.EQ,
		},
	}, 1, 0)
	if err != nil {
		logger.Sugar().Errorw("CreateAdminApps", "error", err)
		return nil, err
	}

	apps = append(apps, churchApp...)

	if len(churchApp) > 0 {
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

	if len(createApps) > 0 {
		span = commontracer.TraceInvoker(span, "admin", "middleware", "CreateApps")

		resp, err := appmgrcli.CreateApps(ctx, createApps)
		if err != nil {
			logger.Sugar().Errorw("CreateAdminApps", "error", err)
			return nil, err
		}
		apps = append(apps, resp...)
	}

	return apps, nil
}

func CreateGenesisRoles(ctx context.Context) ([]*approlecrud.AppRole, error) {
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

	churchRole := approlecrud.AppRole{
		AppID:       uuid.UUID{}.String(),
		CreatedBy:   uuid.UUID{}.String(),
		Role:        constant.ChurchRole,
		Description: "NOT SET",
		Default:     false,
	}

	span = commontracer.TraceInvoker(span, "admin", "middleware", "CreateAppRoles")

	resp, err := approleapp.CreateAppRoles(ctx, []*approlecrud.AppRoleReq{
		{
			AppID:       &genesisRole.AppID,
			CreatedBy:   &genesisRole.CreatedBy,
			Role:        &genesisRole.Role,
			Description: &genesisRole.Description,
			Default:     &genesisRole.Default,
		}, {
			AppID:       &churchRole.AppID,
			CreatedBy:   &churchRole.CreatedBy,
			Role:        &churchRole.Role,
			Description: &churchRole.Description,
			Default:     &churchRole.Default,
		},
	})
	if err != nil {
		logger.Sugar().Errorw("CreateGenesisRole", "error", err)
		return nil, err
	}

	return resp, nil
}
