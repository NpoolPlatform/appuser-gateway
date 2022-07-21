package appuser

import (
	"context"
	"fmt"
	"github.com/NpoolPlatform/api-manager/pkg/db/ent"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	appusercrud "github.com/NpoolPlatform/message/npool/appusermgrv2/appuser"
)

func GetUserByAccount(ctx context.Context, appID, account string) (*appusercrud.AppUser, error) {
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
