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

func (s *Server) GetUsers(ctx context.Context, in *npool.GetUsersRequest) (*npool.GetUsersResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(in.GetAppID()),
		user1.WithOffset(in.GetOffset()),
		user1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetUsers",
			"In", in,
			"Error", err,
		)
		return &npool.GetUsersResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetUsers(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetUsers",
			"In", in,
			"Error", err,
		)
		return &npool.GetUsersResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetUsersResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppUsers(ctx context.Context, in *npool.GetAppUsersRequest) (*npool.GetAppUsersResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(in.GetTargetAppID()),
		user1.WithOffset(in.GetOffset()),
		user1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppUsers",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppUsersResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetUsers(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppUsers",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppUsersResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppUsersResponse{
		Infos: infos,
		Total: total,
	}, nil
}
