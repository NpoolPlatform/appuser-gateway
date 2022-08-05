//nolint:nolintlint,dupl
package role

import (
	"context"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	tracer "github.com/NpoolPlatform/appuser-manager/pkg/tracer/approle"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	approlemgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approle"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/approle"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateRole(ctx context.Context, in *approle.CreateRoleRequest) (*approle.CreateRoleResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateRole")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = tracer.Trace(span, in.GetInfo())

	err = validate(in.GetInfo(), in.GetUserID())
	if err != nil {
		logger.Sugar().Errorw("CreateRole", "err", err)
		return nil, err
	}

	span = commontracer.TraceInvoker(span, "role", "manager", "CreateAppRole")

	resp, err := approlemgrcli.CreateAppRole(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("CreateRole", "err", err)
		return &approle.CreateRoleResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &approle.CreateRoleResponse{
		Info: resp,
	}, nil
}

func (s *Server) CreateAppRole(ctx context.Context, in *approle.CreateAppRoleRequest) (*approle.CreateAppRoleResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateAppRole")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.SetAttributes(attribute.String("TargetAppID", in.GetTargetAppID()))
	span = tracer.Trace(span, in.GetInfo())

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorw("CreateAppRole", "TargetAppID", in.GetTargetAppID(), "err", err)
		return &approle.CreateAppRoleResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	appID := in.GetTargetAppID()
	info := in.GetInfo()
	info.AppID = &appID

	err = validate(in.GetInfo(), in.GetUserID())
	if err != nil {
		logger.Sugar().Errorw("CreateAppRole", "err", err)
		return &approle.CreateAppRoleResponse{}, err
	}

	span = commontracer.TraceInvoker(span, "role", "manager", "CreateAppRole")

	resp, err := approlemgrcli.CreateAppRole(ctx, info)
	if err != nil {
		logger.Sugar().Errorw("CreateAppRole", "err", err)
		return &approle.CreateAppRoleResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &approle.CreateAppRoleResponse{
		Info: resp,
	}, nil
}
