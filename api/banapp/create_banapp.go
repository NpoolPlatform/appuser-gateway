//nolint:nolintlint,dupl
package banapp

import (
	"context"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	"github.com/NpoolPlatform/message/npool/appusergw/banapp"
	appcrud "github.com/NpoolPlatform/message/npool/appusermgrv2/banapp"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateBanApp(ctx context.Context, in *banapp.CreateBanAppRequest) (*banapp.CreateBanAppResponse, error) {
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

	span.AddEvent("call grpc ExistBanAppCondsV2")
	exist, err := grpc.ExistBanAppCondsV2(ctx, &appcrud.Conds{
		AppID: &npool.StringVal{
			Value: in.GetInfo().GetAppID(),
			Op:    cruder.EQ,
		}})
	if err != nil {
		logger.Sugar().Errorw("fail check ban app: %v", err)
		return &banapp.CreateBanAppResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}
	if exist {
		logger.Sugar().Errorw("ban app already exists")
		return &banapp.CreateBanAppResponse{}, status.Error(npool.ErrAlreadyExists, appusergw.ErrMsgAppAlreadyExists)
	}

	span.AddEvent("call grpc CreateBanAppV2")
	resp, err := grpc.CreateBanAppV2(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("fail create ban app: %v", err)
		return &banapp.CreateBanAppResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &banapp.CreateBanAppResponse{
		Info: resp,
	}, nil
}
