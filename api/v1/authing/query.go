package authing

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/authing"

	authing1 "github.com/NpoolPlatform/appuser-gateway/pkg/authing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) Authenticate(ctx context.Context, in *npool.AuthenticateRequest) (*npool.AuthenticateResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("Authenticate", "AppID", in.GetAppID())
		return &npool.AuthenticateResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}
	if in.UserID != nil {
		if _, err := uuid.Parse(in.GetUserID()); err != nil {
			logger.Sugar().Errorw("Authenticate", "UserID", in.GetUserID())
			return &npool.AuthenticateResponse{}, status.Error(codes.InvalidArgument, "UserID is invalid")
		}
	}
	if in.Token != nil && in.GetToken() == "" {
		logger.Sugar().Errorw("Authenticate", "Token", in.GetToken())
		return &npool.AuthenticateResponse{}, status.Error(codes.InvalidArgument, "Token is invalid")
	}
	if in.GetResource() == "" {
		logger.Sugar().Errorw("Authenticate", "Resource", in.GetResource())
		return &npool.AuthenticateResponse{}, status.Error(codes.InvalidArgument, "Resource is invalid")
	}
	if in.GetMethod() == "" {
		logger.Sugar().Errorw("Authenticate", "Method", in.GetMethod())
		return &npool.AuthenticateResponse{}, status.Error(codes.InvalidArgument, "Method is invalid")
	}

	allowed, err := authing1.Authenticate(ctx, in.GetAppID(), in.UserID, in.Token, in.GetResource(), in.GetMethod())
	if err != nil {
		logger.Sugar().Errorw("Authenticate", "error", err)
		return &npool.AuthenticateResponse{}, status.Error(codes.Internal, "fail authenticate")
	}

	return &npool.AuthenticateResponse{
		Info: allowed,
	}, nil
}

func (s *Server) GetAppAuths(ctx context.Context, in *npool.GetAppAuthsRequest) (*npool.GetAppAuthsResponse, error) {
	return nil, nil
}

func (s *Server) GetAppHistories(ctx context.Context, in *npool.GetAppHistoriesRequest) (*npool.GetAppHistoriesResponse, error) {
	return nil, nil
}