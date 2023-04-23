package auth

import (
	"context"

	auth1 "github.com/NpoolPlatform/appuser-gateway/pkg/authing/auth"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/authing/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteAppAuth(ctx context.Context, in *npool.DeleteAppAuthRequest) (resp *npool.DeleteAppAuthResponse, err error) {
	handler, err := auth1.NewHandler(
		ctx,
		auth1.WithID(&in.ID),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteAppAuth",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteAppAuthResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteAuth(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteAppAuth",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteAppAuthResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteAppAuthResponse{
		Info: info,
	}, nil
}
