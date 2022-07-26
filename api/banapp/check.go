package banapp

import (
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	banappcrud "github.com/NpoolPlatform/message/npool/appusermgrv2/banapp"
	"github.com/google/uuid"
	"google.golang.org/grpc/status"
)

func validate(info *banappcrud.BanAppReq) error {
	if _, err := uuid.Parse(info.GetAppID()); err != nil {
		logger.Sugar().Error("AppID is invalid")
		return status.Error(npool.ErrParams, appusergw.ErrMsgAppIDInvalid)
	}

	if info.GetMessage() == "" {
		logger.Sugar().Error("Message is empty")
		return status.Error(npool.ErrParams, appusergw.ErrMsgMessageEmpty)
	}

	return nil
}
