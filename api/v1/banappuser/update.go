//nolint:nolintlint,dupl
package banappuser

import (
	"context"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	tracer "github.com/NpoolPlatform/appuser-manager/pkg/tracer/banappuser"
	"google.golang.org/grpc/codes"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	banappusermgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/banappuser"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/banappuser"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateBanAppUser(ctx context.Context,
	in *banappuser.UpdateBanUserRequest) (*banappuser.UpdateBanUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateBanAppUser")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = tracer.Trace(span, in.GetInfo())

	err = validate(in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("UpdateBanAppUser", "err", err)
		return nil, err
	}

	if _, err := uuid.Parse(in.GetInfo().GetID()); err != nil {
		logger.Sugar().Errorw("UpdateBanAppUser", "ID", in.GetInfo().GetID(), "err", err)
		return &banappuser.UpdateBanUserResponse{}, status.Error(codes.InvalidArgument, "ID is invalid")
	}

	span = commontracer.TraceInvoker(span, "banappuser", "manager", "UpdateBanAppUser")

	resp, err := banappusermgrcli.UpdateBanAppUser(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("UpdateBanAppUser", "err", err)
		return &banappuser.UpdateBanUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &banappuser.UpdateBanUserResponse{
		Info: resp,
	}, nil
}
