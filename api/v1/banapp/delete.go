//nolint:nolintlint,dupl
package banapp

import (
	"context"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/banapp"
	"google.golang.org/grpc/codes"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	banappmgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/banapp"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteBanApp(ctx context.Context, in *banapp.DeleteBanAppRequest) (*banapp.DeleteBanAppResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "DeleteBanApp")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceID(span, in.GetID())

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Errorw("DeleteBanApp", "ID", in.GetID(), "err", err)
		return &banapp.DeleteBanAppResponse{}, status.Error(codes.InvalidArgument, "ID is invalid")
	}

	span = commontracer.TraceInvoker(span, "banapp", "manager", "DeleteBanApp")

	resp, err := banappmgrcli.DeleteBanApp(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("DeleteBanApp", "err", err)
		return &banapp.DeleteBanAppResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &banapp.DeleteBanAppResponse{
		Info: resp,
	}, nil
}
