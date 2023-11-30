package ga

import (
	"context"

	ga1 "github.com/NpoolPlatform/appuser-gateway/pkg/ga"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/ga"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) SetupGoogleAuth(ctx context.Context, in *npool.SetupGoogleAuthRequest) (*npool.SetupGoogleAuthResponse, error) {
	handler, err := ga1.NewHandler(
		ctx,
		ga1.WithAppID(&in.AppID, true),
		ga1.WithUserID(&in.UserID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"SetupGoogleAuth",
			"In", in,
			"Error", err,
		)
		return &npool.SetupGoogleAuthResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.SetupGoogleAuth(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"SetupGoogleAuth",
			"In", in,
			"Error", err,
		)
		return &npool.SetupGoogleAuthResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.SetupGoogleAuthResponse{
		Info: info,
	}, nil
}
