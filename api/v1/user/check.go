package user

import (
	mw "github.com/NpoolPlatform/appuser-middleware/api/v1/user"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validate(info *user.UserReq) error {
	err := mw.Validate(info)
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}
