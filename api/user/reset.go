package user

import (
	"context"

	user1 "github.com/NpoolPlatform/appuser-gateway/pkg/user"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) PreResetUser(ctx context.Context, in *npool.PreResetUserRequest) (*npool.PreResetUserResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(&in.AppID, true),
		user1.WithAccount(&in.Account, true),
		user1.WithAccountType(&in.AccountType, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"PreResetUser",
			"In", in,
			"Error", err,
		)
		return &npool.PreResetUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := handler.PreResetUser(ctx); err != nil {
		logger.Sugar().Errorw(
			"PreResetUser",
			"In", in,
			"Error", err,
		)
		return &npool.PreResetUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.PreResetUserResponse{}, nil
}

func (s *Server) ResetUser(ctx context.Context, in *npool.ResetUserRequest) (*npool.ResetUserResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(&in.AppID, true),
		user1.WithUserID(in.UserID, false),
		user1.WithAccount(&in.Account, true),
		user1.WithAccountType(&in.AccountType, true),
		user1.WithVerificationCode(&in.VerificationCode, true),
		user1.WithPasswordHash(in.PasswordHash, true),
		user1.WithRecoveryCode(in.RecoveryCode, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"ResetUser",
			"In", in,
			"Error", err,
		)
		return &npool.ResetUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := handler.ResetUser(ctx); err != nil {
		logger.Sugar().Errorw(
			"ResetUser",
			"In", in,
			"Error", err,
		)
		return &npool.ResetUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.ResetUserResponse{}, nil
}
