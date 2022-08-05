package admin

import (
	"context"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	serconst "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	"github.com/NpoolPlatform/appuser-manager/pkg/client/approle"
	"github.com/NpoolPlatform/appuser-manager/pkg/client/approleuser"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	approlepb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approle"
	approleuserpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approleuser"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
)

func GetAppGenesisAppRoleUsers(ctx context.Context) ([]*approleuserpb.AppRoleUser, error) {
	var err error

	_, span := otel.Tracer(serconst.ServiceName).Start(ctx, "CreateGenesisRole")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "admin", "middleware", "CreateAppRoles")

	roles, _, err := approle.GetAppRoles(ctx, &approlepb.Conds{
		Roles: &npool.StringSliceVal{
			Op:    cruder.IN,
			Value: []string{constant.GenesisRole, constant.ChurchRole},
		},
	}, 2, 0)
	if err != nil {
		logger.Sugar().Errorw("CreateGenesisRole", "error", err)
		return nil, err
	}

	roleIDs := []string{}
	for _, val := range roles {
		roleIDs = append(roleIDs, val.ID)
	}

	resp, _, err := approleuser.GetAppRoleUsers(ctx, &approleuserpb.Conds{
		AppIDs: &npool.StringSliceVal{
			Op:    cruder.IN,
			Value: []string{constant.GenesisAppID, constant.ChurchAppID},
		},
		RoleIDs: &npool.StringSliceVal{
			Op:    cruder.IN,
			Value: roleIDs,
		},
	}, 2, 0)
	if err != nil {
		logger.Sugar().Errorw("GetAppGenesisAppRoleUsers", "error", err)
		return nil, err
	}

	return resp, nil
}
