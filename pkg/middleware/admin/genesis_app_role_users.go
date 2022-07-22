package admin

import (
	"context"
	"fmt"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	approleusercrud "github.com/NpoolPlatform/message/npool/appusermgrv2/approleuser"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
)

func GetAppGenesisAppRoleUsers(ctx context.Context, appID, userID string) ([]*approleusercrud.AppRoleUser, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppGenesisAppRoleUsers")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.AddEvent("call middleware methods GetGenesisRole")
	role, err := GetGenesisRole(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail get genesis role: %v", err)
	}

	resp, err := grpc.GetAppRoleUsersV2(ctx, &approleusercrud.Conds{
		AppID: &npool.StringVal{
			Value: appID,
			Op:    cruder.EQ,
		},
		RoleID: &npool.StringVal{
			Value: role.GetID(),
			Op:    cruder.EQ,
		},
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
