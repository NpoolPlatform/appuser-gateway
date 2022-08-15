package admin

import (
	"context"
	"encoding/json"
	"fmt"
	tracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer/admin"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/go-service-framework/pkg/config"
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
	genesisApps := []*appcrud.AppReq{}

	_, span := otel.Tracer(serconst.ServiceName).Start(ctx, "CreateAdminApps")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "admin", "apollo", "GetStringValueWithNameSpace")

	hostname := config.GetStringValueWithNameSpace("", config.KeyHostname)
	genesisAppStr := config.GetStringValueWithNameSpace(hostname, constant.KeyGenesisApp)

	err = json.Unmarshal([]byte(genesisAppStr), &genesisApps)
	if err != nil {
		return nil, err
	}

	if len(genesisApps) == 0 {
		return nil, fmt.Errorf("genesis app not available")
	}

	span = commontracer.TraceInvoker(span, "admin", "middleware", "GetApps")

	appIDs := []string{}
	for key := range genesisApps {
		appIDs = append(appIDs, genesisApps[key].GetID())
	}

	appInfos, total, err := appmgrcli.GetApps(ctx, &appcrud.Conds{
		IDs: &npool.StringSliceVal{
			Value: appIDs,
			Op:    cruder.IN,
		},
	}, 0, int32(len(appIDs)))

	if err != nil {
		logger.Sugar().Errorw("CreateAdminApps", "error", err)
		return nil, err
	}

	if total > 0 {
		return appInfos, nil
	}

	span = commontracer.TraceInvoker(span, "admin", "middleware", "CreateApps")

	resp, err := appmgrcli.CreateApps(ctx, genesisApps)
	if err != nil {
		logger.Sugar().Errorw("CreateAdminApps", "error", err)
		return nil, err
	}

	return resp, nil
}

func CreateGenesisRoles(ctx context.Context) ([]*approlepb.AppRole, error) {
	var err error
	genesisRoles := []*approlepb.AppRoleReq{}

	_, span := otel.Tracer(serconst.ServiceName).Start(ctx, "CreateGenesisRole")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "admin", "apollo", "GetStringValueWithNameSpace")

	hostname := config.GetStringValueWithNameSpace("", config.KeyHostname)
	genesisRoleStr := config.GetStringValueWithNameSpace(hostname, constant.KeyGenesisApp)

	err = json.Unmarshal([]byte(genesisRoleStr), &genesisRoles)
	if err != nil {
		return nil, err
	}

	if len(genesisRoles) == 0 {
		return nil, fmt.Errorf("genesis role not available")
	}

	roles := []string{}
	for key := range genesisRoles {
		roles = append(roles, genesisRoles[key].GetRole())
	}

	appRoles, total, err := approlemgrcli.GetAppRoles(ctx, &approlepb.Conds{
		Roles: &npool.StringSliceVal{
			Op:    cruder.IN,
			Value: roles,
		},
	}, 0, int32(len(roles)))
	if err != nil {
		logger.Sugar().Errorw("CreateGenesisRole", "error", err)
		return nil, err
	}

	if total > 0 {
		return appRoles, nil
	}

	span = commontracer.TraceInvoker(span, "admin", "middleware", "CreateAppRoles")

	for key := range genesisRoles {
		genesisRoles[key].AppID = genesisRoles[key].ID

		id := uuid.NewString()
		genesisRoles[key].ID = &id

		defaultVal := false
		genesisRoles[key].Default = &defaultVal
	}
	resp, err := approlemgrcli.CreateAppRoles(ctx, genesisRoles)
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
