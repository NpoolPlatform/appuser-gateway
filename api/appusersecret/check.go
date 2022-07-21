package appusersecret

import (
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw/appusersecret"
	appusersecretcrud "github.com/NpoolPlatform/message/npool/appusermgrv2/appusersecret"
	"github.com/google/uuid"
	"google.golang.org/grpc/status"
)

func validate(info *appusersecretcrud.AppUserSecretReq) error {
	if _, err := uuid.Parse(info.GetAppID()); err != nil {
		logger.Sugar().Error("AppID is invalid")
		return status.Error(npool.ErrParams, appusersecret.ErrMsgAppIDInvalid)
	}

	if _, err := uuid.Parse(info.GetUserID()); err != nil {
		logger.Sugar().Error("UserID is invalid")
		return status.Error(npool.ErrParams, appusersecret.ErrMsgUserIDInvalid)
	}

	if info.GetPasswordHash() == "" {
		logger.Sugar().Error("PasswordHash empty")
		return status.Error(npool.ErrParams, appusersecret.ErrMsgPasswordHashEmpty)
	}

	if info.GetSalt() == "" {
		logger.Sugar().Error("Salt empty")
		return status.Error(npool.ErrParams, appusersecret.ErrMsgSaltEmpty)
	}

	if info.GetGoogleSecret() == "" {
		logger.Sugar().Error("GoogleSecret is empty")
		return status.Error(npool.ErrParams, appusersecret.ErrMsgGoogleSecretEmpty)
	}

	return nil
}
