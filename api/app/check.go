package app

import (
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	appcrud "github.com/NpoolPlatform/message/npool/appusermgrv2/app"
	"github.com/google/uuid"
	"google.golang.org/grpc/status"
)

func validate(info *appcrud.AppReq) error {
	if _, err := uuid.Parse(info.GetCreatedBy()); err != nil {
		logger.Sugar().Error("CreatedBy is invalid")
		return status.Error(npool.ErrParams, appusergw.ErrMsgCreatedByInvalid)
	}

	if info.GetName() == "" {
		logger.Sugar().Error("Name is empty")
		return status.Error(npool.ErrParams, appusergw.ErrMsgAppNameEmpty)
	}

	if info.GetLogo() == "" {
		logger.Sugar().Error("Logo is empty")
		return status.Error(npool.ErrParams, appusergw.ErrMsgAppLogoEmpty)
	}

	return nil
}
