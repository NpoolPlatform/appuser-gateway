package appsubscribe

import (
	"context"

	appsubscribe1 "github.com/NpoolPlatform/appuser-gateway/pkg/subscriber/app/subscribe"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/subscriber/app/subscribe"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetAppSubscribes(ctx context.Context, in *npool.GetAppSubscribesRequest) (*npool.GetAppSubscribesResponse, error) {
	handler, err := appsubscribe1.NewHandler(
		ctx,
		appsubscribe1.WithAppID(in.GetTargetAppID()),
		appsubscribe1.WithOffset(in.GetOffset()),
		appsubscribe1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppSubscribes",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppSubscribesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetAppSubscribes(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppSubscribes",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppSubscribesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppSubscribesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
