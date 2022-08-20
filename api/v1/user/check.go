package user

import (
	"context"

	mw "github.com/NpoolPlatform/appuser-middleware/api/v1/user"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	mgruser "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validate(ctx context.Context, info *mgruser.UserReq) error {
	if info == nil {
		logger.Sugar().Errorw("validate", "err", "params is empty")
		return status.Error(codes.InvalidArgument, "params is empty")
	}

	err := mw.Validate(ctx, info)
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}
