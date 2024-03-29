package oauththirdparty

import (
	"context"

	oauth1 "github.com/NpoolPlatform/appuser-gateway/pkg/oauth/oauththirdparty"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/oauth/oauththirdparty"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateOAuthThirdParty(ctx context.Context, in *npool.UpdateOAuthThirdPartyRequest) (*npool.UpdateOAuthThirdPartyResponse, error) {
	handler, err := oauth1.NewHandler(
		ctx,
		oauth1.WithID(&in.ID, true),
		oauth1.WithEntID(&in.EntID, true),
		oauth1.WithClientName(in.ClientName, false),
		oauth1.WithClientTag(in.ClientTag, false),
		oauth1.WithClientLogoURL(in.ClientLogoURL, false),
		oauth1.WithClientOAuthURL(in.ClientOAuthURL, false),
		oauth1.WithResponseType(in.ResponseType, false),
		oauth1.WithScope(in.Scope, false),
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
