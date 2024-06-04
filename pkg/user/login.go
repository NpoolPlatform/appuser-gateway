package user

import (
	"context"
	"fmt"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	usercodemwcli "github.com/NpoolPlatform/basal-middleware/pkg/client/usercode"
	eventmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/event"
	ivcodemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/invitationcode"
	taskusermwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/task/user"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	loginhispb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user/login/history"
	usercodemwpb "github.com/NpoolPlatform/message/npool/basal/mw/v1/usercode"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	eventmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event"
	ivcodemwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/invitation/invitationcode"
	taskusermwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/task/user"
	thirdmwcli "github.com/NpoolPlatform/third-middleware/pkg/client/verify"

	"github.com/NpoolPlatform/go-service-framework/pkg/pubsub"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	"github.com/google/uuid"
)

type loginHandler struct {
	*Handler
}

func (h *loginHandler) notifyLogin(loginType basetypes.LoginType) {
	clientIP := h.Metadata.ClientIP.String()

	if err := pubsub.WithPublisher(func(publisher *pubsub.Publisher) error {
		req := &loginhispb.HistoryReq{
			AppID:     h.AppID,
			UserID:    h.UserID,
			ClientIP:  &clientIP,
			UserAgent: &h.Metadata.UserAgent,
			LoginType: &loginType,
		}
		return publisher.Update(
			basetypes.MsgID_CreateLoginHistoryReq.String(),
			nil,
			nil,
			nil,
			req,
		)
	}); err != nil {
		logger.Sugar().Errorw(
			"notifyLogin",
			"AppID", *h.AppID,
			"UserID", h.UserID,
			"Account", h.Account,
			"AccountType", h.AccountType,
			"Error", err,
		)
	}
}

func (h *loginHandler) checkLoginReward(ctx context.Context) {
	loginEvent, err := eventmwcli.GetEventOnly(ctx, &eventmwpb.Conds{
		AppID:     &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EventType: &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(basetypes.UsedFor_Signin)},
	})
	if err != nil {
		logger.Sugar().Errorw(
			"checkLoginReward",
			"AppID", *h.AppID,
			"UserID", h.UserID,
			"EventType", basetypes.UsedFor_Signin,
			"Account", h.Account,
			"AccountType", h.AccountType,
			"Error", err,
		)
		return
	}
	if loginEvent == nil {
		return
	}
	now := uint32(time.Now().Unix())
	coolDownDuration := 4 * 60 * 60
	coolDownTime := now - uint32(coolDownDuration)

	exist, err := taskusermwcli.ExistTaskUserConds(ctx, &taskusermwpb.Conds{
		AppID:     &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		UserID:    &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID},
		EventID:   &basetypes.StringVal{Op: cruder.EQ, Value: loginEvent.EntID},
		CreatedAt: &basetypes.Uint32Val{Op: cruder.GTE, Value: coolDownTime},
	})
	if err != nil {
		logger.Sugar().Errorw(
			"checkLoginReward",
			"AppID", *h.AppID,
			"UserID", h.UserID,
			"EventID", loginEvent.EntID,
			"CreatedAt", coolDownTime,
			"Account", h.Account,
			"AccountType", h.AccountType,
			"Error", err,
		)
		return
	}
	if exist {
		return
	}

	h.rewardSignin()
}

func (h *loginHandler) checkConsecutiveLoginReward(ctx context.Context) {
	loginEvent, err := eventmwcli.GetEventOnly(ctx, &eventmwpb.Conds{
		AppID:     &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EventType: &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(basetypes.UsedFor_ConsecutiveLogin)},
	})
	if err != nil {
		logger.Sugar().Errorw(
			"checkConsecutiveLoginReward",
			"AppID", *h.AppID,
			"UserID", h.UserID,
			"EventType", basetypes.UsedFor_ConsecutiveLogin,
			"Account", h.Account,
			"AccountType", h.AccountType,
			"Error", err,
		)
		return
	}
	if loginEvent == nil {
		return
	}
	now := time.Now()
	location := now.Location()
	midnight := uint32(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location).Unix())

	exist, err := taskusermwcli.ExistTaskUserConds(ctx, &taskusermwpb.Conds{
		AppID:     &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		UserID:    &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID},
		EventID:   &basetypes.StringVal{Op: cruder.EQ, Value: loginEvent.EntID},
		CreatedAt: &basetypes.Uint32Val{Op: cruder.GTE, Value: midnight},
	})
	if err != nil {
		logger.Sugar().Errorw(
			"checkConsecutiveLoginReward",
			"AppID", *h.AppID,
			"UserID", h.UserID,
			"EventID", loginEvent.EntID,
			"CreatedAt", midnight,
			"Account", h.Account,
			"AccountType", h.AccountType,
			"Error", err,
		)
		return
	}
	if exist {
		return
	}

	h.rewardConsecutiveLogin()
}

