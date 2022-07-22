package user

import (
	"context"
	"fmt"
	"github.com/NpoolPlatform/api-manager/pkg/db/ent"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	approlecrud "github.com/NpoolPlatform/message/npool/appusermgrv2/approle"
	approleusercrud "github.com/NpoolPlatform/message/npool/appusermgrv2/approleuser"
	appusercrud "github.com/NpoolPlatform/message/npool/appusermgrv2/appuser"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
)

func GetUserByAccount(ctx context.Context, appID, account string) (*appusercrud.AppUser, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetUserByAccount")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.AddEvent("call grpc GetAppUserOnlyV2")
	resp, err := grpc.GetAppUserOnlyV2(ctx, &appusercrud.Conds{
		PhoneNo: &npool.StringVal{
			Value: account,
			Op:    cruder.EQ,
		},
		AppID: &npool.StringVal{
			Value: appID,
			Op:    cruder.EQ,
		},
	})
	if ent.IsNotFound(err) {
		span.AddEvent("call grpc GetAppUserOnlyV2")
		resp, err = grpc.GetAppUserOnlyV2(ctx, &appusercrud.Conds{
			EmailAddress: &npool.StringVal{
				Value: account,
				Op:    cruder.EQ,
			},
			AppID: &npool.StringVal{
				Value: appID,
				Op:    cruder.EQ,
			},
		})
		if err != nil {
			if ent.IsNotFound(err) {
				logger.Sugar().Error("Account not exist")
				return nil, fmt.Errorf("account not exist")
			}
			return nil, err
		}
	} else {
		return nil, err
	}

	return resp, err
}

func GetUserRolesByUser(ctx context.Context, appID, userID string, limit, offset int32) ([]*approlecrud.AppRole, uint32, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetUserByAccount")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.AddEvent("call grpc GetAppRoleUsersV2")
	resp, err := grpc.GetAppRoleUsersV2(ctx, &approleusercrud.Conds{
		AppID: &npool.StringVal{
			Value: appID,
			Op:    cruder.EQ,
		},
		UserID: &npool.StringVal{
			Value: userID,
			Op:    cruder.EQ,
		},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("fail get app role user: %v", err)
	}

	roleIDs := []string{}
	for _, info := range resp {
		roleIDs = append(roleIDs, info.GetRoleID())
	}

	roles, total, err := grpc.GetAppRolesV2(ctx, &approlecrud.Conds{
		IDIn: &npool.StringSlicesVal{
			Value: roleIDs,
			Op:    cruder.IN,
		},
	}, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return roles, total, err
}
