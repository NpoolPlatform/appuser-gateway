//nolint:dupl
package role

import (
	"context"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/role"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetRoles(ctx context.Context, in *role.GetRolesRequest) (*role.GetRolesResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetRoles")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.SetAttributes(attribute.String("AppID", in.GetAppID()))

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("GetRoles", "AppID", in.GetAppID(), "err", err)
		return &role.GetRolesResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	span = commontracer.TraceInvoker(span, "role", "manager", "GetAppRoles")

	infos, total, err := rolemwcli.GetRoles(ctx, in.GetAppID(), in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetRoles", "err", err)
		return &role.GetRolesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &role.GetRolesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppRoles(ctx context.Context, in *role.GetAppRolesRequest) (*role.GetAppRolesResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppRoles")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.SetAttributes(attribute.String("TargetAppID", in.GetTargetAppID()))

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorw("GetAppRoles", "TargetAppID", in.GetTargetAppID(), "err", err)
		return &role.GetAppRolesResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	span = commontracer.TraceInvoker(span, "role", "manager", "GetAppRoles")

	infos, total, err := rolemwcli.GetRoles(ctx, in.GetTargetAppID(), in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetRoles", "err", err)
		return &role.GetAppRolesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &role.GetAppRolesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetRoleUsers(ctx context.Context, in *role.GetRoleUsersRequest) (*role.GetRoleUsersResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetRoleUsers")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.SetAttributes(attribute.String("AppID", in.GetAppID()))
	span.SetAttributes(attribute.String("RoleID", in.GetRoleID()))
	commontracer.TraceOffsetLimit(span, int(in.GetOffset()), int(in.GetLimit()))

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("GetRoleUsers", "AppID", in.GetAppID(), "err", err)
		return &role.GetRoleUsersResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	span = commontracer.TraceInvoker(span, "role", "manager", "GetAppRoleUsers")

	infos, total, err := rolemwcli.GetRoleUsers(ctx, in.GetAppID(), in.GetRoleID(), in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetRoleUsers", "err", err)
		return &role.GetRoleUsersResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &role.GetRoleUsersResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppRoleUsers(ctx context.Context,
	in *role.GetAppRoleUsersRequest) (*role.GetAppRoleUsersResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppRoleUsers")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.SetAttributes(attribute.String("TargetAppID", in.GetTargetAppID()))
	span.SetAttributes(attribute.String("RoleID", in.GetRoleID()))
	commontracer.TraceOffsetLimit(span, int(in.GetOffset()), int(in.GetLimit()))

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorw("GetAppRoleUsers", "TargetAppID", in.GetTargetAppID(), "err", err)
		return &role.GetAppRoleUsersResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	span = commontracer.TraceInvoker(span, "role", "manager", "GetAppRoleUsers")

	infos, total, err := rolemwcli.GetRoleUsers(ctx, in.GetTargetAppID(), in.GetRoleID(), in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetRoleUsers", "err", err)
		return &role.GetAppRoleUsersResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &role.GetAppRoleUsersResponse{
		Infos: infos,
		Total: total,
	}, nil
}
