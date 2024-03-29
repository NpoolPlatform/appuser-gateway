package user

import (
	"context"

	user1 "github.com/NpoolPlatform/appuser-gateway/pkg/user"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) BindUser(ctx context.Context, in *npool.BindUserRequest) (*npool.BindUserResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(&in.AppID, true),
		user1.WithUserID(&in.UserID, true),
		user1.WithAccount(in.Account, true),
		user1.WithAccountType(in.AccountType, true),
		user1.WithNewAccount(in.NewAccount, true),
		user1.WithNewAccountType(in.NewAccountType, true),
		user1.WithNewVerificationCode(in.NewVerificationCode, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"BindUser",
			"In", in,
			"Error", err,
		)
		return &npool.BindUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.BindUser(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"BindUser",
			"In", in,
			"Error", err,
		)
		return &npool.BindUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.BindUserResponse{
		Info: info,
	}, nil
}

func (s *Server) UnbindOAuth(ctx context.Context, in *npool.UnbindOAuthRequest) (*npool.UnbindOAuthResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(&in.AppID, true),
		user1.WithUserID(&in.UserID, true),
		user1.WithAccount(&in.Account, true),
		user1.WithAccountType(&in.AccountType, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UnbindOAuth",
			"In", in,
			"Error", err,
		)
		return &npool.UnbindOAuthResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	err = handler.UnbindOAuth(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UnbindOAuth",
			"In", in,
			"Error", err,
		)
		return &npool.UnbindOAuthResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UnbindOAuthResponse{}, nil
}
