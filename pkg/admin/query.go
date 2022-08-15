package admin

import (
	"context"
	"encoding/json"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	constants "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	approlemgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approle"
	"github.com/NpoolPlatform/appuser-manager/pkg/client/approleuser"
	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/app"
	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/role"
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
	roles := []string{}

	_, span := otel.Tracer(constants.ServiceName).Start(ctx, "GetGenesisRole")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "admin", "manager", "GetAppRoles")

	hostname := config.GetStringValueWithNameSpace("", config.KeyHostname)
	genesisRoleStr := config.GetStringValueWithNameSpace(hostname, constant.KeyGenesisRole)

	err = json.Unmarshal([]byte(genesisRoleStr), &roles)
	if err != nil {
		return nil, err
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

	return rolemwcli.GetManyRoles(ctx, roleIDs)
}

func GetAppGenesisAppRoleUsers(ctx context.Context) ([]*rolemwpb.RoleUser, error) {
	var err error
	genesisRoles := []*approlepb.AppRoleReq{}

	_, span := otel.Tracer(constants.ServiceName).Start(ctx, "CreateGenesisRole")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "admin", "middleware", "CreateAppRoles")

	hostname := config.GetStringValueWithNameSpace("", config.KeyHostname)
	genesisAppStr := config.GetStringValueWithNameSpace(hostname, constant.KeyGenesisApp)

	err = json.Unmarshal([]byte(genesisAppStr), &genesisRoles)
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
		logger.Sugar().Errorw("CreateGenesisRole", "error", err)
		return nil, err
	}

	roleIDs := []string{}
	for _, val := range roleInfos {
		roleIDs = append(roleIDs, val.ID)
	}

	roleUsers, _, err := approleuser.GetAppRoleUsers(ctx, &approleuserpb.Conds{
		AppIDs: &npool.StringSliceVal{
			Op:    cruder.IN,
			Value: []string{constant.GenesisAppID, constant.ChurchAppID},
		},
		RoleIDs: &npool.StringSliceVal{
			Op:    cruder.IN,
			Value: roleIDs,
		},
	}, 0, 0)
	if err != nil {
		logger.Sugar().Errorw("GetAppGenesisAppRoleUsers", "error", err)
		return nil, err
	}

	roleUserIds := []string{}
	for _, val := range roleUsers {
		roleUserIds = append(roleUserIds, val.GetID())
	}

	return rolemwcli.GetManyRoleUsers(ctx, roleUserIds)
}
