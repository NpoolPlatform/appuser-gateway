package user

import (
	"context"

	"github.com/google/uuid"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	user1 "github.com/NpoolPlatform/appuser-gateway/pkg/user"

	usermgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/appuser"
	usermgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/appuser"
	signmethod "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"

	commonpb "github.com/NpoolPlatform/message/npool"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
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

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "AppID", in.GetAppID(), "error", err)
		return &user.SignupResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	if in.GetPasswordHash() == "" {
		logger.Sugar().Errorw("validate", "PasswordHash", in.GetPasswordHash())
		return &user.SignupResponse{}, status.Error(codes.InvalidArgument, "PasswordHash is invalid")
	}
	if in.GetAccount() == "" {
		logger.Sugar().Errorw("validate", "Account", in.GetAccount())
		return &user.SignupResponse{}, status.Error(codes.InvalidArgument, "Account is invalid")
	}
	if in.GetAccountType().String() == "" {
		logger.Sugar().Errorw("validate", "AccountType", in.GetAccountType())
		return &user.SignupResponse{}, status.Error(codes.InvalidArgument, "AccountType is invalid")
	}

	methodExist := false
	for _, val := range signmethod.SignMethodType_name {
		if in.GetAccountType().String() == val {
			methodExist = true
			break
		}
	}
	if !methodExist {
		logger.Sugar().Errorw("validate", "AccountType", in.GetAccountType())
		return &user.SignupResponse{}, status.Error(codes.InvalidArgument, "AccountType is invalid")
	}

	if in.GetVerificationCode() == "" {
		logger.Sugar().Errorw("validate", "VerificationCode", in.GetVerificationCode())
		return &user.SignupResponse{}, status.Error(codes.InvalidArgument, "VerificationCode is invalid")
	}

	exist, err := usermgrcli.ExistAppUserConds(ctx, &usermgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
		PhoneNO: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAccount(),
		},
	})
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &user.SignupResponse{}, status.Error(codes.Internal, err.Error())
	}
	if exist {
		return &user.SignupResponse{}, status.Error(codes.AlreadyExists, "account already exist")
	}

	exist, err = usermgrcli.ExistAppUserConds(ctx, &usermgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
		EmailAddress: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAccount(),
		},
	})
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return &user.SignupResponse{}, status.Error(codes.Internal, err.Error())
	}

	if exist {
		return &user.SignupResponse{}, status.Error(codes.AlreadyExists, "account already exist")
	}

	span = commontracer.TraceInvoker(span, "user", "middleware", "Signup")

	userInfo, err := user1.Signup(ctx, in)
	if err != nil {
		logger.Sugar().Errorw("Signup", "err", err)
		return &user.SignupResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &user.SignupResponse{
		Info: userInfo,
	}, nil
}
