package user

import (
	"context"

	roleuser1 "github.com/NpoolPlatform/appuser-gateway/pkg/role/user"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/role/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, in *npool.CreateUserRequest) (*npool.CreateUserResponse, error) {
	handler, err := roleuser1.NewHandler(
		ctx,
		roleuser1.WithAppID(in.GetAppID()),
		roleuser1.WithUserID(in.GetTargetUserID()),
		roleuser1.WithRoleID(in.GetRoleID()),
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
	handler, err := roleuser1.NewHandler(
		ctx,
		roleuser1.WithAppID(in.GetTargetAppID()),
		roleuser1.WithUserID(in.GetTargetUserID()),
		roleuser1.WithRoleID(in.GetRoleID()),
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
