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

func (s *Server) GetRoles(ctx context.Context, in *npool.GetRolesRequest) (*npool.GetRolesResponse, error) {
	handler, err := role1.NewHandler(
		ctx,
		role1.WithAppID(in.GetAppID()),
		role1.WithOffset(in.GetOffset()),
		role1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetRoles",
			"In", in,
			"Error", err,
		)
		return &npool.GetRolesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetRoles(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetRoles",
			"In", in,
			"Error", err,
		)
		return &npool.GetRolesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetRolesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppRoles(ctx context.Context, in *npool.GetAppRolesRequest) (*npool.GetAppRolesResponse, error) {
	handler, err := role1.NewHandler(
		ctx,
		role1.WithAppID(in.GetTargetAppID()),
		role1.WithOffset(in.GetOffset()),
		role1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppRoles",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppRolesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetRoles(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppRoles",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppRolesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppRolesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
