package user

import (
	"context"

	"github.com/google/uuid"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	user1 "github.com/NpoolPlatform/appuser-gateway/pkg/user"

	signmethod "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
)

func (s *Server) Login(ctx context.Context, in *user.LoginRequest) (*user.LoginResponse, error) { //nolint
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "Login")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "AppID", in.GetAppID(), "error", err)
		return &user.LoginResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	if in.GetPasswordHash() == "" {
		logger.Sugar().Errorw("validate", "PasswordHash", in.GetPasswordHash())
		return &user.LoginResponse{}, status.Error(codes.InvalidArgument, "PasswordHash is invalid")
	}
	if in.GetAccount() == "" {
		logger.Sugar().Errorw("validate", "Account", in.GetAccount())
		return &user.LoginResponse{}, status.Error(codes.InvalidArgument, "Account is invalid")
	}

	switch in.GetAccountType() {
	case signmethod.SignMethodType_Email:
	case signmethod.SignMethodType_Mobile:
	case signmethod.SignMethodType_Twitter:
	case signmethod.SignMethodType_Github:
	case signmethod.SignMethodType_Facebook:
	case signmethod.SignMethodType_Linkedin:
	case signmethod.SignMethodType_Wechat:
	case signmethod.SignMethodType_Google:
	case signmethod.SignMethodType_Username:
	default:
		logger.Sugar().Errorw("validate", "AccountType", in.GetAccountType())
		return &user.LoginResponse{}, status.Error(codes.InvalidArgument, "AccountType is invalid")
	}

	if in.GetManMachineSpec() == "" {
		logger.Sugar().Errorw("validate", "ManMachineSpec", in.GetManMachineSpec())
		return &user.LoginResponse{}, status.Error(codes.InvalidArgument, "ManMachineSpec is invalid")
	}

	if in.GetEnvironmentSpec() == "" {
		logger.Sugar().Errorw("validate", "EnvironmentSpec", in.GetEnvironmentSpec())
		return &user.LoginResponse{}, status.Error(codes.InvalidArgument, "EnvironmentSpec is invalid")
	}

	span = commontracer.TraceInvoker(span, "user", "middleware", "Login")

	info, err := user1.Login(
		ctx,
		in.GetAppID(),
		in.GetAccount(),
		in.GetPasswordHash(),
		in.GetAccountType(),
		in.GetManMachineSpec(),
		in.GetEnvironmentSpec(),
	)
	if err != nil {
		logger.Sugar().Errorw("Login", "err", err)
		return &user.LoginResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &user.LoginResponse{
		Info: info,
	}, nil
}

func (s *Server) Logined(ctx context.Context, in *user.LoginedRequest) (*user.LoginedResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "Logined")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "AppID", in.GetAppID(), "error", err)
		return &user.LoginedResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Errorw("validate", "UserID", in.GetUserID(), "error", err)
		return &user.LoginedResponse{}, status.Error(codes.InvalidArgument, "UserID is invalid")
	}

	if in.GetToken() == "" {
		logger.Sugar().Errorw("validate", "Token", in.GetToken())
		return &user.LoginedResponse{}, status.Error(codes.InvalidArgument, "Token is invalid")
	}

	span = commontracer.TraceInvoker(span, "user", "middleware", "Logined")

	info, err := user1.Logined(
		ctx,
		in.GetAppID(),
		in.GetUserID(),
		in.GetToken(),
	)
	if err != nil {
		logger.Sugar().Errorw("Logined", "err", err)
		return &user.LoginedResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &user.LoginedResponse{
		Info: info,
	}, nil
}

func (s *Server) Logout(ctx context.Context, in *user.LogoutRequest) (*user.LogoutResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "Logout")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "AppID", in.GetAppID(), "error", err)
		return &user.LogoutResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Errorw("validate", "UserID", in.GetUserID(), "error", err)
		return &user.LogoutResponse{}, status.Error(codes.InvalidArgument, "UserID is invalid")
	}

	if in.GetToken() == "" {
		logger.Sugar().Errorw("validate", "Token", in.GetToken())
		return &user.LogoutResponse{}, status.Error(codes.InvalidArgument, "Token is invalid")
	}

	span = commontracer.TraceInvoker(span, "user", "middleware", "Logout")

	info, err := user1.Logout(
		ctx,
		in.GetAppID(),
		in.GetUserID(),
	)
	if err != nil {
		logger.Sugar().Errorw("Logout", "err", err)
		return &user.LogoutResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &user.LogoutResponse{
		Info: info,
	}, nil
}
