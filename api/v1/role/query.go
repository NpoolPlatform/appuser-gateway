//nolint:dupl
package role

import (
	"context"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	approlemgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approle"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/approle"
	approlepb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approle"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetRoles(ctx context.Context, in *approle.GetRolesRequest) (*approle.GetRolesResponse, error) {
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
		return &approle.GetRolesResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	span = commontracer.TraceInvoker(span, "role", "manager", "GetAppRoles")

	resp, _, err := approlemgrcli.GetAppRoles(ctx, &approlepb.Conds{
		AppID: &npool.StringVal{
			Value: in.GetAppID(),
			Op:    cruder.EQ,
		},
	}, in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetRoles", "err", err)
		return &approle.GetRolesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &approle.GetRolesResponse{
		Infos: resp,
	}, nil
}

func (s *Server) GetAppRoles(ctx context.Context, in *approle.GetAppRolesRequest) (*approle.GetAppRolesResponse, error) {
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
		return &approle.GetAppRolesResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	span = commontracer.TraceInvoker(span, "role", "manager", "GetAppRoles")

	resp, _, err := approlemgrcli.GetAppRoles(ctx, &approlepb.Conds{
		AppID: &npool.StringVal{
			Value: in.GetTargetAppID(),
			Op:    cruder.EQ,
		},
	}, in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetAppRoles", "err", err)
		return &approle.GetAppRolesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &approle.GetAppRolesResponse{
		Infos: resp,
	}, nil
}
