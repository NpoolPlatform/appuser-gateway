package app

import (
	"context"

	app1 "github.com/NpoolPlatform/appuser-gateway/pkg/app"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/app"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) BanApp(ctx context.Context, in *npool.BanAppRequest) (*npool.BanAppResponse, error) {
	handler, err := app1.NewHandler(
		ctx,
		app1.WithID(&in.ID),
		app1.WithBanned(&in.Banned),
		app1.WithBanMessage(&in.BanMessage),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"BanApp",
			"In", in,
			"Error", err,
		)
		return &npool.BanAppResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateApp(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"BanApp",
			"In", in,
			"Error", err,
		)
		return &npool.BanAppResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.BanAppResponse{
		Info: info,
	}, nil
}
