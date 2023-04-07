package authing

import (
	"github.com/NpoolPlatform/appuser-manager/api/authing/auth"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	pb "github.com/NpoolPlatform/message/npool/appuser/gw/v1/authing"
	authmgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/authing/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validate(info *pb.CreateAppAuthRequest) error {
	if info == nil {
		logger.Sugar().Errorw("validate", "err", "params is empty")
		return status.Error(codes.InvalidArgument, "params is empty")
	}

	err := auth.Validate(&authmgrpb.AuthReq{
		AppID:    &info.TargetAppID,
		RoleID:   info.RoleID,
		UserID:   info.TargetUserID,
		Resource: &info.Resource,
		Method:   &info.Method,
	})
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}
