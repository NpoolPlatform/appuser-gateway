package thirdparty

import (
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	thirdpartycrud "github.com/NpoolPlatform/message/npool/appusermgrv2/appuserthirdparty"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validate(info *thirdpartycrud.AppUserThirdPartyReq) error {
	if _, err := uuid.Parse(info.GetAppID()); err != nil {
		logger.Sugar().Error("AppID is invalid")
		return status.Error(npool.ErrParams, appusergw.ErrMsgAppIDInvalid)
	}

	if _, err := uuid.Parse(info.GetUserID()); err != nil {
		logger.Sugar().Error("UserID is invalid")
		return status.Error(npool.ErrParams, appusergw.ErrMsgUserIDInvalid)
	}

	if info.GetThirdPartyUserID() == "" {
		logger.Sugar().Error("ThirdPartyUserID is invalid")
		return status.Error(npool.ErrParams, appusergw.ErrMsgThirdPartyUserIDEmpty)
	}

	if info.GetThirdPartyID() == "" {
		logger.Sugar().Error("ThirdPartyID is invalid")
		return status.Error(codes.InvalidArgument, appusergw.ErrMsgThirdPartyIDEmpty)
	}

	return nil
}
