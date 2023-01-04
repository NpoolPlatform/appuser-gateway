package user

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/message/npool/third/mgr/v1/usedfor"
	thirdmwcli "github.com/NpoolPlatform/third-middleware/pkg/client/verify"

	invitationcli "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/client"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	appusermgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/appuser"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"

	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/appuser"
	signmethod "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	commonpb "github.com/NpoolPlatform/message/npool"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func UpdateUser(ctx context.Context, in *npool.UpdateUserRequest) (*usermwpb.User, error) { //nolint
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateUser")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if in.GetOldPasswordHash() != "" {
		if _, err := usermwcli.VerifyUser(
			ctx,
			in.GetAppID(),
			in.GetUserID(),
			in.GetOldPasswordHash(),
		); err != nil {
			logger.Sugar().Infow("UpdateUser", "error", err)
			return nil, status.Error(codes.InvalidArgument, "permission denied")
		}
	}

	user, err := usermwcli.GetUser(ctx, in.GetAppID(), in.GetUserID())
	if err != nil {
		return nil, err
	}

	if in.NewAccount != nil || in.PasswordHash != nil || in.GetNewAccountType() == signmethod.SignMethodType_Google {
		account := in.GetAccount()
		if in.GetAccountType() == signmethod.SignMethodType_Google {
			account = user.GoogleSecret
		}

		if err := thirdmwcli.VerifyCode(
			ctx, in.GetAppID(),
			account,
			in.GetVerificationCode(),
			in.GetAccountType(),
			usedfor.UsedFor_Update,
		); err != nil {
			return nil, err
		}
	}

	if in.NewAccount != nil || in.GetNewAccountType() == signmethod.SignMethodType_Google {
		account := in.GetNewAccount()
		if in.GetNewAccountType() == signmethod.SignMethodType_Google {
			account = user.GoogleSecret
		}

		if err := thirdmwcli.VerifyCode(
			ctx, in.GetAppID(),
			account,
			in.GetNewVerificationCode(),
			in.GetNewAccountType(),
			usedfor.UsedFor_Update,
		); err != nil {
			return nil, err
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
		return nil, err
	}

	if in.InvitationCodeID != nil && in.InvitationCodeConfirmed != nil {
		if _, err = uuid.Parse(in.GetInvitationCodeID()); err != nil {
			logger.Sugar().Errorw("UpdateUser", "err", err)
			return nil, err
		}

		invite, err := invitationcli.UpdateUserInvitationCode(
			ctx,
			in.GetInvitationCodeID(),
			in.GetInvitationCodeConfirmed(),
		)
		if err != nil {
			logger.Sugar().Errorw("UpdateUser", "err", err)
			return nil, err
		}

		info.InvitationCodeConfirmed = invite.Confirmed
	}

	_ = UpdateCache(ctx, info)

	return info, nil
}

func ResetUser(ctx context.Context, in *npool.ResetUserRequest) error {
	conds := &appusermgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
	}

	switch in.GetAccountType() {
	case signmethod.SignMethodType_Email:
		conds.EmailAddress = &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAccount(),
		}
	case signmethod.SignMethodType_Mobile:
		conds.PhoneNO = &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAccount(),
		}
	default:
		return fmt.Errorf("invalid account type")
	}

	auser, err := appusermgrcli.GetAppUserOnly(ctx, conds)
	if err != nil {
		return fmt.Errorf("invalid user")
	}

	if err := thirdmwcli.VerifyCode(
		ctx, in.GetAppID(),
		in.GetAccount(),
		in.GetVerificationCode(),
		in.GetAccountType(),
		usedfor.UsedFor_Update,
	); err != nil {
		return err
	}

	_, err = usermwcli.UpdateUser(ctx, &usermwpb.UserReq{
		ID:           &auser.ID,
		AppID:        &in.AppID,
		PasswordHash: in.PasswordHash,
	})

	return err
}

func UpdateUserKol(ctx context.Context, in *npool.UpdateUserKolRequest) (*usermwpb.User, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateUserKol")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "role", "middleware", "UpdateUserKol")

	req := &usermwpb.UserReq{
		ID:    &in.TargetUserID,
		AppID: &in.AppID,
		Kol:   &in.Kol,
	}

	info, err := usermwcli.UpdateUser(ctx, req)
	if err != nil {
		logger.Sugar().Errorw("UpdateUserKol", "err", err)
		return nil, err
	}

	return info, nil
}
