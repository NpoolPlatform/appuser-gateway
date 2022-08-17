package user

import (
	"context"
	"fmt"

	rolemgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approle"
	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	inspirecli "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"
	rolemgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approle"
	"github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"
	usermwp "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"
	thirdgwconst "github.com/NpoolPlatform/third-gateway/pkg/const"
	"github.com/google/uuid"
)

func Signup(ctx context.Context, in *user.SignupRequest) (*usermwp.User, error) {
	app, err := appmwcli.GetApp(ctx, in.GetAppID())
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, fmt.Errorf("invalid app")
	}

	inviterID, err := checkInvitationCode(ctx, in.GetAppID(), in.GetInvitationCode(), app.InvitationCodeMust)
	if err != nil {
		logger.Sugar().Errorw("Signup", "err", err)
		return nil, err
	}

	err = VerifyCode(
		ctx,
		in.GetAppID(),
		"",
		in.GetAccount(),
		in.GetAccountType().String(),
		in.GetVerificationCode(),
		thirdgwconst.UsedForSignup,
		false,
	)
	if err != nil {
		return nil, err
	}

	emailAddress := ""
	phoneNO := ""

	if in.GetAccountType().String() == signmethod.SignMethodType_Mobile.String() {
		phoneNO = in.GetAccount()
	} else if in.GetAccountType().String() == signmethod.SignMethodType_Email.String() {
		emailAddress = in.GetAccount()
	}

	role, err := rolemgrcli.GetAppRoleOnly(ctx, &rolemgrpb.Conds{
		AppID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
		Default: &npool.BoolVal{
			Op:    cruder.EQ,
			Value: true,
		},
	})
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, fmt.Errorf("invalid default role")
	}

	userID := uuid.NewString()

	userInfo, err := usermwcli.CreateUser(ctx, &usermwp.UserReq{
		ID:           &userID,
		AppID:        &in.AppID,
		EmailAddress: &emailAddress,
		PhoneNO:      &phoneNO,
		PasswordHash: &in.PasswordHash,
		RoleIDs:      []string{role.ID},
	})
	if err != nil {
		return nil, err
	}

	if in.GetInvitationCode() == "" || inviterID == "" {
		return userInfo, nil
	}

	// TODO: revert user info
	_, err = inspirecli.CreateInvitation(ctx, in.AppID, inviterID, userID)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

func checkInvitationCode(ctx context.Context, appID, code string, must bool) (string, error) {
	if must && code == "" {
		return "", fmt.Errorf("invitation code is must")
	}

	if code == "" {
		return "", nil
	}

	ivc, err := inspirecli.GetUserInvitationCodeOnly(ctx, &inspirepb.Conds{
		InvitationCode: &npool.StringVal{
			Op:    cruder.EQ,
			Value: code,
		},
	})
	if err != nil {
		return "", err
	}

	if code == nil && must {
		return "", fmt.Errorf("invalid code")
	}

	if code.AppID != appID {
		return "", fmt.Errorf("invalid invitation code")
	}

	return code.UserID, nil
}
