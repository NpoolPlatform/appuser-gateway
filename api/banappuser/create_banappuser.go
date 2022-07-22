package banappuser

import (
	"context"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw/banappuser"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateBanAppUser(ctx context.Context, in *banappuser.CreateBanAppUserRequest) (*banappuser.CreateBanAppUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateExtra")
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

	span.AddEvent("call grpc CreateBanAppUserV2")
	resp, err := grpc.CreateBanAppUserV2(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("fail create ban app user: %v", err)
		return &banappuser.CreateBanAppUserResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &banappuser.CreateBanAppUserResponse{
		Info: resp,
	}, nil
}
