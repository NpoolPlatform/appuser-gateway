package extra

import (
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	appuserextracrud "github.com/NpoolPlatform/message/npool/appusermgrv2/appuserextra"
	"github.com/google/uuid"
	"google.golang.org/grpc/status"
)

func validate(info *appuserextracrud.AppUserExtraReq) error {
	if _, err := uuid.Parse(info.GetAppID()); err != nil {
		logger.Sugar().Error("AppID is invalid")
		return status.Error(npool.ErrParams, appusergw.ErrMsgAppIDInvalid)
	}

	if _, err := uuid.Parse(info.GetUserID()); err != nil {
		logger.Sugar().Error("UserID is invalid")
		return status.Error(npool.ErrParams, appusergw.ErrMsgUserIDInvalid)
	}

	return nil
}
