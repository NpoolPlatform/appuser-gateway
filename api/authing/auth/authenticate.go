//nolint:dupl
package auth

import (
	"context"

	auth1 "github.com/NpoolPlatform/appuser-gateway/pkg/authing/auth"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/authing/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Authenticate(ctx context.Context, in *npool.AuthenticateRequest) (*npool.AuthenticateResponse, error) {
	handler, err := auth1.NewHandler(
		ctx,
		auth1.WithAppID(in.GetAppID()),
		auth1.WithUserID(in.UserID),
		auth1.WithToken(in.Token),
		auth1.WithResource(in.GetResource()),
		auth1.WithMethod(in.GetMethod()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"Authenticate",
			"In", in,
			"Error", err,
		)
		return &npool.AuthenticateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	allowed, err := handler.Authenticate(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"Authenticate",
			"In", in,
			"Error", err,
		)
		return &npool.AuthenticateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AuthenticateResponse{
		Info: allowed,
	}, nil
}
