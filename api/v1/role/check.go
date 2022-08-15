package role

// import (
//	"context"
//
//	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
//	"github.com/NpoolPlatform/appuser-manager/api/v2/approle"
//	approle2 "github.com/NpoolPlatform/appuser-manager/pkg/client/approle"
//	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
//	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
//	"github.com/NpoolPlatform/message/npool"
//	approlepb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approle"
//	"github.com/google/uuid"
//	"google.golang.org/grpc/codes"
//	"google.golang.org/grpc/status"
//)
//
// func validate(ctx context.Context, info *approlepb.AppRoleReq, userID string) error {
//	if _, err := uuid.Parse(userID); err != nil {
//		logger.Sugar().Errorw("validate", "userId", userID, "err", err)
//		return status.Error(codes.InvalidArgument, "UserID is invalid")
//	}
//
//	if info.GetRole() == constant.GenesisRole {
//		logger.Sugar().Errorw("validate", "Role", info.GetRole())
//		return status.Error(codes.PermissionDenied, "permission denied")
//	}
//
//	if userID != info.GetCreatedBy() {
//		logger.Sugar().Errorw("validate", "userId", userID, "CreatedBy", info.GetCreatedBy())
//		return status.Error(codes.PermissionDenied, "permission denied")
//	}
//
//	err := approle.Validate(info)
//	if err != nil {
//		logger.Sugar().Errorw("validate", "err", err)
//		return status.Error(codes.InvalidArgument, err.Error())
//	}
//
//	exist, err := approle2.ExistAppRoleConds(ctx, &approlepb.Conds{
//		AppID: &npool.StringVal{
//			Op:    cruder.EQ,
//			Value: info.GetAppID(),
//		},
//		Role: &npool.StringVal{
//			Op:    cruder.EQ,
//			Value: info.GetRole(),
//		},
//	})
//	if err != nil {
//		return status.Error(codes.Internal, err.Error())
//	}
//
//	if exist {
//		return status.Error(codes.AlreadyExists, "app role already exists")
//	}
//	return nil
//}
