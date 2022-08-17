package admin

import (
	"context"
	"encoding/json"
	"fmt"

	tracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer/admin"
	appmgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/app"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/admin"
	appmrgpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/app"
	rolemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role"
	"github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	serconst "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	appusermgrconst "github.com/NpoolPlatform/appuser-manager/pkg/message/const"

	approlemgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approle"
	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	approlepb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approle"
	appmw "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
)

func CreateAdminApps(ctx context.Context) ([]*appmw.App, error) {
	var err error
	genesisApps := []*appmrgpb.AppReq{}

	_, span := otel.Tracer(serconst.ServiceName).Start(ctx, "CreateAdminApps")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "admin", "apollo", "GetStringValueWithNameSpace")

	hostname := config.GetStringValueWithNameSpace(appusermgrconst.ServiceName, config.KeyHostname)
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

	appInfos, total, err := appmwcli.GetManyApps(ctx, appIDs)
	if err != nil {
		logger.Sugar().Errorw("CreateAdminApps", "error", err)
		return nil, err
	}

	if total > 0 {
		return appInfos, nil
	}

	span = commontracer.TraceInvoker(span, "admin", "middleware", "CreateApps")

	createBy := uuid.UUID{}.String()
	logo := "NOT SET"
	for key := range genesisApps {
		genesisApps[key].CreatedBy = &createBy
		genesisApps[key].Logo = &logo
	}

	_, err = appmgrcli.CreateApps(ctx, genesisApps)
	if err != nil {
		logger.Sugar().Errorw("CreateAdminApps", "error", err)
		return nil, err
	}

	appInfos, _, err = appmwcli.GetManyApps(ctx, appIDs)
	if err != nil {
		logger.Sugar().Errorw("CreateAdminApps", "error", err)
		return nil, err
	}

	return appInfos, nil
}

func CreateGenesisRoles(ctx context.Context) ([]*rolemwpb.Role, error) {
	var err error
	genesisRoles := []*approlepb.AppRoleReq{}

	_, span := otel.Tracer(serconst.ServiceName).Start(ctx, "CreateGenesisRoles")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "admin", "apollo", "GetStringValueWithNameSpace")

	hostname := config.GetStringValueWithNameSpace(appusermgrconst.ServiceName, config.KeyHostname)
	genesisRoleStr := config.GetStringValueWithNameSpace(hostname, constant.KeyGenesisRole)

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

	respRoles, _, err := approlemgrcli.GetAppRoles(ctx, &approlepb.Conds{
		Roles: &npool.StringSliceVal{
			Op:    cruder.IN,
			Value: roles,
		},
	}, 0, int32(len(roles)))
	if err != nil {
		logger.Sugar().Errorw("CreateGenesisRoles", "error", err)
		return nil, err
	}

	appRoleIDs := []string{}
	for _, val := range respRoles {
		appRoleIDs = append(appRoleIDs, val.GetID())
	}

	if len(respRoles) > 0 {
		appRoles, _, err := rolemwcli.GetManyRoles(ctx, appRoleIDs)
		if err != nil {
			logger.Sugar().Errorw("CreateGenesisRoles", "error", err)
			return nil, err
		}

		if len(appRoles) > 0 {
			return appRoles, nil
		}
	}

	span = commontracer.TraceInvoker(span, "admin", "middleware", "CreateAppRoles")

	for key := range genesisRoles {
		defaultVal := false
		genesisRoles[key].Default = &defaultVal

		createBy := uuid.NewString()
		genesisRoles[key].CreatedBy = &createBy

		description := "NOT SET"
		genesisRoles[key].Description = &description
	}

	respRoles, err = approlemgrcli.CreateAppRoles(ctx, genesisRoles)
	if err != nil {
		logger.Sugar().Errorw("CreateGenesisRoles", "error", err)
		return nil, err
	}

	for _, val := range respRoles {
		appRoleIDs = append(appRoleIDs, val.GetID())
	}

	appRoles, _, err := rolemwcli.GetManyRoles(ctx, appRoleIDs)
	if err != nil {
		logger.Sugar().Errorw("CreateGenesisRoles", "error", err)
		return nil, err
	}

	return appRoles, nil
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
	if err != nil {
		logger.Sugar().Errorw("CreateGenesisUser", "error", err, "appRole", appRole)
		return nil, err
	}
	if appRole == nil {
		return nil, fmt.Errorf("fail  get app role")
	}

	roleID := appRole.ID
	userID := uuid.NewString()

	userInfo, err := usermwcli.CreateUser(ctx, &user.UserReq{
		ID:           &userID,
		AppID:        &appID,
		EmailAddress: &emailAddress,
		PasswordHash: &passwordHash,
		RoleIDs:      []string{roleID},
	})
	if err != nil {
		logger.Sugar().Errorw("CreateGenesisUser", "error", err)
		return nil, err
	}

	return userInfo, nil
}
