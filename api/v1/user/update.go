package user

import (
	"context"

	invitationcli "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/client"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"

	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"
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

	if in.GetEmailAddress() != "" && in.GetPhoneNO() != "" {
		logger.Sugar().Infow("UpdateUser", "EmailAddress", in.GetEmailAddress(), "PhoneNO", in.GetPhoneNO())
		return &npool.UpdateUserResponse{}, status.Error(codes.InvalidArgument, "Can't update email and phone numbers together")
	}

	if in.GetEmailAddress() != "" || in.GetPhoneNO() != "" || in.GetPasswordHash() != "" {
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
	}

	account := in.GetEmailAddress()
	if account == "" {
		account = in.GetPhoneNO()
	}

	if account != "" {
		if err := user1.VerifyCode(
			ctx,
			in.GetAppID(),
			in.GetUserID(),
			account,
			in.GetAccountType(),
			in.GetNewVerificationCode(),
			thirdgwconst.UsedForUpdate,
		); err != nil {
			logger.Sugar().Infow("UpdateUser", "VerificationCode", in.GetVerificationCode())
			return &npool.UpdateUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	span = commontracer.TraceInvoker(span, "role", "middleware", "UpdateUser")

	info, err := usermwcli.UpdateUser(ctx, &usermwpb.UserReq{
		ID:               &in.UserID,
		AppID:            &in.AppID,
		EmailAddress:     in.EmailAddress,
		PhoneNO:          in.PhoneNO,
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
	})
	if err != nil {
		logger.Sugar().Errorw("UpdateUser", "err", err)
		return &npool.UpdateUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	if in.InvitationCodeID != nil {
		if in.InvitationCodeConfirmed == nil {
			logger.Sugar().Errorw("UpdateUser", "err", err)
			return &npool.UpdateUserResponse{}, status.Error(codes.InvalidArgument, "InvitationCodeConfirmed empty")
		}

		if _, err = uuid.Parse(in.GetInvitationCodeID()); err != nil {
			logger.Sugar().Errorw("UpdateUser", "err", err)
			return &npool.UpdateUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}

		invite, err := invitationcli.UpdateUserInvitationCode(ctx, in.GetInvitationCodeID(), in.GetInvitationCodeConfirmed())
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
