//nolint:nolintlint,dupl
package user

import (
	"context"
	"fmt"

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

func (s *Server) Signup(ctx context.Context, in *user.SignupRequest) (*user.SignupResponse, error) {
	// TODO: Wait for authing-gateway refactoring to complete the API
	return &user.SignupResponse{}, status.Error(codes.Internal, fmt.Errorf("NOT IMPLEMENTED").Error())
}

func (s *Server) CreateUser(ctx context.Context, in *user.CreateUserRequest) (*user.CreateUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateUser")
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
		logger.Sugar().Errorw("CreateUser", "err", err)
		return &user.CreateUserResponse{}, err
	}

	span = commontracer.TraceInvoker(span, "role", "middleware", "CreateUser")

	resp, err := usermwcli.CreateUser(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("CreateUser", "err", err)
		return &user.CreateUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &user.CreateUserResponse{
		Info: resp,
	}, nil
}
