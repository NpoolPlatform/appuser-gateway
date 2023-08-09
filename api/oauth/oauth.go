package oauth

import (
	"context"

	oauth1 "github.com/NpoolPlatform/appuser-gateway/pkg/oauth"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/oauth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetOAuthLoginURL(ctx context.Context, in *npool.GetOAuthLoginURLRequest) (*npool.GetOAuthLoginURLResponse, error) {
	handler, err := oauth1.NewHandler(
		ctx,
		oauth1.WithAppID(in.AppID),
		oauth1.WithClientName(&in.ClientName),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetOAuthLoginList",
			"In", in,
			"Error", err,
		)
		return &npool.GetOAuthLoginURLResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.GetOAuthURL(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetOAuthLoginList",
			"In", in,
			"Error", err,
		)
		return &npool.GetOAuthLoginURLResponse{}, status.Error(codes.Aborted, err.Error())
	}

	return &npool.GetOAuthLoginURLResponse{
		Info: info,
	}, nil
}

func (s *Server) OAuthLogin(ctx context.Context, in *npool.OAuthLoginRequest) (*npool.OAuthLoginResponse, error) {
	handler, err := oauth1.NewHandler(
		ctx,
		oauth1.WithAppID(in.AppID),
		oauth1.WithCode(&in.Code),
		oauth1.WithState(&in.State),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"Login",
			"In", in,
			"Error", err,
		)
		return &npool.OAuthLoginResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.OAuthLogin(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"Login",
			"In", in,
			"Error", err,
		)
		return &npool.OAuthLoginResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &npool.OAuthLoginResponse{
		Info: info,
	}, nil
}
