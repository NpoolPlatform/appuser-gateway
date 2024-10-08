package user

import (
	"context"
	"fmt"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role"

	appusersvcname "github.com/NpoolPlatform/appuser-middleware/pkg/servicename"
	inspiremwsvcname "github.com/NpoolPlatform/inspire-middleware/pkg/servicename"

	rolemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	eventmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event"
	ivcodemwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/invitation/invitationcode"
	registrationmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/invitation/registration"

	dtmcli "github.com/NpoolPlatform/dtm-cluster/pkg/dtm"
	"github.com/dtm-labs/dtm/client/dtmcli/dtmimp"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/go-service-framework/pkg/pubsub"
	redis2 "github.com/NpoolPlatform/go-service-framework/pkg/redis"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	"github.com/google/uuid"
)

type signupHandler struct {
	*Handler
	inviterID   *string
	defaultRole *rolemwpb.Role
}

func (h *signupHandler) withCreateInvitationCode(dispose *dtmcli.SagaDispose) {
	if h.App.CreateInvitationCodeWhen != basetypes.CreateInvitationCodeWhen_Registration {
		return
	}

	id := uuid.NewString()
	req := &ivcodemwpb.InvitationCodeReq{
		EntID:  &id,
		AppID:  h.AppID,
		UserID: h.UserID,
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
		EntID:        h.UserID,
		AppID:        h.AppID,
		EmailAddress: h.EmailAddress,
		PhoneNO:      h.PhoneNO,
		PasswordHash: h.PasswordHash,
		RoleIDs:      []string{h.defaultRole.EntID},
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

	id := uuid.NewString()
	req := &registrationmwpb.RegistrationReq{
		EntID:     &id,
		AppID:     h.AppID,
		InviterID: h.inviterID,
		InviteeID: h.UserID,
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
	role, err := rolemwcli.GetRoleOnly(ctx, &rolemwpb.Conds{
		AppID:   &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		Default: &basetypes.BoolVal{Op: cruder.EQ, Value: true},
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

//nolint:dupl
func (h *signupHandler) rewardSignup() {
	if err := pubsub.WithPublisher(func(publisher *pubsub.Publisher) error {
		req := &eventmwpb.CalcluateEventRewardsRequest{
			AppID:       *h.AppID,
			UserID:      *h.UserID,
			EventType:   basetypes.UsedFor_Signup,
			Consecutive: 1,
		}
		return publisher.Update(
			basetypes.MsgID_CalculateEventRewardReq.String(),
			nil,
			nil,
			nil,
			req,
		)
	}); err != nil {
		logger.Sugar().Errorw(
			"rewardSignup",
			"AppID", *h.AppID,
			"UserID", h.UserID,
			"Account", h.Account,
			"AccountType", h.AccountType,
			"Error", err,
		)
	}
}

//nolint:dupl
func (h *signupHandler) rewardAffiliateSignup() {
	if err := pubsub.WithPublisher(func(publisher *pubsub.Publisher) error {
		req := &eventmwpb.CalcluateEventRewardsRequest{
			AppID:       *h.AppID,
			UserID:      *h.UserID,
			EventType:   basetypes.UsedFor_AffiliateSignup,
			Consecutive: 1,
		}
		return publisher.Update(
			basetypes.MsgID_CalculateEventRewardReq.String(),
			nil,
			nil,
			nil,
			req,
		)
	}); err != nil {
		logger.Sugar().Errorw(
			"rewardAffiliateSignup",
			"AppID", *h.AppID,
			"UserID", h.UserID,
			"Account", h.Account,
			"AccountType", h.AccountType,
			"Error", err,
		)
	}
}

// Signup
//  1 Create invitation code according to application configuration
//  2 Create user
//  3 Create registration invitation
//  4 Reward user's signup event and affiliate signup event
func (h *Handler) Signup(ctx context.Context) (info *usermwpb.User, err error) {
	signupHandler := &signupHandler{
		Handler: h,
	}

	key := fmt.Sprintf(
		"%v:%v:%v:%v",
		basetypes.Prefix_PrefixUserAccount,
		*h.AppID,
		basetypes.UsedFor_Signup,
		*h.Account,
	)
	if err := redis2.TryLock(key, 0); err != nil {
		return nil, err
	}
	defer func() {
		_ = redis2.Unlock(key)
	}()

	if err := signupHandler.CheckUser(ctx); err != nil {
		return nil, err
	}

	inviterID, err := h.CheckInvitationCode(ctx)
	if err != nil {
		return nil, err
	}

	signupHandler.inviterID = inviterID
	id := uuid.NewString()
	signupHandler.UserID = &id

	if err := h.VerifyUserCode(ctx, basetypes.UsedFor_Signup); err != nil {
		return nil, err
	}

	if err := signupHandler.getDefaultRole(ctx); err != nil {
		return nil, err
	}

	sagaDispose := dtmcli.NewSagaDispose(dtmimp.TransOptions{
		WaitResult:     true,
		TimeoutToFail:  60,
		RequestTimeout: *signupHandler.RequestTimeoutSeconds,
	})
	signupHandler.withCreateInvitationCode(sagaDispose)
	signupHandler.withCreateUser(sagaDispose)
	signupHandler.withCreateRegistrationInvitation(sagaDispose)

	if err := dtmcli.WithSaga(ctx, sagaDispose); err != nil {
		return nil, err
	}

	signupHandler.rewardSignup()
	signupHandler.rewardAffiliateSignup()

	return h.GetUser(ctx)
}
