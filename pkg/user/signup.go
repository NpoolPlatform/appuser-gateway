package user

import (
	"context"
	"fmt"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	rolemgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approle"
	usermgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/appuser"

	appusersvcname "github.com/NpoolPlatform/appuser-middleware/pkg/servicename"
	inspiremwsvcname "github.com/NpoolPlatform/inspire-middleware/pkg/servicename"

	appctrlmgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/appcontrol"
	rolemgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approle"
	usermgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/appuser"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	ivcodemgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/invitationcode"
	registrationmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/registration"
	ivcodemwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/invitation/invitationcode"
	registrationmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/invitation/registration"

	dtmcli "github.com/NpoolPlatform/dtm-cluster/pkg/dtm"
	"github.com/dtm-labs/dtmcli/dtmimp"

	_ "github.com/NpoolPlatform/go-service-framework/pkg/pubsub"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"

	"github.com/google/uuid"
)

type signupHandler struct {
	*Handler
	inviterID   *string
	defaultRole *rolemgrpb.AppRole
}

func (h *signupHandler) withCreateInvitationCode(dispose *dtmcli.SagaDispose) {
	if h.App.CreateInvitationCodeWhen != appctrlmgrpb.CreateInvitationCodeWhen_Registration {
		return
	}

	id := uuid.NewString()
	req := &ivcodemgrpb.InvitationCodeReq{
		ID:     &id,
		AppID:  &h.AppID,
		UserID: &h.UserID,
	}

	dispose.Add(
		inspiremwsvcname.ServiceDomain,
		"inspire.middleware.invitation.invitationcode.v1.Middleware/CreateInvitationCode",
		"inspire.middleware.invitation.invitationcode.v1.Middleware/DeleteInvitationCode",
		&ivcodemwpb.CreateInvitationCodeRequest{
			Info: req,
		},
	)
}

func (h *signupHandler) withCreateUser(dispose *dtmcli.SagaDispose) {
	req := &usermwpb.UserReq{
		ID:           &h.UserID,
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
		&usermwpb.CreateUserRequest{
			Info: req,
		},
	)
}

func (h *signupHandler) withCreateRegistrationInvitation(dispose *dtmcli.SagaDispose) {
	if h.inviterID == nil {
		return
	}

	req := &registrationmgrpb.RegistrationReq{
		AppID:     &h.AppID,
		InviterID: h.inviterID,
		InviteeID: &h.UserID,
	}

	dispose.Add(
		inspiremwsvcname.ServiceDomain,
		"inspire.middleware.invitation.registration.v1.Middleware/CreateRegistration",
		"inspire.middleware.invitation.registration.v1.Middleware/DeleteRegistration",
		&registrationmwpb.CreateRegistrationRequest{
			Info: req,
		},
	)
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

func (h *signupHandler) checkUser(ctx context.Context) error {
	conds := &usermgrpb.Conds{
		AppID: &commonpb.StringVal{Op: cruder.EQ, Value: h.AppID},
	}
	if h.EmailAddress != nil {
		conds.EmailAddress = &commonpb.StringVal{Op: cruder.EQ, Value: *h.EmailAddress}
	}
	if h.PhoneNO != nil {
		conds.PhoneNO = &commonpb.StringVal{Op: cruder.EQ, Value: *h.PhoneNO}
	}

	exist, err := usermgrcli.ExistAppUserConds(ctx, conds)
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("user already exist")
	}
	return nil
}

/// Signup
///  1 Create invitation code according to application configuration
///  2 Create user
///  3 Create registration invitation
func (h *Handler) Signup(ctx context.Context) (info *usermwpb.User, err error) {
	signupHandler := &signupHandler{
		Handler: h,
	}

	h.UserID = uuid.NewString()

	if err := signupHandler.checkUser(ctx); err != nil {
		return nil, err
	}

	inviterID, err := h.CheckInvitationCode(ctx)
	if err != nil {
		return nil, err
	}

	signupHandler.inviterID = inviterID

	if err := h.VerifyUserCode(ctx, basetypes.UsedFor_Signup); err != nil {
		return nil, err
	}

	if err := signupHandler.getDefaultRole(ctx); err != nil {
		return nil, err
	}

	sagaDispose := dtmcli.NewSagaDispose(dtmimp.TransOptions{
		WaitResult:     true,
		RequestTimeout: signupHandler.RequestTimeoutSeconds,
	})
	signupHandler.withCreateInvitationCode(sagaDispose)
	signupHandler.withCreateUser(sagaDispose)
	signupHandler.withCreateRegistrationInvitation(sagaDispose)

	if err := dtmcli.WithSaga(ctx, sagaDispose); err != nil {
		return nil, err
	}

	/// TODO: if newbie has coupon, send event to allocate coupon, and we don't care about allocate result

	return h.GetUser(ctx)
}
