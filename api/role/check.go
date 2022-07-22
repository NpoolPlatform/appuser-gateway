package role

import (
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	"github.com/NpoolPlatform/message/npool/appusermgrv2/approle"
	"github.com/google/uuid"
	"google.golang.org/grpc/status"
)

func validate(info *approle.AppRoleReq, userID string) error {
	if _, err := uuid.Parse(userID); err != nil {
		logger.Sugar().Error("UserID is invalid")
		return status.Error(npool.ErrParams, appusergw.ErrMsgUserIDInvalid)
	}

	if info.GetRole() == constant.GenesisRole {
		return status.Error(npool.ErrPermissionDenied, appusergw.ErrMsgPermissionDenied)
	}

	if userID != info.GetCreatedBy() {
		return status.Error(npool.ErrPermissionDenied, appusergw.ErrMsgPermissionDenied)
	}

	if _, err := uuid.Parse(info.GetAppID()); err != nil {
		logger.Sugar().Error("AppID is invalid")
		return status.Error(npool.ErrParams, appusergw.ErrMsgAppIDInvalid)
	}

	if _, err := uuid.Parse(info.GetCreatedBy()); err != nil {
		logger.Sugar().Error("CreatedBy is invalid")
		return status.Error(npool.ErrParams, appusergw.ErrMsgCreatedByInvalid)
	}

	if _, err := uuid.Parse(info.GetRole()); err != nil {
		logger.Sugar().Error("Role empty")
		return status.Error(npool.ErrParams, appusergw.ErrMsgMessageEmpty)
	}

	return nil
}
