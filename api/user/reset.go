package user

import (
	"context"

	user1 "github.com/NpoolPlatform/appuser-gateway/pkg/user"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) ResetUser(ctx context.Context, in *npool.ResetUserRequest) (*npool.ResetUserResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(in.GetAppID()),
		user1.WithUserID(in.UserID),
		user1.WithAccount(&in.Account, &in.AccountType),
		user1.WithVerificationCode(&in.VerificationCode),
		user1.WithPasswordHash(in.PasswordHash),
		user1.WithRecoveryCode(in.RecoveryCode),
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
