package roleuser

import (
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	"github.com/NpoolPlatform/message/npool/appusermgrv2/approleuser"
	"github.com/google/uuid"
	"google.golang.org/grpc/status"
)

func validate(info *approleuser.AppRoleUserReq) error {
	if _, err := uuid.Parse(info.GetAppID()); err != nil {
		logger.Sugar().Error("UserID is invalid")
		return status.Error(npool.ErrParams, appusergw.ErrMsgUserIDInvalid)
	}

	if _, err := uuid.Parse(info.GetUserID()); err != nil {
		logger.Sugar().Error("UserID is invalid")
		return status.Error(npool.ErrParams, appusergw.ErrMsgUserIDInvalid)
	}

	if _, err := uuid.Parse(info.GetRoleID()); err != nil {
		logger.Sugar().Error("RoleID is invalid")
		return status.Error(npool.ErrParams, appusergw.ErrMsgRoleIDInvalid)
	}

	return nil
}
