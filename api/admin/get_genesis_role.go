package admin

import (
	"context"
	"github.com/NpoolPlatform/api-manager/pkg/db/ent"
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

func (s *Server) GetGenesisRole(ctx context.Context, in *admin.GetGenesisRoleRequest) (*admin.GetGenesisRoleResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateExtra")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.AddEvent("call middleware CreateGenesisRole")
	resp, err := mw.GetGenesisRole(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			logger.Sugar().Errorw("genesis role not found: %v", err)
			return &admin.GetGenesisRoleResponse{}, status.Error(npool.ErrService, appusergw.ErrMsgAdminAppNotFound)
		}
		logger.Sugar().Errorw("fail get genesis role : %v", err)
		return &admin.GetGenesisRoleResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &admin.GetGenesisRoleResponse{
		Info: resp,
	}, nil
}
