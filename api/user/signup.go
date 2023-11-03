package user

import (
	"context"

	user1 "github.com/NpoolPlatform/appuser-gateway/pkg/user"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Signup(ctx context.Context, in *user.SignupRequest) (*user.SignupResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(&in.AppID, true),
		user1.WithPasswordHash(&in.PasswordHash, true),
		user1.WithAccount(&in.Account, true),
		user1.WithAccountType(&in.AccountType, true),
		user1.WithVerificationCode(&in.VerificationCode, true),
		user1.WithInvitationCode(in.InvitationCode, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"Signup",
			"In", in,
			"Error", err,
		)
		return &user.SignupResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.Signup(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"Signup",
			"In", in,
			"Error", err,
		)
		return &user.SignupResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &user.SignupResponse{
		Info: info,
	}, nil
}
