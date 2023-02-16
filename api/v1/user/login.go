package user

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"

	user1 "github.com/NpoolPlatform/appuser-gateway/pkg/user"

	ivcodemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/invitationcode"
	ivcodemwpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/invitationcode"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"

	historycli "github.com/NpoolPlatform/appuser-manager/pkg/client/login/history"
	historypb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/login/history"

	"github.com/NpoolPlatform/message/npool"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"

	"github.com/google/uuid"
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
	case basetypes.SignMethod_Email:
	case basetypes.SignMethod_Mobile:
	case basetypes.SignMethod_Twitter:
	case basetypes.SignMethod_Github:
	case basetypes.SignMethod_Facebook:
	case basetypes.SignMethod_Linkedin:
	case basetypes.SignMethod_Wechat:
	case basetypes.SignMethod_Google:
	case basetypes.SignMethod_Username:
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

func (s *Server) LoginVerify(ctx context.Context, in *user.LoginVerifyRequest) (*user.LoginVerifyResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "LoginVerify")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "AppID", in.GetAppID(), "error", err)
		return &user.LoginVerifyResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Errorw("validate", "UserID", in.GetUserID(), "error", err)
		return &user.LoginVerifyResponse{}, status.Error(codes.InvalidArgument, "UserID is invalid")
	}

	if in.GetToken() == "" {
		logger.Sugar().Errorw("validate", "Token", in.GetToken())
		return &user.LoginVerifyResponse{}, status.Error(codes.InvalidArgument, "Token is invalid")
	}

	if in.GetVerificationCode() == "" {
		logger.Sugar().Errorw("validate", "VerificationCode", in.GetVerificationCode())
		return &user.LoginVerifyResponse{}, status.Error(codes.InvalidArgument, "VerificationCode is invalid")
	}

	span = commontracer.TraceInvoker(span, "user", "middleware", "LoginVerify")

	info, err := user1.LoginVerify(
		ctx,
		in.GetAppID(),
		in.GetUserID(),
		in.GetToken(),
		in.GetAccount(),
		in.GetAccountType(),
		in.GetVerificationCode(),
	)
	if err != nil {
		logger.Sugar().Errorw("LoginVerify", "err", err)
		return &user.LoginVerifyResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &user.LoginVerifyResponse{
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
		logger.Sugar().Errorw("Logined", "error", err)
		return &user.LoginedResponse{}, status.Error(codes.Internal, err.Error())
	}
	if info == nil {
		return &user.LoginedResponse{}, nil
	}

	code, err := ivcodemwcli.GetInvitationCodeOnly(ctx, &ivcodemwpb.Conds{
		AppID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
		UserID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: in.GetUserID(),
		},
	})
	if err != nil {
		logger.Sugar().Errorw("UpdateCache", "error", err)
		return &user.LoginedResponse{}, status.Error(codes.Internal, err.Error())
	}
	if code != nil {
		info.InvitationCode = &code.InvitationCode
	}

	_ = user1.UpdateCache(ctx, info)

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

func (s *Server) GetLoginHistories(ctx context.Context, in *user.GetLoginHistoriesRequest) (*user.GetLoginHistoriesResponse, error) {
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
		return &user.GetLoginHistoriesResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Errorw("validate", "UserID", in.GetUserID(), "error", err)
		return &user.GetLoginHistoriesResponse{}, status.Error(codes.InvalidArgument, "UserID is invalid")
	}

	infos, total, err := historycli.GetHistories(ctx, &historypb.Conds{
		AppID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
		UserID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: in.GetUserID(),
		},
	}, in.GetOffset(), in.GetLimit())
	if err != nil {
		return nil, err
	}

	return &user.GetLoginHistoriesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
