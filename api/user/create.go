package user

import (
	"context"

	user1 "github.com/NpoolPlatform/appuser-gateway/pkg/user"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, in *npool.CreateUserRequest) (*npool.CreateUserResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(in.GetAppID()),
		user1.WithEmailAddress(in.EmailAddress),
		user1.WithPhoneNO(in.PhoneNO),
		user1.WithPasswordHash(in.PasswordHash),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateUser",
			"In", in,
			"Error", err,
		)
		return &npool.CreateUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateUser(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateUser",
			"In", in,
			"Error", err,
		)
		return &npool.CreateUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateUserResponse{
		Info: info,
	}, nil
}

func (s *Server) CreateAppUser(ctx context.Context, in *npool.CreateAppUserRequest) (*npool.CreateAppUserResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(in.GetTargetAppID()),
		user1.WithEmailAddress(in.EmailAddress),
		user1.WithPhoneNO(in.PhoneNO),
		user1.WithPasswordHash(&in.PasswordHash),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAppUser",
			"In", in,
			"Error", err,
		)
		return &npool.CreateAppUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateUser(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAppUser",
			"In", in,
			"Error", err,
		)
		return &npool.CreateAppUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateAppUserResponse{
		Info: info,
	}, nil
}
