package ga

import (
	"context"

	ga1 "github.com/NpoolPlatform/appuser-gateway/pkg/ga"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/ga"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) VerifyGoogleAuth(ctx context.Context, in *npool.VerifyGoogleAuthRequest) (*npool.VerifyGoogleAuthResponse, error) {
	handler, err := ga1.NewHandler(
		ctx,
		ga1.WithAppID(&in.AppID, true),
		ga1.WithUserID(&in.UserID, true),
		ga1.WithCode(&in.Code, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"VerifyGoogleAuth",
			"In", in,
			"Error", err,
		)
		return &npool.VerifyGoogleAuthResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.VerifyGoogleAuth(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"VerifyGoogleAuth",
			"In", in,
			"Error", err,
		)
		return &npool.VerifyGoogleAuthResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.VerifyGoogleAuthResponse{
		Info: info,
	}, nil
}
