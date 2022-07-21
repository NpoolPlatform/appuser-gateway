package banapp

import (
	"context"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw/banapp"
	banappcrud "github.com/NpoolPlatform/message/npool/appusermgrv2/banapp"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateBanApp(ctx context.Context, in *banapp.UpdateBanAppRequest) (*banapp.UpdateBanAppResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateBanApp")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	err = validate(in.GetInfo())
	if err != nil {
		return nil, err
	}

	if _, err := uuid.Parse(in.GetInfo().GetID()); err != nil {
		logger.Sugar().Error("ID is invalid")
		return &banapp.UpdateBanAppResponse{}, status.Error(npool.ErrParams, banapp.ErrMsgBanAppIDInvalid)
	}

	if in.GetInfo().GetMessage() == "" {
		logger.Sugar().Error("Message is empty")
		return &banapp.UpdateBanAppResponse{}, status.Error(npool.ErrParams, banapp.ErrMsgBanAppMessageEmpty)
	}

	span.AddEvent("call grpc ExistBanAppCondsV2")
	resp, err := grpc.UpdateBanAppV2(ctx, &banappcrud.BanAppReq{
		Message: in.GetInfo().Message,
	})
	if err != nil {
		return &banapp.UpdateBanAppResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &banapp.UpdateBanAppResponse{
		Info: resp,
	}, nil
}
