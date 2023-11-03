package user

import (
	"context"

	user1 "github.com/NpoolPlatform/appuser-gateway/pkg/user"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteUser(ctx context.Context, in *npool.DeleteUserRequest) (*npool.DeleteUserResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(&in.AppID, true),
		user1.WithUserID(&in.TargetUserID, true),
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
