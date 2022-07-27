//nolint:nolintlint,dupl
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

func (s *Server) UpdateBanAppUser(ctx context.Context,
	in *banappuser.UpdateBanAppUserRequest) (*banappuser.UpdateBanAppUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateBanAppUser")
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
		return &banappuser.UpdateBanAppUserResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgBanAppIDInvalid)
	}

	span.AddEvent("call grpc UpdateBanAppUserV2")
	resp, err := grpc.UpdateBanAppUserV2(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Error("fail update ban app user:$v", err)
		return &banappuser.UpdateBanAppUserResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &banappuser.UpdateBanAppUserResponse{
		Info: resp,
	}, nil
}
