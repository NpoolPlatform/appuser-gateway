package subscriber

import (
	"context"

	subscriber1 "github.com/NpoolPlatform/appuser-gateway/pkg/subscriber"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/subscriber"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteSubscriber(ctx context.Context, in *npool.DeleteSubscriberRequest) (*npool.DeleteSubscriberResponse, error) {
	handler, err := subscriber1.NewHandler(
		ctx,
		subscriber1.WithAppID(in.GetAppID()),
		subscriber1.WithEmailAddress(in.GetEmailAddress()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteSubscriber",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteSubscriberResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteSubscriber(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteSubscriber",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteSubscriberResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteSubscriberResponse{
		Info: info,
	}, nil
}
