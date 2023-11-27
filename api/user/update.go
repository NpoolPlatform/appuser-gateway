package user

import (
	"context"

	user1 "github.com/NpoolPlatform/appuser-gateway/pkg/user"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"
	appusertypes "github.com/NpoolPlatform/message/npool/basetypes/appuser/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateUser(ctx context.Context, in *npool.UpdateUserRequest) (*npool.UpdateUserResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithID(&in.ID, true),
		user1.WithEntID(&in.UserID, true),
		user1.WithAppID(&in.AppID, true),
		user1.WithUserID(&in.UserID, true),
		user1.WithAccount(in.Account, false),
		user1.WithAccountType(in.AccountType, false),
		user1.WithNewAccount(in.NewAccount, false),
		user1.WithNewAccountType(in.NewAccountType, false),
		user1.WithPasswordHash(in.PasswordHash, false),
		user1.WithOldPasswordHash(in.OldPasswordHash, false),
		user1.WithVerificationCode(in.VerificationCode, false),
		user1.WithNewVerificationCode(in.NewVerificationCode, false),
		user1.WithUsername(in.Username, false),
		user1.WithAddressFields(in.AddressFields, false),
		user1.WithGender(in.Gender, false),
		user1.WithPostalCode(in.PostalCode, false),
		user1.WithAge(in.Age, false),
		user1.WithBirthday(in.Birthday, false),
		user1.WithAvatar(in.Avatar, false),
		user1.WithOrganization(in.Organization, false),
		user1.WithFirstName(in.FirstName, false),
		user1.WithLastName(in.LastName, false),
		user1.WithIDNumber(in.IDNumber, false),
		user1.WithSigninVerifyType(in.SigninVerifyType, false),
		user1.WithKolConfirmed(in.KolConfirmed, false),
		user1.WithSelectedLangID(in.SelectedLangID, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateUser",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateUser(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateUser",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateUserResponse{
		Info: info,
	}, nil
}

func (s *Server) UpdateAppUser(ctx context.Context, in *npool.UpdateAppUserRequest) (*npool.UpdateAppUserResponse, error) {
	updateCacheMode := appusertypes.UpdateCacheMode_DontUpdateCache
	handler, err := user1.NewHandler(
		ctx,
		user1.WithID(&in.ID, true),
		user1.WithEntID(&in.TargetUserID, true),
		user1.WithAppID(&in.AppID, true),
		user1.WithUserID(&in.TargetUserID, true),
		user1.WithKol(in.Kol, false),
		user1.WithEmailAddress(in.EmailAddress, false),
		user1.WithUpdateCacheMode(&updateCacheMode, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAppUser",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateAppUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateUser(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAppUser",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateAppUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateAppUserResponse{
		Info: info,
	}, nil
}
