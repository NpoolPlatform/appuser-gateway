package auth

import (
	"context"

	auth1 "github.com/NpoolPlatform/appuser-gateway/pkg/authing/auth"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/authing/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetAppAuths(ctx context.Context, in *npool.GetAppAuthsRequest) (resp *npool.GetAppAuthsResponse, err error) {
	handler, err := auth1.NewHandler(
		ctx,
		auth1.WithAppID(in.GetTargetAppID()),
		auth1.WithOffset(in.GetOffset()),
		auth1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppAuths",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppAuthsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetAuths(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppAuths",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppAuthsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppAuthsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
