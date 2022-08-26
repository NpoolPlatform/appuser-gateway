package user

import (
	"context"

	invitationcli "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/client"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"

	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"
	signmethod "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	user1 "github.com/NpoolPlatform/appuser-gateway/pkg/user"
	thirdgwconst "github.com/NpoolPlatform/third-gateway/pkg/const"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

//nolint:gocyclo
func (s *Server) UpdateUser(ctx context.Context, in *npool.UpdateUserRequest) (*npool.UpdateUserResponse, error) {
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

	switch in.GetNewAccountType() {
	case signmethod.SignMethodType_Google:
		fallthrough //nolint
	case signmethod.SignMethodType_Email:
		fallthrough //nolint
	case signmethod.SignMethodType_Mobile:
		if in.GetNewVerificationCode() == "" {
			logger.Sugar().Infow("UpdateUser", "NewVerificationCode", in.GetNewVerificationCode())
			return &npool.UpdateUserResponse{}, status.Error(codes.InvalidArgument, "NewVerificationCode is invalid")
		}
	}

	if in.NewAccount != nil || in.PasswordHash != nil || in.GetNewAccountType() == signmethod.SignMethodType_Google {
		if err := user1.VerifyCode(
			ctx,
			in.GetAppID(),
			in.GetUserID(),
			in.GetAccount(),
			in.GetAccountType(),
			in.GetVerificationCode(),
			thirdgwconst.UsedForUpdate,
		); err != nil {
			logger.Sugar().Infow("UpdateUser", "VerificationCode", in.GetVerificationCode())
			return &npool.UpdateUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}

		if err := user1.VerifyCode(
			ctx,
			in.GetAppID(),
			in.GetUserID(),
			in.GetNewAccount(),
			in.GetNewAccountType(),
			in.GetNewVerificationCode(),
			thirdgwconst.UsedForUpdate,
		); err != nil {
			logger.Sugar().Infow("UpdateUser", "NewVerificationCode", in.GetNewVerificationCode())
			return &npool.UpdateUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	span = commontracer.TraceInvoker(span, "role", "middleware", "UpdateUser")

	req := &usermwpb.UserReq{
		ID:               &in.UserID,
		AppID:            &in.AppID,
		Username:         in.Username,
		AddressFields:    in.AddressFields,
		Gender:           in.Gender,
		PostalCode:       in.PostalCode,
		Age:              in.Age,
		Birthday:         in.Birthday,
		Avatar:           in.Avatar,
		Organization:     in.Organization,
		FirstName:        in.FirstName,
		LastName:         in.LastName,
		IDNumber:         in.IDNumber,
		SigninVerifyType: in.SigninVerifyType,
		PasswordHash:     in.PasswordHash,
	}
	switch in.GetNewAccountType() {
	case signmethod.SignMethodType_Google:
		verified := true
		req.GoogleAuthVerified = &verified
	case signmethod.SignMethodType_Email:
		req.EmailAddress = in.NewAccount
	case signmethod.SignMethodType_Mobile:
		req.PhoneNO = in.NewAccount
	}

	info, err := usermwcli.UpdateUser(ctx, req)
	if err != nil {
		logger.Sugar().Errorw("UpdateUser", "err", err)
		return &npool.UpdateUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	if in.InvitationCodeID != nil && in.InvitationCodeConfirmed != nil {
		if _, err = uuid.Parse(in.GetInvitationCodeID()); err != nil {
			logger.Sugar().Errorw("UpdateUser", "err", err)
			return &npool.UpdateUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}

		invite, err := invitationcli.UpdateUserInvitationCode(
			ctx,
			in.GetInvitationCodeID(),
			in.GetInvitationCodeConfirmed(),
		)
		if err != nil {
			logger.Sugar().Errorw("UpdateUser", "err", err)
			return &npool.UpdateUserResponse{}, status.Error(codes.Internal, err.Error())
		}

		info.InvitationCodeConfirmed = invite.Confirmed
	}

	_ = user1.UpdateCache(ctx, info)

	return &npool.UpdateUserResponse{
		Info: info,
	}, nil
}
