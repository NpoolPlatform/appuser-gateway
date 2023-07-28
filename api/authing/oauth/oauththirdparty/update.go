package oauththirdparty

import (
	"context"

	oauth1 "github.com/NpoolPlatform/appuser-gateway/pkg/authing/oauth/oauththirdparty"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/authing/oauth/oauththirdparty"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateOAuthThirdParty(ctx context.Context, in *npool.UpdateOAuthThirdPartyRequest) (*npool.UpdateOAuthThirdPartyResponse, error) {
	handler, err := oauth1.NewHandler(
		ctx,
		oauth1.WithID(&in.ID),
		oauth1.WithClientName(in.ClientName),
		oauth1.WithClientTag(in.ClientTag),
		oauth1.WithClientLogoURL(in.ClientLogoURL),
		oauth1.WithClientOAuthURL(in.ClientOAuthURL),
		oauth1.WithResponseType(in.ResponseType),
		oauth1.WithScope(in.Scope),
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
