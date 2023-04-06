package user

import (
	"context"
	"fmt"
	usercodemwcli "github.com/NpoolPlatform/basal-middleware/pkg/client/usercode"
	usercodemwpb "github.com/NpoolPlatform/message/npool/basal/mw/v1/usercode"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	ivcodemwpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/invitationcode"
	registrationmwpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/registration"

	rolemgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approle"
	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	ivcodemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/invitationcode"
	commonpb "github.com/NpoolPlatform/message/npool"
	appctrlmgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/appcontrol"
	rolemgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approle"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	"github.com/google/uuid"

	"github.com/NpoolPlatform/go-service-framework/pkg/pubsub"
	pubsubmsg "github.com/NpoolPlatform/message/npool/pubsub/v1"
)

//nolint
func Signup(
	ctx context.Context,
	appID, account, passwordHash string,
	accountType basetypes.SignMethod,
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

	if err := usercodemwcli.VerifyUserCode(ctx, &usercodemwpb.VerifyUserCodeRequest{
		Prefix:      basetypes.Prefix_PrefixUserCode.String(),
		AppID:       appID,
		Account:     account,
		AccountType: accountType,
		UsedFor:     basetypes.UsedFor_Signup,
		Code:        verificationCode,
	}); err != nil {
		return nil, err
	}

	emailAddress := ""
	phoneNO := ""

	if accountType.String() == basetypes.SignMethod_Mobile.String() {
		phoneNO = account
	} else if accountType.String() == basetypes.SignMethod_Email.String() {
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

	var code *ivcodemwpb.InvitationCode
	if app.CreateInvitationCodeWhen == appctrlmgrpb.CreateInvitationCodeWhen_Registration {
		_code, err := ivcodemwcli.CreateInvitationCode(ctx, &ivcodemwpb.InvitationCodeReq{
			AppID:  &appID,
			UserID: &userID,
		})

		if err != nil {
			return nil, err
		}
		code = _code
	}

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

	if code != nil {
		userInfo.InvitationCode = &code.InvitationCode
	}

	if invitationCode == nil || *invitationCode == "" || inviterID == "" {
		return userInfo, nil
	}

	err = pubsub.Publish(pubsubmsg.MessageID_SignupInvitation.String(), registrationmwpb.RegistrationReq{
		AppID:     &appID,
		InviterID: &inviterID,
		InviteeID: &userID,
	})
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

	ivc, err := ivcodemwcli.GetInvitationCodeOnly(ctx, &ivcodemwpb.Conds{
		InvitationCode: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: *code,
		},
	})
	if err != nil {
		return "", err
	}

	if ivc == nil {
		return "", fmt.Errorf("invalid code")
	}

	if ivc.AppID != appID {
		return "", fmt.Errorf("invalid invitation code")
	}

	return ivc.UserID, nil
}

func (h *Handler) Signup(ctx context.Context) (info *usermwpb.User, err error) {
	inviterID, err := h.CheckInvitationCode(ctx)
	if err != nil {
		return nil, err
	}

	if err := h.VerifyUserCode(ctx, basetypes.UsedFor_Signup); err != nil {
		return nil, err
	}

	signupHandler := &signupHandler{
		Handler:   h,
		inviterID: inviterID,
	}

	defer func() {
		if err == nil {
			return
		}

		if err := signupHandler.cancel(ctx); err != nil {
			logger.Sugar().Errorw(
				"Signup",
				"Step", "Cancel",
				"Error", err,
			)
		}
	}()

	err = signupHandler.try(ctx)
	if err != nil {
		return nil, err
	}

	info, err = signupHandler.confirm(ctx)
	if err != nil {
		return nil, err
	}

	return info, nil
}

type signupHandler struct {
	*Handler
	inviterID *string
}

/// Signup
///  1 Create invitation code according to application configuration
///  2 Create user
///  3 Create registration invitation

func (h *signupHandler) try(ctx context.Context) error {
	return nil
}

func (h *signupHandler) confirm(ctx context.Context) (*usermwpb.User, error) {
	return nil, nil
}

func (h *signupHandler) cancel(ctx context.Context) error {
	return nil
}
