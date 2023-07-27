package oauth

import (
	"context"

	oauth1 "github.com/NpoolPlatform/appuser-gateway/pkg/authing/oauth"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/authing/oauth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetOAuthLoginList(ctx context.Context, in *npool.GetOAuthLoginListRequest) (*npool.GetOAuthLoginListResponse, error) {
	handler, err := oauth1.NewHandler(
		ctx,
		oauth1.WithAppID(in.AppID),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetOAuthLoginList",
			"In", in,
			"Error", err,
		)
		return &npool.GetOAuthLoginListResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	infoLists, err := handler.GetOAuthLoginList(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetOAuthLoginList",
			"In", in,
			"Error", err,
		)
		return &npool.GetOAuthLoginListResponse{}, status.Error(codes.Aborted, err.Error())
	}

	return &npool.GetOAuthLoginListResponse{
		OAuthLoginInfoLists: infoLists,
	}, nil
}

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

	oauthLoginURL, err := handler.GetOAuthURL(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetOAuthLoginList",
			"In", in,
			"Error", err,
		)
		return &npool.GetOAuthLoginURLResponse{}, status.Error(codes.Aborted, err.Error())
	}

	// authURI := "https://github.com/login/oauth/authorize"
	// clientID := "25881c93d384676c0473"
	// redirectURI := "http://localhost:8080/oauth/callback"
	// responseType := "code"
	// state := uuid.NewString()
	// redirectURL := fmt.Sprintf(
	// 	"%s?client_id=%s&redirect_uri=%s&response_type=%s&state=%s",
	// 	authURI, clientID, redirectURI, responseType, state,
	// )
	return &npool.GetOAuthLoginURLResponse{
		OAuthLoginURL: oauthLoginURL,
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
