package app

import (
	"context"

	app1 "github.com/NpoolPlatform/appuser-gateway/pkg/app"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/app"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteApp(ctx context.Context, in *app.DeleteAppRequest) (*app.DeleteAppResponse, error) {
	handler, err := app1.NewHandler(
		ctx,
		app1.WithID(&in.ID),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteApp",
			"In", in,
			"Error", err,
		)
		return &app.DeleteAppResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	info, err := handler.DeleteApp(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteApp",
			"In", in,
			"Error", err,
		)
		return &app.DeleteAppResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &app.DeleteAppResponse{
		Info: info,
	}, nil
}
