package subscriber

import (
	"context"

	subscriber1 "github.com/NpoolPlatform/appuser-gateway/pkg/subscriber"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/subscriber"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateSubscriber(ctx context.Context, in *npool.CreateSubscriberRequest) (*npool.CreateSubscriberResponse, error) {
	handler, err := subscriber1.NewHandler(
		ctx,
		subscriber1.WithAppID(in.GetAppID()),
		subscriber1.WithSubscribeAppID(in.SubscribeAppID),
		subscriber1.WithEmailAddress(in.GetEmailAddress()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateSubscriber",
			"In", in,
			"Error", err,
		)
		return &npool.CreateSubscriberResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateSubscriber(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateSubscriber",
			"In", in,
			"Error", err,
		)
		return &npool.CreateSubscriberResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateSubscriberResponse{
		Info: info,
	}, nil
}
