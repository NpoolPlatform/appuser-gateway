package app

import (
	"context"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	tracer "github.com/NpoolPlatform/appuser-middleware/pkg/tracer/app"
	"google.golang.org/grpc/codes"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/app"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateApp(ctx context.Context, in *app.CreateAppRequest) (*app.CreateAppResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateAppV2")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = tracer.Trace(span, in.GetInfo())

	err = validate(ctx, in.Info)
	if err != nil {
		return nil, err
	}

	span = commontracer.TraceInvoker(span, "app", "middleware", "CreateApp")

	resp, err := appmwcli.CreateApp(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("CreateApp", "err", err)
		return &app.CreateAppResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &app.CreateAppResponse{
		Info: resp,
	}, nil
}
