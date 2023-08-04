package appoauththirdparty

import (
	"context"

	oauth1 "github.com/NpoolPlatform/appuser-gateway/pkg/oauth/appoauththirdparty"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/oauth/appoauththirdparty"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteOAuthThirdParty(ctx context.Context, in *npool.DeleteOAuthThirdPartyRequest) (*npool.DeleteOAuthThirdPartyResponse, error) {
	handler, err := oauth1.NewHandler(
		ctx,
		oauth1.WithID(&in.ID),
		oauth1.WithAppID(in.GetTargetAppID()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteOAuthThirdParty",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteOAuthThirdPartyResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteOAuthThirdParty(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteOAuthThirdParty",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteOAuthThirdPartyResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteOAuthThirdPartyResponse{
		Info: info,
	}, nil
}
