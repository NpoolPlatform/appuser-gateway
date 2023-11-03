package appsubscribe

import (
	"context"

	appsubscribe1 "github.com/NpoolPlatform/appuser-gateway/pkg/subscriber/app/subscribe"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/subscriber/app/subscribe"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateAppSubscribe(ctx context.Context, in *npool.CreateAppSubscribeRequest) (*npool.CreateAppSubscribeResponse, error) {
	handler, err := appsubscribe1.NewHandler(
		ctx,
		appsubscribe1.WithAppID(&in.TargetAppID, true),
		appsubscribe1.WithSubscribeAppID(&in.SubscribeAppID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAppSubscribe",
			"In", in,
			"Error", err,
		)
		return &npool.CreateAppSubscribeResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateAppSubscribe(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAppSubscribe",
			"In", in,
			"Error", err,
		)
		return &npool.CreateAppSubscribeResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateAppSubscribeResponse{
		Info: info,
	}, nil
}
