package appoauththirdparty

import (
	"context"

	oauth1 "github.com/NpoolPlatform/appuser-gateway/pkg/oauth/appoauththirdparty"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/oauth/appoauththirdparty"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateOAuthThirdParty(ctx context.Context, in *npool.CreateOAuthThirdPartyRequest) (resp *npool.CreateOAuthThirdPartyResponse, err error) {
	handler, err := oauth1.NewHandler(
		ctx,
		oauth1.WithAppID(in.GetTargetAppID()),
		oauth1.WithClientID(&in.ClientID),
		oauth1.WithClientSecret(&in.ClientSecret),
		oauth1.WithCallbackURL(&in.CallbackURL),
		oauth1.WithThirdPartyID(&in.ThirdPartyID),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateOAuthThirdParty",
			"In", in,
			"Error", err,
		)
		return &npool.CreateOAuthThirdPartyResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateOAuthThirdParty(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateOAuthThirdParty",
			"In", in,
			"Error", err,
		)
		return &npool.CreateOAuthThirdPartyResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateOAuthThirdPartyResponse{
		Info: info,
	}, nil
}
