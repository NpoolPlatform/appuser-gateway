package oauththirdparty

import (
	"context"

	oauth1 "github.com/NpoolPlatform/appuser-gateway/pkg/authing/oauth/oauththirdparty"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/authing/oauth/oauththirdparty"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetOAuthThirdParties(ctx context.Context, in *npool.GetOAuthThirdPartiesRequest) (resp *npool.GetOAuthThirdPartiesResponse, err error) {
	handler, err := oauth1.NewHandler(
		ctx,
		oauth1.WithOffset(in.GetOffset()),
		oauth1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetOAuthThirdParties",
			"In", in,
			"Error", err,
		)
		return &npool.GetOAuthThirdPartiesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetOAuthThirdParties(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetOAuthThirdParties",
			"In", in,
			"Error", err,
		)
		return &npool.GetOAuthThirdPartiesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetOAuthThirdPartiesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
