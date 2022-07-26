package role

import (
	"context"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	"github.com/NpoolPlatform/message/npool/appusergw/approle"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateRole(ctx context.Context, in *approle.UpdateRoleRequest) (*approle.UpdateRoleResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateRole")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetInfo().GetID()); err != nil {
		logger.Sugar().Error("ID is invalid")
		return &approle.UpdateRoleResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgBanAppIDInvalid)
	}

	span.AddEvent("call grpc UpdateAppRoleV2")
	resp, err := grpc.UpdateAppRoleV2(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Error("fail update app role:%v", err)
		return &approle.UpdateRoleResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &approle.UpdateRoleResponse{
		Info: resp,
	}, nil
}
