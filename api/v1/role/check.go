package role

import (
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	"github.com/NpoolPlatform/appuser-manager/api/v2/approle"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	approlepb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approle"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validate(info *approlepb.AppRoleReq, userID string) error {
	if _, err := uuid.Parse(userID); err != nil {
		logger.Sugar().Errorw("validate", "userId", userID, "err", err)
		return status.Error(codes.InvalidArgument, "UserID is invalid")
	}

	if info.GetRole() == constant.GenesisRole {
		logger.Sugar().Errorw("validate", "Role", info.GetRole())
		return status.Error(codes.PermissionDenied, "permission denied")
	}

	if userID != info.GetCreatedBy() {
		logger.Sugar().Errorw("validate", "userId", userID, "CreatedBy", info.GetCreatedBy())
		return status.Error(codes.PermissionDenied, "permission denied")
	}

	err := approle.Validate(info)
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}
