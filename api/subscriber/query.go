//nolint:dupl
package subscriber

import (
	"context"

	subscriber1 "github.com/NpoolPlatform/appuser-gateway/pkg/subscriber"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/subscriber"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetSubscriberes(ctx context.Context, in *npool.GetSubscriberesRequest) (*npool.GetSubscriberesResponse, error) {
	handler, err := subscriber1.NewHandler(
		ctx,
		subscriber1.WithAppID(in.GetAppID()),
		subscriber1.WithOffset(in.GetOffset()),
		subscriber1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetSubscriberes",
			"In", in,
			"Error", err,
		)
		return &npool.GetSubscriberesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetSubscriberes(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetSubscriberes",
			"In", in,
			"Error", err,
		)
		return &npool.GetSubscriberesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetSubscriberesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
