package appoauththirdparty

import (
	"context"

	oauth1 "github.com/NpoolPlatform/appuser-gateway/pkg/oauth/appoauththirdparty"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/oauth/appoauththirdparty"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateOAuthThirdParty(ctx context.Context, in *npool.UpdateOAuthThirdPartyRequest) (*npool.UpdateOAuthThirdPartyResponse, error) {
	handler, err := oauth1.NewHandler(
		ctx,
		oauth1.WithID(&in.ID, true),
		oauth1.WithEntID(&in.EntID, true),
		oauth1.WithAppID(&in.TargetAppID, true),
		oauth1.WithClientID(in.ClientID, false),
		oauth1.WithClientSecret(in.ClientSecret, false),
		oauth1.WithCallbackURL(in.CallbackURL, false),
		oauth1.WithThirdPartyID(in.ThirdPartyID, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateOAuthThirdParty",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateOAuthThirdPartyResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateOAuthThirdParty(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateOAuthThirdParty",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateOAuthThirdPartyResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateOAuthThirdPartyResponse{
		Info: info,
	}, nil
}
