package admin

import (
	"context"

	tracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer/admin"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/admin"
	"github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	serconst "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"

	appmgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/app"
	approlemgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approle"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	appcrud "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/app"
	approlepb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approle"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
)

func CreateAdminApps(ctx context.Context) ([]*appcrud.App, error) {
	var err error
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

	appInfos, total, err := appmgrcli.GetApps(ctx, &appcrud.Conds{
		IDs: &npool.StringSliceVal{
			Op:    cruder.IN,
			Value: []string{constant.GenesisAppID, constant.ChurchAppID},
		},
	}, 0, 2) // nolint

	if err != nil {
		logger.Sugar().Errorw("CreateAdminApps", "error", err)
		return nil, err
	}

	if total == 2 { // nolint
		return appInfos, nil
	}

	createGenesis := true
	createChurch := true
	for _, val := range appInfos {
		if val.ID == constant.GenesisAppID {
			createGenesis = false
		}
		if val.ID == constant.ChurchAppID {
			createChurch = false
		}
	}

	if createGenesis {
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

	if createChurch {
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

	span = commontracer.TraceInvoker(span, "admin", "middleware", "CreateApps")

	resp, err := appmgrcli.CreateApps(ctx, createApps)
	if err != nil {
		logger.Sugar().Errorw("CreateAdminApps", "error", err)
		return nil, err
	}

	return resp, nil
}

func CreateGenesisRoles(ctx context.Context) ([]*approlepb.AppRole, error) {
	var err error

	_, span := otel.Tracer(serconst.ServiceName).Start(ctx, "CreateGenesisRole")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	appRoles, total, err := approlemgrcli.GetAppRoles(ctx, &approlepb.Conds{
		Roles: &npool.StringSliceVal{
			Op:    cruder.IN,
			Value: []string{constant.GenesisRole, constant.ChurchRole},
		},
	}, 0, 2) // nolint
	if err != nil {
		logger.Sugar().Errorw("CreateGenesisRole", "error", err)
		return nil, err
	}

	if total == 2 { // nolint
		return appRoles, nil
	}

	createGenesis := true
	createChurch := true
	for _, val := range appRoles {
		if val.AppID == constant.GenesisAppID {
			createGenesis = false
		}
		if val.AppID == constant.ChurchAppID {
			createChurch = false
		}
	}

	createAppRoles := []*approlepb.AppRoleReq{}

	if createGenesis {
		genesisRole := approlepb.AppRole{
			AppID:       constant.GenesisAppID,
			CreatedBy:   uuid.UUID{}.String(),
			Role:        constant.GenesisRole,
			Description: "NOT SET",
			Default:     false,
		}

		createAppRoles = append(createAppRoles, &approlepb.AppRoleReq{
			AppID:       &genesisRole.AppID,
			CreatedBy:   &genesisRole.CreatedBy,
			Role:        &genesisRole.Role,
			Description: &genesisRole.Description,
			Default:     &genesisRole.Default,
		})
	}

	if createChurch {
		churchRole := approlepb.AppRole{
			AppID:       constant.ChurchAppID,
			CreatedBy:   uuid.UUID{}.String(),
			Role:        constant.ChurchRole,
			Description: "NOT SET",
			Default:     false,
		}

		createAppRoles = append(createAppRoles, &approlepb.AppRoleReq{
			AppID:       &churchRole.AppID,
			CreatedBy:   &churchRole.CreatedBy,
			Role:        &churchRole.Role,
			Description: &churchRole.Description,
			Default:     &churchRole.Default,
		})
	}

	span = commontracer.TraceInvoker(span, "admin", "middleware", "CreateAppRoles")

	resp, err := approlemgrcli.CreateAppRoles(ctx, createAppRoles)
	if err != nil {
		logger.Sugar().Errorw("CreateGenesisRole", "error", err)
		return nil, err
	}

	return resp, nil
}

func CreateGenesisUser(ctx context.Context, appID, emailAddress, passwordHash string) (*user.User, error) {
	var err error

	_, span := otel.Tracer(serconst.ServiceName).Start(ctx, "CreateGenesisUser")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = tracer.Trace(span, &admin.CreateGenesisUserRequest{
		TargetAppID:  appID,
		EmailAddress: emailAddress,
		PasswordHash: passwordHash,
	})

	span = commontracer.TraceInvoker(span, "app", "db", "CreateGenesisUser")

	appRole, err := approlemgrcli.GetAppRoleOnly(ctx, &approlepb.Conds{
		AppID: &npool.StringVal{
			Value: appID,
			Op:    cruder.EQ,
		},
	})
	if err != nil || appRole == nil {
		logger.Sugar().Errorw("CreateGenesisUser", "error", err, "appRole", appRole)
		return nil, err
	}

	roleID := appRole.ID
	importFromApp := uuid.UUID{}.String()
	userID := uuid.NewString()

	userInfo, err := usermwcli.CreateUser(ctx, &user.UserReq{
		ID:                &userID,
		AppID:             &appID,
		EmailAddress:      &emailAddress,
		ImportedFromAppID: &importFromApp,
		PasswordHash:      &passwordHash,
		RoleIDs:           []string{roleID},
	})
	if err != nil {
		logger.Sugar().Errorw("CreateGenesisUser", "error", err)
		return nil, err
	}

	return userInfo, nil
}
