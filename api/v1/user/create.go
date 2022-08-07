//nolint:nolintlint,dupl
package user

import (
	"context"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	tracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer/user"
	mgrtracer "github.com/NpoolPlatform/appuser-middleware/pkg/tracer/user"
	"google.golang.org/grpc/codes"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	mw "github.com/NpoolPlatform/appuser-gateway/pkg/user"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Signup(ctx context.Context, in *user.SignupRequest) (*user.SignupResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "Signup")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = tracer.Trace(span, in)

	err = signUpValidate(ctx, in)
	if err != nil {
		logger.Sugar().Errorw("Signup", "err", err)
		return &user.SignupResponse{}, err
	}

	span = commontracer.TraceInvoker(span, "user", "middleware", "Signup")

	userInfo, err := mw.Signup(ctx, in)
	if err != nil {
		return nil, err
	}
	return &user.SignupResponse{
		Info: userInfo,
	}, nil
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

	span = mgrtracer.Trace(span, in.GetInfo())

	err = validate(in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("CreateUser", "err", err)
		return &user.CreateUserResponse{}, err
	}

	span = commontracer.TraceInvoker(span, "user", "middleware", "CreateUser")

	resp, err := usermwcli.CreateUser(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("CreateUser", "err", err)
		return &user.CreateUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &user.CreateUserResponse{
		Info: resp,
	}, nil
}
