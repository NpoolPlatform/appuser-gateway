//nolint:dupl
package user

import (
	"context"

	roleuser1 "github.com/NpoolPlatform/appuser-gateway/pkg/role/user"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/role/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteUser(ctx context.Context, in *npool.DeleteUserRequest) (*npool.DeleteUserResponse, error) {
	handler, err := roleuser1.NewHandler(
		ctx,
		roleuser1.WithID(&in.ID, true),
		roleuser1.WithEntID(&in.TargetUserID, true),
		roleuser1.WithAppID(&in.AppID, true),
		roleuser1.WithUserID(&in.TargetUserID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteUser",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteUser(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteUser",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteUserResponse{
		Info: info,
	}, nil
}

func (s *Server) DeleteAppUser(ctx context.Context, in *npool.DeleteAppUserRequest) (*npool.DeleteAppUserResponse, error) {
	handler, err := roleuser1.NewHandler(
		ctx,
		roleuser1.WithID(&in.ID, true),
		roleuser1.WithEntID(&in.TargetUserID, true),
		roleuser1.WithAppID(&in.TargetAppID, true),
		roleuser1.WithUserID(&in.TargetUserID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteAppUser",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteAppUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteUser(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteAppUser",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteAppUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteAppUserResponse{
		Info: info,
	}, nil
}
