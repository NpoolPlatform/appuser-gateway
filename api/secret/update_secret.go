//nolint:nolintlint,dupl
package secret

import (
	"context"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	"github.com/NpoolPlatform/message/npool/appusergw/appusersecret"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateSecret(ctx context.Context, in *appusersecret.UpdateSecretRequest) (*appusersecret.UpdateSecretResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateSecret")
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
		return &appusersecret.UpdateSecretResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgBanAppIDInvalid)
	}

	span.AddEvent("call grpc UpdateAppUserSecretV2")
	resp, err := grpc.UpdateAppUserSecretV2(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Error("fail update secret:%v", err)
		return &appusersecret.UpdateSecretResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &appusersecret.UpdateSecretResponse{
		Info: resp,
	}, nil
}
