package banappuser

import (
	"github.com/NpoolPlatform/appuser-manager/api/v2/banappuser"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	banappuserpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/banappuser"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validate(info *banappuserpb.BanAppUserReq) error {
	err := banappuser.Validate(info)
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}
