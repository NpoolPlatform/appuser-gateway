//nolint:dupl
package user

import (
	"context"

	user1 "github.com/NpoolPlatform/appuser-gateway/pkg/user"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) BanUser(ctx context.Context, in *npool.BanUserRequest) (*npool.BanUserResponse, error) {
	updateCacheMode := user1.UpdateCacheIfExist
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(&in.AppID, true),
		user1.WithUserID(&in.TargetUserID, true),
		user1.WithBanned(&in.Banned, true),
		user1.WithBanMessage(&in.BanMessage, true),
		user1.WithUpdateCacheMode(&updateCacheMode, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"BanUser",
			"In", in,
			"Error", err,
		)
		return &npool.BanUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if _, err := handler.UpdateUser(ctx); err != nil {
		logger.Sugar().Errorw(
			"BanUser",
			"In", in,
			"Error", err,
		)
		return &npool.BanUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.BanUserResponse{}, nil
}

func (s *Server) BanAppUser(ctx context.Context, in *npool.BanAppUserRequest) (*npool.BanAppUserResponse, error) {
	updateCacheMode := user1.UpdateCacheIfExist
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(&in.TargetAppID, true),
		user1.WithUserID(&in.TargetUserID, true),
		user1.WithBanned(&in.Banned, true),
		user1.WithBanMessage(&in.BanMessage, true),
		user1.WithUpdateCacheMode(&updateCacheMode, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"BanAppUser",
			"In", in,
			"Error", err,
		)
		return &npool.BanAppUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if _, err := handler.UpdateUser(ctx); err != nil {
		logger.Sugar().Errorw(
			"BanAppUser",
			"In", in,
			"Error", err,
		)
		return &npool.BanAppUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.BanAppUserResponse{}, nil
}
