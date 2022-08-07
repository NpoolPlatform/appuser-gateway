package user

import (
	"context"

	usermgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/appuser"
	mw "github.com/NpoolPlatform/appuser-middleware/api/v1/user"
	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"
	usermgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/appuser"
	"github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"
	mgruser "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validate(info *mgruser.UserReq) error {
	err := mw.Validate(info)
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}

//nolint:gocyclo,nolintlint
func signUpValidate(ctx context.Context, info *user.SignupRequest) error {
	if _, err := uuid.Parse(info.GetAppID()); err != nil {
		logger.Sugar().Errorw("validate", "AppID", info.GetAppID(), "error", err)
		return status.Error(codes.InvalidArgument, "AppID is invalid")
	}
	if info.GetUsername() == "" {
		logger.Sugar().Errorw("validate", "Username", info.GetUsername())
		return status.Error(codes.InvalidArgument, "Username is invalid")
	}
	if info.GetPasswordHash() == "" {
		logger.Sugar().Errorw("validate", "PasswordHash", info.GetPasswordHash())
		return status.Error(codes.InvalidArgument, "PasswordHash is invalid")
	}
	if info.GetAccount() == "" {
		logger.Sugar().Errorw("validate", "Account", info.GetAccount())
		return status.Error(codes.InvalidArgument, "Account is invalid")
	}
	if info.GetAccountType() == "" {
		logger.Sugar().Errorw("validate", "AccountType", info.GetAccountType())
		return status.Error(codes.InvalidArgument, "AccountType is invalid")
	}

	methodExist := false
	for _, val := range signmethod.SignMethodType_name {
		if info.GetAccountType() == val {
			methodExist = true
		}
	}
	if !methodExist {
		logger.Sugar().Errorw("validate", "AccountType", info.GetAccountType())
		return status.Error(codes.InvalidArgument, "signup method not exist")
	}

	if info.GetVerificationCode() == "" {
		logger.Sugar().Errorw("validate", "VerificationCode", info.GetVerificationCode())
		return status.Error(codes.InvalidArgument, "VerificationCode is invalid")
	}

	existP, err := usermgrcli.ExistAppUserConds(ctx, &usermgrpb.Conds{
		AppID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: info.GetAppID(),
		},
		PhoneNo: &npool.StringVal{
			Op:    cruder.EQ,
			Value: info.GetAccount(),
		},
	})
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return status.Error(codes.Internal, err.Error())
	}

	existA, err := usermgrcli.ExistAppUserConds(ctx, &usermgrpb.Conds{
		AppID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: info.GetAppID(),
		},
		EmailAddress: &npool.StringVal{
			Op:    cruder.EQ,
			Value: info.GetAccount(),
		},
	})
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return status.Error(codes.Internal, err.Error())
	}

	if existP || existA {
		return status.Error(codes.AlreadyExists, "account already exist")
	}

	app, err := appmwcli.GetApp(ctx, info.GetAppID())
	if err != nil {
		logger.Sugar().Errorw("validate", "err", err)
		return status.Error(codes.Internal, err.Error())
	}
	if app == nil {
		logger.Sugar().Errorw("validate", "fail get app")
		return status.Error(codes.Internal, "fail get app")
	}

	if app.InvitationCodeMust {
		if info.GetInvitationCode() == "" {
			logger.Sugar().Errorw("validate", "invitation code is must")
			return status.Error(codes.InvalidArgument, "invitation code is must")
		}
	}

	return nil
}
