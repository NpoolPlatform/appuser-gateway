package admin

import (
	"context"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	mw "github.com/NpoolPlatform/appuser-gateway/pkg/middleware/admin"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	"github.com/NpoolPlatform/message/npool/appusergw/admin"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetAdminApps(ctx context.Context, in *admin.GetAdminAppsRequest) (*admin.GetAdminAppsResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetExtra")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.AddEvent("call middleware GetAdminApps")
	resp, err := mw.GetAdminApps(ctx)
	if err != nil {
		logger.Sugar().Errorw("fail Get admin apps: %v", err)
		return &admin.GetAdminAppsResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	if len(resp) == 0 {
		return nil, status.Error(npool.ErrService, appusergw.ErrMsgAdminAppNotFound)
	}

	return &admin.GetAdminAppsResponse{
		Infos: resp,
	}, nil
}
