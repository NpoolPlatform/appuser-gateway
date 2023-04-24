package role

import (
	"context"

	role1 "github.com/NpoolPlatform/appuser-gateway/pkg/role"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/role"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateRole(ctx context.Context, in *npool.CreateRoleRequest) (*npool.CreateRoleResponse, error) {
	handler, err := role1.NewHandler(
		ctx,
		role1.WithAppID(in.GetAppID()),
		role1.WithCreatedBy(&in.UserID),
		role1.WithRole(&in.RoleName),
		role1.WithDescription(&in.Description),
		role1.WithDefault(&in.Default),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateRole",
			"In", in,
			"Error", err,
		)
		return &npool.CreateRoleResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateRole(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateRole",
			"In", in,
			"Error", err,
		)
		return &npool.CreateRoleResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateRoleResponse{
		Info: info,
	}, nil
}

func (s *Server) CreateAppRole(ctx context.Context, in *npool.CreateAppRoleRequest) (*npool.CreateAppRoleResponse, error) {
	handler, err := role1.NewHandler(
		ctx,
		role1.WithAppID(in.GetTargetAppID()),
		role1.WithCreatedBy(&in.UserID),
		role1.WithRole(&in.RoleName),
		role1.WithDescription(&in.Description),
		role1.WithDefault(&in.Default),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAppRole",
			"In", in,
			"Error", err,
		)
		return &npool.CreateAppRoleResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateRole(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAppRole",
			"In", in,
			"Error", err,
		)
		return &npool.CreateAppRoleResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateAppRoleResponse{
		Info: info,
	}, nil
}
