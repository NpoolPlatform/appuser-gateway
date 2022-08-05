//nolint:nolintlint,dupl
package roleuser

import (
	"context"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	approleusermgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approleuser"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/approleuser"
	approleusercrud "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approleuser"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetRoleUsers(ctx context.Context, in *approleuser.GetRoleUsersRequest) (*approleuser.GetRoleUsersResponse, error) {
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
		return &approleuser.GetRoleUsersResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	span = commontracer.TraceInvoker(span, "role", "manager", "GetAppRoleUsers")

	resp, _, err := approleusermgrcli.GetAppRoleUsers(ctx, &approleusercrud.Conds{
		AppID: &npool.StringVal{
			Value: in.GetAppID(),
			Op:    cruder.EQ,
		},
		RoleID: &npool.StringVal{
			Value: in.GetRoleID(),
			Op:    cruder.EQ,
		},
	}, in.GetLimit(), in.GetOffset())
	if err != nil {
		logger.Sugar().Errorw("GetRoleUsers", "err", err)
		return &approleuser.GetRoleUsersResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	return &approleuser.GetRoleUsersResponse{
		Infos: resp,
	}, nil
}

func (s *Server) GetAppRoleUsers(ctx context.Context,
	in *approleuser.GetAppRoleUsersRequest) (*approleuser.GetAppRoleUsersResponse, error) {
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
		return &approleuser.GetAppRoleUsersResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	span = commontracer.TraceInvoker(span, "role", "manager", "GetAppRoleUsers")

	resp, _, err := approleusermgrcli.GetAppRoleUsers(ctx, &approleusercrud.Conds{
		AppID: &npool.StringVal{
			Value: in.GetTargetAppID(),
			Op:    cruder.EQ,
		},
		RoleID: &npool.StringVal{
			Value: in.GetRoleID(),
			Op:    cruder.EQ,
		},
	}, in.GetLimit(), in.GetOffset())
	if err != nil {
		logger.Sugar().Errorw("GetAppRoleUsers", "err", err)
		return &approleuser.GetAppRoleUsersResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &approleuser.GetAppRoleUsersResponse{
		Infos: resp,
	}, nil
}
