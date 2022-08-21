package user

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	rolemgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approle"
	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	inspirecli "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/client"

	commonpb "github.com/NpoolPlatform/message/npool"
	rolemgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approle"
	signmethod "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"

	thirdgwconst "github.com/NpoolPlatform/third-gateway/pkg/const"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	"github.com/google/uuid"
)

func Signup(
	ctx context.Context,
	appID, account, passwordHash string,
	accountType signmethod.SignMethodType,
	verificationCode string,
	invitationCode *string,
) (
	*usermwpb.User, error,
) {
	app, err := appmwcli.GetApp(ctx, appID)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, fmt.Errorf("invalid app")
	}

	inviterID, err := checkInvitationCode(ctx, appID, invitationCode, app.InvitationCodeMust)
	if err != nil {
		logger.Sugar().Errorw("Signup", "err", err)
		return nil, err
	}

	err = verifyCode(
		ctx, appID, "",
		account, accountType,
		verificationCode,
		thirdgwconst.UsedForSignup,
		false,
	)
	if err != nil {
		return nil, err
	}

	emailAddress := ""
	phoneNO := ""

	if accountType.String() == signmethod.SignMethodType_Mobile.String() {
		phoneNO = account
	} else if accountType.String() == signmethod.SignMethodType_Email.String() {
		emailAddress = account
	}

	role, err := rolemgrcli.GetAppRoleOnly(ctx, &rolemgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
		Default: &commonpb.BoolVal{
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

	userInfo, err := usermwcli.CreateUser(ctx, &usermwpb.UserReq{
		ID:           &userID,
		AppID:        &appID,
		EmailAddress: &emailAddress,
		PhoneNO:      &phoneNO,
		PasswordHash: &passwordHash,
		RoleIDs:      []string{role.ID},
	})
	if err != nil {
		return nil, err
	}

	if invitationCode == nil || *invitationCode == "" || inviterID == "" {
		return userInfo, nil
	}

	// TODO: revert user info
	_, err = inspirecli.CreateInvitation(ctx, appID, inviterID, userID)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

func checkInvitationCode(ctx context.Context, appID string, code *string, must bool) (string, error) {
	if must && (code == nil || *code == "") {
		return "", fmt.Errorf("invitation code is must")
	}

	if code == nil || *code == "" {
		return "", nil
	}

	ivc, err := inspirecli.GetUserInvitationCodeOnly(ctx, &inspirepb.Conds{
		InvitationCode: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: *code,
		},
	})
	if err != nil {
		return "", err
	}

	if ivc == nil && must {
		return "", fmt.Errorf("invalid code")
	}

	if ivc.AppID != appID {
		return "", fmt.Errorf("invalid invitation code")
	}

	return ivc.UserID, nil
}
