package banapp

import (
	"context"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	"github.com/NpoolPlatform/message/npool/appusergw/banapp"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetBanApp(ctx context.Context, in *banapp.GetBanAppRequest) (*banapp.GetBanAppResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetBanApp")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Error("ID is invalid")
		return &banapp.GetBanAppResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgIDInvalid)
	}

	span.AddEvent("call grpc GetBanAppV2")
	resp, err := grpc.GetBanAppV2(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Error("fail get ban app")
		return &banapp.GetBanAppResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &banapp.GetBanAppResponse{
		Info: resp,
	}, nil
}
