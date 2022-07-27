//nolint:nolintlint,dupl
package thirdparty

import (
	"context"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw/thirdparty"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateThirdParty(ctx context.Context,
	in *thirdparty.CreateThirdPartyRequest) (*thirdparty.CreateThirdPartyResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateThirdParty")
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

	span.AddEvent("call grpc CreateAppUserThirdPartyV2")
	resp, err := grpc.CreateAppUserThirdPartyV2(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("fail create app user third party: %v", err)
		return &thirdparty.CreateThirdPartyResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}
	return &thirdparty.CreateThirdPartyResponse{
		Info: resp,
	}, nil
}
