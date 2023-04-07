package admin

import (
	"context"
	"encoding/json"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	constant1 "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	servicename "github.com/NpoolPlatform/appuser-gateway/pkg/servicename"
	servicename2 "github.com/NpoolPlatform/appuser-manager/pkg/servicename"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"

	approlemgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approle"
	approlemgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approle"

	approleusermgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approleuser"
	approleusermgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approleuser"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	appcrud "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/app"
	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"

	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role"
	rolemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
)

func GetAdminApps(ctx context.Context) ([]*appmwpb.App, error) {
	var err error
	genesisApps := []*appcrud.App{}

	_, span := otel.Tracer(servicename.ServiceDomain).Start(ctx, "GetAdminApps")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "admin", "manager", "GetApps")

	genesisAppStr := config.GetStringValueWithNameSpace(servicename2.ServiceDomain, constant1.KeyGenesisApp)

	err = json.Unmarshal([]byte(genesisAppStr), &genesisApps)
	if err != nil {
		return nil, err
	}

	appIDs := []string{}
	for key := range genesisApps {
		appIDs = append(appIDs, genesisApps[key].GetID())
	}

	apps, _, err := appmwcli.GetManyApps(ctx, appIDs)
	if err != nil {
		return nil, err
	}

	return apps, nil
}

func GetGenesisRoles(ctx context.Context) ([]*rolemwpb.Role, uint32, error) {
	var err error
	appRoles := []*approlemgrpb.AppRoleReq{}

	_, span := otel.Tracer(servicename.ServiceDomain).Start(ctx, "GetGenesisRoles")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "admin", "manager", "GetGenesisRoles")

	genesisRoleStr := config.GetStringValueWithNameSpace(servicename2.ServiceDomain, constant1.KeyGenesisRole)

	err = json.Unmarshal([]byte(genesisRoleStr), &appRoles)
	if err != nil {
		return nil, 0, err
	}

	roles := []string{}
	for _, val := range appRoles {
		roles = append(roles, val.GetRole())
	}

	roleInfos, _, err := approlemgrcli.GetAppRoles(ctx, &approlemgrpb.Conds{
		Roles: &commonpb.StringSliceVal{
			Value: roles,
			Op:    cruder.IN,
		},
	}, 0, int32(len(roles)))
	if err != nil {
		logger.Sugar().Errorw("GetGenesisRole", "err", err)
		return nil, 0, err
	}
	if len(roleInfos) == 0 {
		return []*rolemwpb.Role{}, 0, nil
	}

	roleIDs := []string{}
	for _, val := range roleInfos {
		roleIDs = append(roleIDs, val.GetID())
	}

	infos, total, err := rolemwcli.GetManyRoles(ctx, roleIDs)
	if err != nil {
		logger.Sugar().Errorw("GetGenesisRoles", "error", err)
		return nil, 0, err
	}

	return infos, total, nil
}

func GetGenesisUsers(ctx context.Context) ([]*usermwpb.User, uint32, error) {
	var err error
	genesisRoles := []*approlemgrpb.AppRoleReq{}

	_, span := otel.Tracer(servicename.ServiceDomain).Start(ctx, "GetGenesisUsers")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "admin", "middleware", "GetGenesisUsers")

	genesisRoleStr := config.GetStringValueWithNameSpace(servicename2.ServiceDomain, constant1.KeyGenesisRole)

	err = json.Unmarshal([]byte(genesisRoleStr), &genesisRoles)
	if err != nil {
		return nil, 0, err
	}

	roles := []string{}
	for key := range genesisRoles {
		roles = append(roles, genesisRoles[key].GetRole())
	}

	roleInfos, _, err := approlemgrcli.GetAppRoles(ctx, &approlemgrpb.Conds{
		Roles: &commonpb.StringSliceVal{
			Op:    cruder.IN,
			Value: roles,
		},
	}, 0, int32(len(roles)))
	if err != nil {
		logger.Sugar().Errorw("GetGenesisUsers", "error", err)
		return nil, 0, err
	}

	roleIDs := []string{}
	appIDs := []string{}
	for _, val := range roleInfos {
		roleIDs = append(roleIDs, val.ID)
		appIDs = append(appIDs, val.AppID)
	}

	roleUsers, _, err := approleusermgrcli.GetAppRoleUsers(ctx, &approleusermgrpb.Conds{
		AppIDs: &commonpb.StringSliceVal{
			Op:    cruder.IN,
			Value: appIDs,
		},
		RoleIDs: &commonpb.StringSliceVal{
			Op:    cruder.IN,
			Value: roleIDs,
		},
	}, 0, 0)
	if err != nil {
		logger.Sugar().Errorw("GetGenesisUsers", "error", err)
		return nil, 0, err
	}
	if len(roleUsers) == 0 {
		return []*usermwpb.User{}, 0, nil
	}

	userIds := []string{}
	for _, val := range roleUsers {
		userIds = append(userIds, val.GetUserID())
	}

	infos, total, err := usermwcli.GetManyUsers(ctx, userIds)
	if err != nil {
		logger.Sugar().Errorw("GetGenesisUsers", "error", err)
		return nil, 0, err
	}

	return infos, total, nil
}
