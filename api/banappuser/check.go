package banappuser

import (
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	"github.com/NpoolPlatform/message/npool/appusermgrv2/banappuser"
	"github.com/google/uuid"
	"google.golang.org/grpc/status"
)

func validate(info *banappuser.BanAppUserReq) error {
	if _, err := uuid.Parse(info.GetAppID()); err != nil {
		logger.Sugar().Error("AppID is invalid")
		return status.Error(npool.ErrParams, appusergw.ErrMsgAppIDInvalid)
	}

	if _, err := uuid.Parse(info.GetUserID()); err != nil {
		logger.Sugar().Error("UserID is invalid")
		return status.Error(npool.ErrParams, appusergw.ErrMsgUserIDInvalid)
	}

	if _, err := uuid.Parse(info.GetMessage()); err != nil {
		logger.Sugar().Error("Message empty")
		return status.Error(npool.ErrParams, appusergw.ErrMsgMessageEmpty)
	}

	return nil
}
