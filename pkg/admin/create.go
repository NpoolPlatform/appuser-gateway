package admin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	tracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer/admin"
	appmgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/app"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	appmrgpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/app"
	rolemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role"

	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/admin"
	"github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	servicename "github.com/NpoolPlatform/appuser-gateway/pkg/servicename"
	servicename2 "github.com/NpoolPlatform/appuser-manager/pkg/servicename"

	approlemgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approle"
	roleusermgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approleuser"
	authmgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/authing/auth"
	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	authmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/authing"
	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role"
	roleusermgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approleuser"

	approlepb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approle"
	authmgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/authing/auth"
	appmw "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
	authingmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/authing"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
)

func GetGenesisApps() ([]*appmrgpb.AppReq, error) {
	genesisApps := []*appmrgpb.AppReq{}
	genesisAppStr := config.GetStringValueWithNameSpace(servicename2.ServiceDomain, constant.KeyGenesisApp)

	err := json.Unmarshal([]byte(genesisAppStr), &genesisApps)
	if err != nil {
		return nil, err
	}

	if len(genesisApps) == 0 {
		return nil, fmt.Errorf("genesis app not available")
	}
	return genesisApps, nil
}

func CreateAdminApps(ctx context.Context) ([]*appmw.App, error) {
	var err error

	_, span := otel.Tracer(servicename.ServiceDomain).Start(ctx, "CreateAdminApps")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "admin", "apollo", "GetStringValueWithNameSpace")

	genesisApps, err := GetGenesisApps()
	if err != nil {
		return nil, err
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

	_, span := otel.Tracer(servicename.ServiceDomain).Start(ctx, "CreateGenesisRoles")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "admin", "apollo", "GetStringValueWithNameSpace")

	genesisRoleStr := config.GetStringValueWithNameSpace(servicename2.ServiceDomain, constant.KeyGenesisRole)

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

		genesis := true
		genesisRoles[key].Genesis = &genesis
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

	_, span := otel.Tracer(servicename.ServiceDomain).Start(ctx, "CreateGenesisUser")
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

	roles, _, err := approlemgrcli.GetAppRoles(ctx, &approlepb.Conds{
		AppID: &npool.StringVal{
			Value: appID,
			Op:    cruder.EQ,
		},
	}, 0, 100) // nolint
	if err != nil {
		logger.Sugar().Errorw("CreateGenesisUser", "error", err)
		return nil, err
	}
	if len(roles) == 0 {
		return nil, fmt.Errorf("fail get app role")
	}

	userID := uuid.NewString()
	roleIDs := []string{}

	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
	}

	userInfo, err := usermwcli.CreateUser(ctx, &user.UserReq{
		ID:           &userID,
		AppID:        &appID,
		EmailAddress: &emailAddress,
		PasswordHash: &passwordHash,
		RoleIDs:      roleIDs,
	})
	if err != nil {
		logger.Sugar().Errorw("CreateGenesisUser", "error", err)
		return nil, err
	}

	return userInfo, nil
}

type genesisURL struct {
	Path   string
	Method string
}

func createGenesisAuths(ctx context.Context, appID string) (infos []*authingmwpb.Auth, total uint32, err error) {
	_, span := otel.Tracer(servicename.ServiceDomain).Start(ctx, "createGenesisAuths")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	apisJSON := config.GetStringValueWithNameSpace(servicename2.ServiceDomain, constant.KeyGenesisAuthingAPIs)
	apis := []genesisURL{}
	err = json.Unmarshal([]byte(apisJSON), &apis)
	if err != nil {
		return nil, 0, err
	}
	if len(apis) == 0 {
		return nil, 0, fmt.Errorf("genesis authing apis not available")
	}

	span = commontracer.TraceInvoker(span, "admin", "manager", "GetAppRoleUsers")

	roleUsers, _, err := roleusermgrcli.GetAppRoleUsers(ctx, &roleusermgrpb.Conds{
		AppID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
	}, 0, 0)

	auths := []*authmgrpb.AuthReq{}

	for _, val := range roleUsers {
		for _, _api := range apis {
			api := _api
			auths = append(auths, &authmgrpb.AuthReq{
				AppID:    &val.AppID,
				Resource: &api.Path,
				Method:   &api.Method,
				UserID:   &val.UserID,
				RoleID:   &val.RoleID,
			})
		}
	}

	span = commontracer.TraceInvoker(span, "admin", "manager", "CreateAuths")

	_, err = authmgrcli.CreateAuths(ctx, auths)
	if err != nil {
		return nil, 0, err
	}

	return authmwcli.GetAuths(ctx, appID, 0, 0)
}

func AuthorizeGenesis(ctx context.Context) (infos []*authingmwpb.Auth, total uint32, err error) {
	_, span := otel.Tracer(servicename.ServiceDomain).Start(ctx, "AuthorizeGenesis")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "admin", "admin", "getGenesisApps")

	genesisApps, err := GetGenesisApps()
	if err != nil {
		return nil, 0, err
	}

	for _, val := range genesisApps {
		span = commontracer.TraceInvoker(span, "admin", "admin", "createGenesisAuths")
		authInfos, _, err := createGenesisAuths(ctx, val.GetID())
		if err != nil {
			return nil, 0, err
		}

		infos = append(infos, authInfos...)
	}

	return infos, uint32(len(infos)), nil
}
