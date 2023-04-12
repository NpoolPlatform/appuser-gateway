package role

import (
	"context"

	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	tracer "github.com/NpoolPlatform/appuser-manager/pkg/tracer/approle"
	"google.golang.org/grpc/codes"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	approlemgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approle"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/role"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateRole(ctx context.Context, in *role.UpdateRoleRequest) (*role.UpdateRoleResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateRole")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = tracer.Trace(span, in.GetInfo())

	if _, err := uuid.Parse(in.GetInfo().GetID()); err != nil {
		logger.Sugar().Errorw("GetAppRoles", "ID", in.GetInfo().GetID(), "err", err)
		return &role.UpdateRoleResponse{}, status.Error(codes.InvalidArgument, "ID is invalid")
	}

	span = commontracer.TraceInvoker(span, "role", "manager", "UpdateAppRole")

	appRole, err := approlemgrcli.UpdateAppRole(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("GetAppRoles", "err", err)
		return &role.UpdateRoleResponse{}, status.Error(codes.Internal, err.Error())
	}

	info, err := rolemwcli.GetRole(ctx, appRole.ID)
	if err != nil {
		logger.Sugar().Errorw("CreateRole", "err", err)
		return &role.UpdateRoleResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &role.UpdateRoleResponse{
		Info: info,
	}, nil
}
