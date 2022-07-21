package appusersecret

import (
	"context"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw/appusersecret"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateSecret(ctx context.Context, in *appusersecret.CreateSecretRequest) (*appusersecret.CreateSecretResponse, error) {
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

	span.AddEvent("call grpc CreateAppUserSecretV2")
	resp, err := grpc.CreateAppUserSecretV2(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("fail create app user secret: %v", err)
		return &appusersecret.CreateSecretResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}
	return &appusersecret.CreateSecretResponse{
		Info: resp,
	}, nil
}
