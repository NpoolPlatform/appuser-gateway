//nolint:nolintlint,dupl
package user

import (
	"context"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	tracer "github.com/NpoolPlatform/appuser-middleware/pkg/tracer/user"
	"google.golang.org/grpc/codes"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateUser(ctx context.Context, in *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateUser")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = tracer.Trace(span, in.GetInfo())
	span = commontracer.TraceInvoker(span, "role", "middleware", "UpdateUser")

	resp, err := usermwcli.UpdateUser(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("UpdateUser", "err", err)
		return &user.UpdateUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &user.UpdateUserResponse{
		Info: resp,
	}, nil
}
