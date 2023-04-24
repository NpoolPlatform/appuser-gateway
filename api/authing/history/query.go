package history

import (
	"context"

	history1 "github.com/NpoolPlatform/appuser-gateway/pkg/authing/history"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/authing/history"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetAppAuthHistories(
	ctx context.Context,
	in *npool.GetAppAuthHistoriesRequest,
) (
	*npool.GetAppAuthHistoriesResponse,
	error,
) {
	handler, err := history1.NewHandler(
		ctx,
		history1.WithAppID(in.GetTargetAppID()),
		history1.WithOffset(in.GetOffset()),
		history1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppAuthHistories",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppAuthHistoriesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetAuthHistories(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppAuthHistories",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppAuthHistoriesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppAuthHistoriesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
