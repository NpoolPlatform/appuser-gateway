package role

import (
	"context"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	"google.golang.org/grpc/codes"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	approleusermgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approleuser"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/role"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteRoleUser(ctx context.Context, in *role.DeleteRoleUserRequest) (*role.DeleteRoleUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "DeleteRoleUser")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	commontracer.TraceID(span, in.GetID())

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Errorw("DeleteRoleUser", "ID", in.GetID(), "err", err)
		return &role.DeleteRoleUserResponse{}, status.Error(codes.InvalidArgument, "ID is invalid")
	}

	span = commontracer.TraceInvoker(span, "role", "manager", "DeleteAppRoleUser")

	resp, err := approleusermgrcli.DeleteAppRoleUser(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("DeleteRoleUser", "err", err)
		return &role.DeleteRoleUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &role.DeleteRoleUserResponse{
		Info: resp,
	}, nil
}