//nolint:dupl
func (h *loginHandler) rewardSignin() {
	if err := pubsub.WithPublisher(func(publisher *pubsub.Publisher) error {
		req := &eventmwpb.CalcluateEventRewardsRequest{
			AppID:       *h.AppID,
			UserID:      *h.UserID,
			EventType:   basetypes.UsedFor_Signin,
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
			"rewardSignin",
			"AppID", *h.AppID,
			"UserID", h.UserID,
			"Account", h.Account,
			"AccountType", h.AccountType,
			"Error", err,
		)
	}
}

//nolint:dupl
func (h *loginHandler) rewardConsecutiveLogin() {
	if err := pubsub.WithPublisher(func(publisher *pubsub.Publisher) error {
		req := &eventmwpb.CalcluateEventRewardsRequest{
			AppID:       *h.AppID,
			UserID:      *h.UserID,
			EventType:   basetypes.UsedFor_ConsecutiveLogin,
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
			"rewardConsecutiveLogin",
			"AppID", *h.AppID,
			"UserID", h.UserID,
			"Account", h.Account,
			"AccountType", h.AccountType,
			"Error", err,
		)
	}
}

func (h *loginHandler) verifyRecaptcha(ctx context.Context) error {
	if h.ManMachineSpec == nil {
		return fmt.Errorf("invalid manmachinespec")
	}
	switch h.App.RecaptchaMethod {
	case basetypes.RecaptchaMethod_GoogleRecaptchaV3:
		return thirdmwcli.VerifyGoogleRecaptchaV3(ctx, *h.ManMachineSpec)
	case basetypes.RecaptchaMethod_NoRecaptcha:
	default:
	}
	return nil
}

func (h *loginHandler) verifyAccount(ctx context.Context) error {
	if h.Account == nil || h.AccountType == nil || h.PasswordHash == nil {
		return fmt.Errorf("invalid account or password")
	}
	info, err := usermwcli.VerifyAccount(
		ctx,
		*h.AppID,
		*h.Account,
		*h.AccountType,
		*h.PasswordHash,
	)
	if err != nil {
		return err
	}
	if info == nil || info.Banned {
		return fmt.Errorf("invalid user")
	}
	if _, err = uuid.Parse(info.EntID); err != nil {
		return err
	}
	h.UserID = &info.EntID
	h.User = info
	return nil
}

func (h *loginHandler) prepareMetadata(ctx context.Context) error {
	meta, err := MetadataFromContext(ctx)
	if err != nil {
		return err
	}
	meta.AppID = uuid.MustParse(*h.AppID)
	meta.Account = *h.Account
	meta.AccountType = h.AccountType.String()
	meta.UserID = uuid.MustParse(*h.UserID)
	h.Metadata = meta
	return nil
}

func (h *loginHandler) formalizeUser() {
	h.User.Logined = true
	h.User.LoginAccount = *h.Account
	h.User.LoginAccountType = *h.AccountType
	h.User.LoginToken = *h.Token
	h.User.LoginClientIP = h.Metadata.ClientIP.String()
	h.User.LoginClientUserAgent = h.Metadata.UserAgent

	if !h.App.SigninVerifyEnable {
		h.User.LoginVerified = true
	}

	h.User.GoogleOTPAuth = fmt.Sprintf(
		"otpauth://totp/%s?secret=%s",
		*h.Account,
		h.User.GoogleSecret,
	)

	h.Metadata.User = h.User
}

func (h *loginHandler) getInvitationCode(ctx context.Context) error {
	code, err := ivcodemwcli.GetInvitationCodeOnly(
		ctx,
		&ivcodemwpb.Conds{
			AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			UserID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID},
		},
	)
	if err != nil {
		return err
	}
	if code == nil {
		return nil
	}

	h.User.InvitationCode = &code.InvitationCode
	return nil
}

func (h *Handler) Login(ctx context.Context) (info *usermwpb.User, err error) {
	handler := &loginHandler{
		Handler: h,
	}

	if err := handler.verifyRecaptcha(ctx); err != nil {
		return nil, err
	}

	if err := handler.verifyAccount(ctx); err != nil {
		return nil, err
	}
	if err := handler.prepareMetadata(ctx); err != nil {
		return nil, err
	}
	token, err := createToken(h.Metadata)
	if err != nil {
		return nil, err
	}
	h.Token = &token
	handler.formalizeUser()
	if err := handler.getInvitationCode(ctx); err != nil {
		return nil, err
	}
	if err := h.CreateCache(ctx); err != nil {
		return nil, err
	}

	handler.notifyLogin(basetypes.LoginType_FreshLogin)
	handler.checkLoginReward(ctx)
	handler.checkConsecutiveLoginReward(ctx)

	return h.User, nil
}

func (h *loginHandler) getThirdUser(ctx context.Context) error {
	info, err := usermwcli.GetUserOnly(
		ctx,
		&usermwpb.Conds{
			AppID:            &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			ThirdPartyUserID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.Account},
			ThirdPartyID:     &basetypes.StringVal{Op: cruder.EQ, Value: *h.ThirdPartyID},
		},
	)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("user not exist")
	}
	h.User = info
	return nil
}

