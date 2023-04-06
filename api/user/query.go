//nolint:dupl
package user

import (
	"context"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"

	constant1 "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"

	mw "github.com/NpoolPlatform/appuser-gateway/pkg/user"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetUsers(ctx context.Context, in *user.GetUsersRequest) (*user.GetUsersResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateUser")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.SetAttributes(attribute.String("AppID", in.GetAppID()))
	commontracer.TraceOffsetLimit(span, int(in.GetOffset()), int(in.GetLimit()))

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("GetUsers", "AppID", in.GetAppID(), "err", err)
		return &user.GetUsersResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	span = commontracer.TraceInvoker(span, "role", "middleware", "CreateUser")

	limit := constant1.DefaultRowLimit
	if in.GetLimit() > 0 {
		limit = in.GetLimit()
	}

	infos, total, err := mw.GetUsers(ctx, in.GetAppID(), in.GetOffset(), limit)
	if err != nil {
		logger.Sugar().Errorw("GetUsers", "err", err)
		return &user.GetUsersResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &user.GetUsersResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppUsers(ctx context.Context, in *user.GetAppUsersRequest) (*user.GetAppUsersResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateUser")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.SetAttributes(attribute.String("TargetAppID", in.GetTargetAppID()))
	commontracer.TraceOffsetLimit(span, int(in.GetOffset()), int(in.GetLimit()))

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorw("GetAppUsers", "TargetAppID", in.GetTargetAppID(), "err", err)
		return &user.GetAppUsersResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	span = commontracer.TraceInvoker(span, "role", "middleware", "CreateUser")

	limit := constant1.DefaultRowLimit
	if in.GetLimit() > 0 {
		limit = in.GetLimit()
	}

	infos, total, err := mw.GetUsers(ctx, in.GetTargetAppID(), in.GetOffset(), limit)
	if err != nil {
		logger.Sugar().Errorw("GetAppUsers", "err", err)
		return &user.GetAppUsersResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &user.GetAppUsersResponse{
		Infos: infos,
		Total: total,
	}, nil
}
