package user

import (
	"context"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"
	signmethod "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"

	user1 "github.com/NpoolPlatform/appuser-gateway/pkg/user"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) UpdateUser(ctx context.Context, in *npool.UpdateUserRequest) (*npool.UpdateUserResponse, error) { //nolint
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateUser")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Infow("UpdateUser", "AppID", in.GetAppID())
		return &npool.UpdateUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Infow("UpdateUser", "UserID", in.GetUserID())
		return &npool.UpdateUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if in.NewAccount != nil && in.GetNewAccount() == "" {
		logger.Sugar().Infow("UpdateUser", "NewAccount", in.GetNewAccount())
		return &npool.UpdateUserResponse{}, status.Error(codes.InvalidArgument, "NewAccount is invalid")
	}
	if in.PasswordHash != nil && in.GetPasswordHash() == "" {
		logger.Sugar().Infow("UpdateUser", "PasswordHash", in.GetPasswordHash())
		return &npool.UpdateUserResponse{}, status.Error(codes.InvalidArgument, "PasswordHash is invalid")
	}
	if in.PasswordHash != nil {
		if in.GetOldPasswordHash() == "" && in.GetVerificationCode() == "" {
			logger.Sugar().Infow("UpdateUser", "PasswordHash", in.GetPasswordHash())
			return &npool.UpdateUserResponse{}, status.Error(codes.InvalidArgument, "permission denied")
		}
	}

	switch in.GetNewAccountType() {
	case signmethod.SignMethodType_Google:
		fallthrough //nolint
	case signmethod.SignMethodType_Email:
		fallthrough //nolint
	case signmethod.SignMethodType_Mobile:
		if in.GetNewVerificationCode() == "" || in.GetVerificationCode() == "" {
			logger.Sugar().Infow("UpdateUser", "NewVerificationCode", in.GetNewVerificationCode())
			return &npool.UpdateUserResponse{}, status.Error(codes.InvalidArgument, "NewVerificationCode is invalid")
		}
	}

	span = commontracer.TraceInvoker(span, "role", "middleware", "UpdateUser")

	info, err := user1.UpdateUser(ctx, in)
	if err != nil {
		logger.Sugar().Infow("UpdateUser", "error", err)
		return &npool.UpdateUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	_ = user1.UpdateCache(ctx, info)

	return &npool.UpdateUserResponse{
		Info: info,
	}, nil
}

func (s *Server) ResetUser(ctx context.Context, in *npool.ResetUserRequest) (*npool.ResetUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "ResetUser")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Infow("ResetUser", "AppID", in.GetAppID())
		return &npool.ResetUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Infow("ResetUser", "UserID", in.GetUserID())
		return &npool.ResetUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if in.GetAccount() == "" {
		logger.Sugar().Infow("ResetUser", "Account", in.GetAccount())
		return &npool.ResetUserResponse{}, status.Error(codes.InvalidArgument, "Account is invalid")
	}

	switch in.GetAccountType() {
	case signmethod.SignMethodType_Google:
	case signmethod.SignMethodType_Email:
	case signmethod.SignMethodType_Mobile:
	default:
		logger.Sugar().Infow("ResetUser", "AccountType", in.GetAccountType())
		return &npool.ResetUserResponse{}, status.Error(codes.InvalidArgument, "AccountType is invalid")
	}

	if in.GetVerificationCode() == "" {
		logger.Sugar().Infow("ResetUser", "VerificationCode", in.GetVerificationCode())
		return &npool.ResetUserResponse{}, status.Error(codes.InvalidArgument, "VerificationCode is invalid")
	}

	if in.GetPasswordHash() == "" {
		logger.Sugar().Infow("ResetUser", "PasswordHash", in.GetPasswordHash())
		return &npool.ResetUserResponse{}, status.Error(codes.InvalidArgument, "PasswordHash is invalid")
	}

	span = commontracer.TraceInvoker(span, "role", "middleware", "ResetUser")

	err = user1.ResetUser(ctx, in)
	if err != nil {
		logger.Sugar().Infow("ResetUser", "error", err)
		return &npool.ResetUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &npool.ResetUserResponse{}, nil
}
