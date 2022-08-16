package ban

import (
	"github.com/NpoolPlatform/appuser-manager/api/v2/banapp"
	"github.com/NpoolPlatform/appuser-manager/api/v2/banappuser"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	banapppb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/banapp"
	banappuserpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/banappuser"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validate(info *banapppb.BanAppReq) error {
	err := banapp.Validate(info)
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}

func validateBanUser(info *banappuserpb.BanAppUserReq) error {
	err := banappuser.Validate(info)
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}
