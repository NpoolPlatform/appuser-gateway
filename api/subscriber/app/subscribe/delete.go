package appsubscribe

import (
	"context"

	appsubscribe1 "github.com/NpoolPlatform/appuser-gateway/pkg/subscriber/app/subscribe"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/subscriber/app/subscribe"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteAppSubscribe(ctx context.Context, in *npool.DeleteAppSubscribeRequest) (*npool.DeleteAppSubscribeResponse, error) {
	handler, err := appsubscribe1.NewHandler(
		ctx,
		appsubscribe1.WithID(&in.ID),
		appsubscribe1.WithAppID(in.GetTargetAppID()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteAppSubscribe",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteAppSubscribeResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteAppSubscribe(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteAppSubscribe",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteAppSubscribeResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteAppSubscribeResponse{
		Info: info,
	}, nil
}
