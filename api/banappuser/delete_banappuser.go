package banappuser

import (
	"context"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	"github.com/NpoolPlatform/message/npool/appusergw/banappuser"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteBanAppUser(ctx context.Context, in *banappuser.DeleteBanAppUserRequest) (*banappuser.DeleteBanAppUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "DeleteBanApp")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Error("ID is invalid")
		return &banappuser.DeleteBanAppUserResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgAppIDInvalid)
	}

	span.AddEvent("call grpc DeleteBanAppV2")
	resp, err := grpc.DeleteBanAppUserV2(ctx, in.GetID())
	if err != nil {
		return &banappuser.DeleteBanAppUserResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &banappuser.DeleteBanAppUserResponse{
		Info: resp,
	}, nil
}
