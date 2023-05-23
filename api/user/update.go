package user

import (
	"context"

	user1 "github.com/NpoolPlatform/appuser-gateway/pkg/user"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateUser(ctx context.Context, in *npool.UpdateUserRequest) (*npool.UpdateUserResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(in.GetAppID()),
		user1.WithUserID(&in.UserID),
		user1.WithAccount(in.Account, in.AccountType),
		user1.WithNewAccount(in.NewAccount, in.NewAccountType),
		user1.WithPasswordHash(in.PasswordHash),
		user1.WithOldPasswordHash(in.OldPasswordHash),
		user1.WithVerificationCode(in.VerificationCode),
		user1.WithNewVerificationCode(in.NewVerificationCode),
		user1.WithUsername(in.Username),
		user1.WithAddressFields(in.AddressFields),
		user1.WithGender(in.Gender),
		user1.WithPostalCode(in.PostalCode),
		user1.WithAge(in.Age),
		user1.WithBirthday(in.Birthday),
		user1.WithAvatar(in.Avatar),
		user1.WithOrganization(in.Organization),
		user1.WithFirstName(in.FirstName),
		user1.WithLastName(in.LastName),
		user1.WithIDNumber(in.IDNumber),
		user1.WithSigninVerifyType(in.SigninVerifyType),
		user1.WithKolConfirmed(in.KolConfirmed),
		user1.WithSelectedLangID(in.SelectedLangID),
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
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(in.GetAppID()),
		user1.WithUserID(&in.TargetUserID),
		user1.WithKol(in.Kol),
		user1.WithEmailAddress(in.EmailAddress),
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
