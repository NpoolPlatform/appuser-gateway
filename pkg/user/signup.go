package user

import (
	"context"
	"fmt"

	usercodemwcli "github.com/NpoolPlatform/basal-middleware/pkg/client/usercode"
	usercodemwpb "github.com/NpoolPlatform/message/npool/basal/mw/v1/usercode"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	ivcodemwpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/invitationcode"

	rolemgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approle"
	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	ivcodemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/invitationcode"
	commonpb "github.com/NpoolPlatform/message/npool"
	appctrlmgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/appcontrol"
	rolemgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approle"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	appusersvcname "github.com/NpoolPlatform/appuser-middleware/pkg/servicename"
	inspiremwsvcname "github.com/NpoolPlatform/inspire-middleware/pkg/servicename"

	ivcodemgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/invitationcode"
	registrationmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/registration"

	dtmcli "github.com/NpoolPlatform/dtm-cluster/pkg/dtm"
	"github.com/dtm-labs/dtmcli/dtmimp"

	"github.com/NpoolPlatform/go-service-framework/pkg/pubsub"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	"github.com/google/uuid"
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

	/*
		err = pubsub.Publish(
			basetypes.MsgID_CreateRegistrationInvitationConfirm.String(),
			nil,
			nil,
			nil,
			registrationmwpb.RegistrationReq{
				AppID:     &appID,
				InviterID: &inviterID,
				InviteeID: &userID,
			},
		)
		if err != nil {
			return nil, err
		}
	*/

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

func (h *signupHandler) withCreateInvitationCode(ctx context.Context, dispose *dtmcli.SagaDispose) {
	if h.App.CreateInvitationCodeWhen != appctrlmgrpb.CreateInvitationCodeWhen_Registration {
		return
	}

	id := uuid.NewString()
	req := &ivcodemgrpb.InvitationCodeReq{
		ID:     &id,
		AppID:  &h.AppID,
		UserID: &h.userID,
	}

	dispose.Add(
		inspiremwsvcname.ServiceDomain,
		"inspire.middleware.invitation.invitationcode.v1.Middleware/CreateInvitationCode",
		"inspire.middleware.invitation.invitationcode.v1.Middleware/DeleteInvitationCode",
		req,
	)
}

func (h *signupHandler) withCreateUser(ctx context.Context, dispose *dtmcli.SagaDispose) {
	req := &usermwpb.UserReq{
		ID:           &h.userID,
		AppID:        &h.AppID,
		EmailAddress: h.EmailAddress,
		PhoneNO:      h.PhoneNO,
		PasswordHash: &h.PasswordHash,
		RoleIDs:      []string{h.defaultRole.ID},
	}

	dispose.Add(
		appusersvcname.ServiceDomain,
		"appuser.middleware.user.v1.Middleware/CreateUser",
		"appuser.middleware.user.v1.Middleware/DeleteUser",
		req,
	)
}

func (h *signupHandler) withCreateRegistrationInvitation(ctx context.Context, dispose *dtmcli.SagaDispose) {
	if h.inviterID == nil {
		return
	}

	req := &registrationmgrpb.RegistrationReq{
		AppID:     &h.AppID,
		InviterID: h.inviterID,
		InviteeID: &h.userID,
	}

	dispose.Add(
		inspiremwsvcname.ServiceDomain,
		"inspire.middleware.invitation.registration.v1.Middleware/CreateRegistration",
		"inspire.middleware.invitation.registration.v1.Middleware/DeleteRegistration",
		req,
	)
}

/// Signup
///  1 Create invitation code according to application configuration
///  2 Create user
///  3 Create registration invitation
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
		userID:    uuid.NewString(),
	}

	if err := signupHandler.getDefaultRole(ctx); err != nil {
		return nil, err
	}

	sagaDispose := dtmcli.NewSagaDispose(dtmimp.TransOptions{})
	signupHandler.withCreateInvitationCode(ctx, sagaDispose)
	signupHandler.withCreateUser(ctx, sagaDispose)
	signupHandler.withCreateRegistrationInvitation(ctx, sagaDispose)

	if err := dtmcli.WithSaga(ctx, sagaDispose); err != nil {
		return nil, err
	}

	/// TODO: if newbie has coupon, send event to allocate coupon, and we don't care about allocate result

	return info, nil
}

type signupHandler struct {
	*Handler
	inviterID   *string
	userID      string
	publisher   *pubsub.Publisher
	defaultRole *rolemgrpb.AppRole
	tryResp     map[basetypes.MsgID]uuid.UUID
}

func (h *signupHandler) getDefaultRole(ctx context.Context) error {
	role, err := rolemgrcli.GetAppRoleOnly(ctx, &rolemgrpb.Conds{
		AppID:   &commonpb.StringVal{Op: cruder.EQ, Value: h.AppID},
		Default: &commonpb.BoolVal{Op: cruder.EQ, Value: true},
	})
	if err != nil {
		return err
	}
	if role == nil {
		return fmt.Errorf("invalid default role")
	}
	h.defaultRole = role
	return nil
}
