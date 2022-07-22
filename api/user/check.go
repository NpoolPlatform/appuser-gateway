package user

import (
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	appusercrud "github.com/NpoolPlatform/message/npool/appusermgrv2/appuser"
	"github.com/google/uuid"
	"google.golang.org/grpc/status"
)

func validate(info *appusercrud.AppUserReq) error {
	if _, err := uuid.Parse(info.GetAppID()); err != nil {
		logger.Sugar().Error("AppID is invalid")
		return status.Error(npool.ErrParams, appusergw.ErrMsgAppIDInvalid)
	}

	if info.GetPhoneNo() == "" && info.GetEmailAddress() == "" {
		logger.Sugar().Error("PhoneNo and EmailAddress is empty")
		return status.Error(npool.ErrParams, appusergw.ErrMsgPhoneAndEmailMustExistOne)
	}
	return nil
}
