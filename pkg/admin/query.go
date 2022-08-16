package admin

import (
	"context"
	"encoding/json"
	"fmt"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	constants "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	approlemgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approle"
	"github.com/NpoolPlatform/appuser-manager/pkg/client/approleuser"
	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role"
	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	appcrud "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/app"
	approlepb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approle"
	approleuserpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approleuser"
	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
	rolemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
)

func GetAdminApps(ctx context.Context) ([]*appmwpb.App, error) {
	var err error
	genesisApps := []*appcrud.App{}

	_, span := otel.Tracer(constants.ServiceName).Start(ctx, "GetAdminApps")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "admin", "manager", "GetApps")

	hostname := config.GetStringValueWithNameSpace("", config.KeyHostname)
	genesisAppStr := config.GetStringValueWithNameSpace(hostname, constant.KeyGenesisApp)

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

func GetGenesisRoles(ctx context.Context) ([]*rolemwpb.Role, error) {
	var err error
	appRoles := []*approlepb.AppRoleReq{}

	_, span := otel.Tracer(constants.ServiceName).Start(ctx, "GetGenesisRoles")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "admin", "manager", "GetGenesisRoles")

	hostname := config.GetStringValueWithNameSpace("", config.KeyHostname)
	genesisRoleStr := config.GetStringValueWithNameSpace(hostname, constant.KeyGenesisRole)

	err = json.Unmarshal([]byte(genesisRoleStr), &appRoles)
	if err != nil {
		return nil, err
	}

	roles := []string{}
	for _, val := range appRoles {
		roles = append(roles, val.GetRole())
	}

	roleInfos, _, err := approlemgrcli.GetAppRoles(ctx, &approlepb.Conds{
		Roles: &npool.StringSliceVal{
			Value: roles,
			Op:    cruder.IN,
		},
	}, 0, int32(len(roles)))
	if err != nil {
		logger.Sugar().Errorw("GetGenesisRole", "err", err)
		return nil, err
	}

	roleIDs := []string{}
	for _, val := range roleInfos {
		roleIDs = append(roleIDs, val.GetID())
	}

	infos, _, err := rolemwcli.GetManyRoles(ctx, roleIDs)
	if err != nil {
		logger.Sugar().Errorw("GetGenesisRoles", "error", err)
		return nil, err
	}

	return infos, nil
}

func GetGenesisRoleUsers(ctx context.Context) ([]*rolemwpb.RoleUser, error) {
	var err error
	genesisRoles := []*approlepb.AppRoleReq{}

	_, span := otel.Tracer(constants.ServiceName).Start(ctx, "GetGenesisRoleUsers")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "admin", "middleware", "GetGenesisRoleUsers")

	hostname := config.GetStringValueWithNameSpace("", config.KeyHostname)
	genesisRoleStr := config.GetStringValueWithNameSpace(hostname, constant.KeyGenesisRole)

	err = json.Unmarshal([]byte(genesisRoleStr), &genesisRoles)
	if err != nil {
		return nil, err
	}

	roles := []string{}
	for key := range genesisRoles {
		roles = append(roles, genesisRoles[key].GetRole())
	}

	roleInfos, _, err := approlemgrcli.GetAppRoles(ctx, &approlepb.Conds{
		Roles: &npool.StringSliceVal{
			Op:    cruder.IN,
			Value: roles,
		},
	}, 0, int32(len(roles)))
	if err != nil {
		logger.Sugar().Errorw("GetGenesisRoleUsers", "error", err)
		return nil, err
	}

	roleIDs := []string{}
	appIDs := []string{}
	for _, val := range roleInfos {
		roleIDs = append(roleIDs, val.ID)
		appIDs = append(appIDs, val.AppID)
	}

	roleUsers, _, err := approleuser.GetAppRoleUsers(ctx, &approleuserpb.Conds{
		AppIDs: &npool.StringSliceVal{
			Op:    cruder.IN,
			Value: appIDs,
		},
		RoleIDs: &npool.StringSliceVal{
			Op:    cruder.IN,
			Value: roleIDs,
		},
	}, 0, 0)
	if err != nil {
		logger.Sugar().Errorw("GetGenesisRoleUsers", "error", err)
		return nil, err
	}
	if len(roleUsers) == 0 {
		logger.Sugar().Errorw("GetGenesisRoleUsers", "error", "not found")
		return nil, fmt.Errorf("role user not found")
	}
	roleUserIds := []string{}
	for _, val := range roleUsers {
		roleUserIds = append(roleUserIds, val.GetID())
	}

	infos, _, err := rolemwcli.GetManyRoleUsers(ctx, roleUserIds)
	if err != nil {
		logger.Sugar().Errorw("GetGenesisRoleUsers", "error", err)
		return nil, err
	}

	return infos, nil
}
