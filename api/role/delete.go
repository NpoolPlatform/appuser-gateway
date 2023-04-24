package role

import (
	"context"

	role1 "github.com/NpoolPlatform/appuser-gateway/pkg/role"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/role"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteRole(ctx context.Context, in *npool.DeleteRoleRequest) (*npool.DeleteRoleResponse, error) {
	handler, err := role1.NewHandler(
		ctx,
		role1.WithID(&in.ID),
		role1.WithAppID(in.GetAppID()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteRole",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteRoleResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteRole(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteRole",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteRoleResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteRoleResponse{
		Info: info,
	}, nil
}

func (s *Server) DeleteAppRole(ctx context.Context, in *npool.DeleteAppRoleRequest) (*npool.DeleteAppRoleResponse, error) {
	handler, err := role1.NewHandler(
		ctx,
		role1.WithID(&in.ID),
		role1.WithAppID(in.GetTargetAppID()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteAppRole",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteAppRoleResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteRole(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteAppRole",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteAppRoleResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteAppRoleResponse{
		Info: info,
	}, nil
}
