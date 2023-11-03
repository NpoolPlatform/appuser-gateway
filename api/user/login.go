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

func (s *Server) Login(ctx context.Context, in *npool.LoginRequest) (*npool.LoginResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(&in.AppID, true),
		user1.WithAccount(&in.Account, true),
		user1.WithAccountType(&in.AccountType, true),
		user1.WithPasswordHash(&in.PasswordHash, true),
		user1.WithManMachineSpec(&in.ManMachineSpec, true),
		user1.WithEnvironmentSpec(&in.EnvironmentSpec, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"Login",
			"In", in,
			"Error", err,
		)
		return &npool.LoginResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.Login(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"Login",
			"In", in,
			"Error", err,
		)
		return &npool.LoginResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.LoginResponse{
		Info: info,
	}, nil
}

func (s *Server) LoginVerify(ctx context.Context, in *npool.LoginVerifyRequest) (*npool.LoginVerifyResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(&in.AppID, true),
		user1.WithUserID(&in.UserID, true),
		user1.WithAccount(&in.Account, true),
		user1.WithAccountType(&in.AccountType, true),
		user1.WithToken(&in.Token, true),
		user1.WithVerificationCode(&in.VerificationCode, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"LoginVerify",
			"In", in,
			"Error", err,
		)
		return &npool.LoginVerifyResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.LoginVerify(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"LoginVerify",
			"In", in,
			"Error", err,
		)
		return &npool.LoginVerifyResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.LoginVerifyResponse{
		Info: info,
	}, nil
}

func (s *Server) Logined(ctx context.Context, in *npool.LoginedRequest) (*npool.LoginedResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(&in.AppID, true),
		user1.WithUserID(&in.UserID, true),
		user1.WithToken(&in.Token, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"Logined",
			"In", in,
			"Error", err,
		)
		return &npool.LoginedResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.Logined(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"Logined",
			"In", in,
			"Error", err,
		)
		return &npool.LoginedResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.LoginedResponse{
		Info: info,
	}, nil
}

func (s *Server) Logout(ctx context.Context, in *npool.LogoutRequest) (*npool.LogoutResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(&in.AppID, true),
		user1.WithUserID(&in.UserID, true),
		user1.WithToken(&in.Token, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"Logout",
			"In", in,
			"Error", err,
		)
		return &npool.LogoutResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.Logout(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"Logout",
			"In", in,
			"Error", err,
		)
		return &npool.LogoutResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.LogoutResponse{
		Info: info,
	}, nil
}

func (s *Server) GetLoginHistories(ctx context.Context, in *npool.GetLoginHistoriesRequest) (*npool.GetLoginHistoriesResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(&in.AppID, true),
		user1.WithUserID(&in.UserID, true),
		user1.WithOffset(in.GetOffset()),
		user1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetLoginHistories",
			"In", in,
			"Error", err,
		)
		return &npool.GetLoginHistoriesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetLoginHistories(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetLoginHistories",
			"In", in,
			"Error", err,
		)
		return &npool.GetLoginHistoriesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetLoginHistoriesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
