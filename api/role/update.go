//nolint:dupl
package role

import (
	"context"

	role1 "github.com/NpoolPlatform/appuser-gateway/pkg/role"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/role"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateRole(ctx context.Context, in *npool.UpdateRoleRequest) (*npool.UpdateRoleResponse, error) {
	handler, err := role1.NewHandler(
		ctx,
		role1.WithID(&in.ID, true),
		role1.WithEntID(&in.EntID, true),
		role1.WithAppID(&in.AppID, true),
		role1.WithRole(in.RoleName, false),
		role1.WithDescription(in.Description, false),
		role1.WithDefault(in.Default, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateRole",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateRoleResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateRole(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateRole",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateRoleResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateRoleResponse{
		Info: info,
	}, nil
}

func (s *Server) UpdateAppRole(ctx context.Context, in *npool.UpdateAppRoleRequest) (*npool.UpdateAppRoleResponse, error) {
	handler, err := role1.NewHandler(
		ctx,
		role1.WithID(&in.ID, true),
		role1.WithEntID(&in.EntID, true),
		role1.WithAppID(&in.TargetAppID, true),
		role1.WithRole(in.RoleName, false),
		role1.WithDescription(in.Description, false),
		role1.WithDefault(in.Default, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAppRole",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateAppRoleResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateRole(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAppRole",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateAppRoleResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateAppRoleResponse{
		Info: info,
	}, nil
}