func (h *Handler) ThirdLogin(ctx context.Context) (info *usermwpb.User, err error) {
	handler := &loginHandler{
		Handler: h,
	}

	if err := handler.getThirdUser(ctx); err != nil {
		return nil, err
	}

	if err := handler.prepareMetadata(ctx); err != nil {
		return nil, err
	}
	token, err := createToken(h.Metadata)
	if err != nil {
		return nil, err
	}
	h.Token = &token
	handler.formalizeUser()
	if err := handler.getInvitationCode(ctx); err != nil {
		return nil, err
	}
	if err := h.CreateCache(ctx); err != nil {
		return nil, err
	}
	handler.notifyLogin(basetypes.LoginType_FreshLogin)
	handler.checkLoginReward(ctx)
	handler.checkConsecutiveLoginReward(ctx)
	return h.User, nil
}

func (h *loginHandler) mustQueryMetadata(ctx context.Context) (err error) {
	h.Metadata, err = h.QueryCache(ctx)
	if err != nil {
		return err
	}
	if h.Metadata == nil || h.Metadata.User == nil {
		return fmt.Errorf("metadata not exist")
	}

	h.User = h.Metadata.User

	return nil
}

func (h *loginHandler) verifyUserCode(ctx context.Context) error {
	if h.AccountType == nil {
		return fmt.Errorf("invalid account type")
	}
	if h.VerificationCode == nil {
		return fmt.Errorf("invalid verification code")
	}

	switch *h.AccountType {
	case basetypes.SignMethod_Email:
		fallthrough //nolint
	case basetypes.SignMethod_Mobile:
		if h.Account == nil {
			return fmt.Errorf("invalid account")
		}
	case basetypes.SignMethod_Google:
	default:
		return fmt.Errorf("not supported")
	}

	switch *h.AccountType {
	case basetypes.SignMethod_Email:
		if *h.Account != h.Metadata.User.EmailAddress {
			return fmt.Errorf("invalid account")
		}
	case basetypes.SignMethod_Mobile:
		if *h.Account != h.Metadata.User.PhoneNO {
			return fmt.Errorf("invalid account")
		}
	case basetypes.SignMethod_Google:
		h.Account = &h.Metadata.User.GoogleSecret
	}

	if err := usercodemwcli.VerifyUserCode(ctx, &usercodemwpb.VerifyUserCodeRequest{
		Prefix:      basetypes.Prefix_PrefixUserCode.String(),
		AppID:       *h.AppID,
		Account:     *h.Account,
		AccountType: *h.AccountType,
		UsedFor:     basetypes.UsedFor_Signin,
		Code:        *h.VerificationCode,
	}); err != nil {
		return err
	}

	return nil
}

func (h *Handler) LoginVerify(ctx context.Context) (*usermwpb.User, error) {
	if h.Token == nil {
		return nil, fmt.Errorf("invalid token")
	}
	handler := &loginHandler{
		Handler: h,
	}
	if err := handler.mustQueryMetadata(ctx); err != nil {
		return nil, err
	}
	if err := verifyToken(h.Metadata, *h.Token); err != nil {
		return nil, err
	}
	if err := handler.verifyUserCode(ctx); err != nil {
		return nil, err
	}
	h.User.LoginVerified = true
	h.Metadata.User.LoginVerified = true
	if err := h.CreateCache(ctx); err != nil {
		return nil, err
	}
	return h.User, nil
}

func (h *Handler) Logined(ctx context.Context) (*usermwpb.User, error) {
	if h.Token == nil {
		return nil, fmt.Errorf("invalid token")
	}

	handler := &loginHandler{
		Handler: h,
	}

	if err := handler.mustQueryMetadata(ctx); err != nil {
		return nil, err
	}
	if err := verifyToken(h.Metadata, *h.Token); err != nil {
		return nil, err
	}
	if !h.User.LoginVerified {
		return nil, fmt.Errorf("not verified")
	}
	if err := h.CreateCache(ctx); err != nil {
		return nil, err
	}

	return h.User, nil
}

func (h *Handler) Logout(ctx context.Context) (*usermwpb.User, error) {
	if h.Token == nil {
		return nil, fmt.Errorf("invalid token")
	}

	handler := &loginHandler{
		Handler: h,
	}

	if err := handler.mustQueryMetadata(ctx); err != nil {
		return nil, err
	}

	if err := h.DeleteCache(ctx); err != nil {
		return nil, err
	}

	return h.User, nil
}
