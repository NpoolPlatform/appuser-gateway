package user

import (
	"context"

	appuserextracli "github.com/NpoolPlatform/appuser-manager/pkg/client/appuserextra"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"
	appuserextrapb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/appuserextra"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	regmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/registration"
	regmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/registration"

	user1 "github.com/NpoolPlatform/appuser-gateway/pkg/user"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

//nolint:dupl
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

	if in.IDNumber != nil {
		if in.GetIDNumber() == "" {
			logger.Sugar().Infow("UpdateUser", "IDNumber", in.GetIDNumber())
			return &npool.UpdateUserResponse{}, status.Error(codes.InvalidArgument, "IDNumber is invalid")
		}
		exist, err := appuserextracli.ExistAppUserExtraConds(ctx, &appuserextrapb.Conds{
			IDNumber: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: in.GetIDNumber(),
			},
		})
		if err != nil {
			logger.Sugar().Infow("CreateUser", "exist", exist, "err", err)
			return &npool.UpdateUserResponse{}, status.Error(codes.Internal, err.Error())
		}
		if exist {
			logger.Sugar().Infow("CreateUser", "IDNumber", in.GetIDNumber())
			return &npool.UpdateUserResponse{}, status.Error(codes.InvalidArgument, "IDNumber is already exists")
		}
	}

	switch in.GetNewAccountType() {
	case basetypes.SignMethod_Google:
		fallthrough //nolint
	case basetypes.SignMethod_Email:
		fallthrough //nolint
	case basetypes.SignMethod_Mobile:
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
	if in.UserID != nil {
		if _, err := uuid.Parse(in.GetUserID()); err != nil {
			logger.Sugar().Infow("ResetUser", "UserID", in.GetUserID())
			return &npool.ResetUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	if in.GetAccount() == "" {
		logger.Sugar().Infow("ResetUser", "Account", in.GetAccount())
		return &npool.ResetUserResponse{}, status.Error(codes.InvalidArgument, "Account is invalid")
	}

	switch in.GetAccountType() {
	case basetypes.SignMethod_Email:
	case basetypes.SignMethod_Mobile:
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

func (s *Server) UpdateUserKol(ctx context.Context, in *npool.UpdateUserKolRequest) (*npool.UpdateUserKolResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateUserKol")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Infow("UpdateUserKol", "AppID", in.GetAppID())
		return &npool.UpdateUserKolResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Infow("UpdateUserKol", "UserID", in.GetUserID())
		return &npool.UpdateUserKolResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetTargetUserID()); err != nil {
		logger.Sugar().Infow("UpdateUserKol", "TargetUserID", in.GetTargetUserID())
		return &npool.UpdateUserKolResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	// TODO: check if user is KOL, and target user is invited by user
	reg, err := regmwcli.GetRegistrationOnly(ctx, &regmgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
		InviterID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetUserID(),
		},
		InviteeID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetTargetUserID(),
		},
	})
	if err != nil {
		return &npool.UpdateUserKolResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if reg == nil {
		return &npool.UpdateUserKolResponse{}, status.Error(codes.InvalidArgument, "permission denied")
	}

	span = commontracer.TraceInvoker(span, "role", "middleware", "UpdateUser")

	info, err := user1.UpdateUserKol(ctx, in)
	if err != nil {
		logger.Sugar().Infow("UpdateUser", "error", err)
		return &npool.UpdateUserKolResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	_ = user1.UpdateCache(ctx, info)

	return &npool.UpdateUserKolResponse{
		Info: info,
	}, nil
}

func (s *Server) UpdateAppUserKol(ctx context.Context, in *npool.UpdateAppUserKolRequest) (*npool.UpdateAppUserKolResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateAppUserKol")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Infow("UpdateAppUserKol", "AppID", in.GetAppID())
		return &npool.UpdateAppUserKolResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetTargetUserID()); err != nil {
		logger.Sugar().Infow("UpdateAppUserKol", "TargetUserID", in.GetTargetUserID())
		return &npool.UpdateAppUserKolResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	span = commontracer.TraceInvoker(span, "role", "middleware", "UpdateAppUser")

	info, err := user1.UpdateUserKol(ctx, &npool.UpdateUserKolRequest{
		AppID:        in.GetAppID(),
		TargetUserID: in.GetTargetUserID(),
		Kol:          in.GetKol(),
	})
	if err != nil {
		logger.Sugar().Infow("UpdateAppUser", "error", err)
		return &npool.UpdateAppUserKolResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &npool.UpdateAppUserKolResponse{
		Info: info,
	}, nil
}
