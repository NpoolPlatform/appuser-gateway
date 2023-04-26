package user

import (
	"context"
	"fmt"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	loginhispb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user/login/history"
	ivcodemwpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/invitationcode"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	ivcodemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/invitationcode"
	thirdmwcli "github.com/NpoolPlatform/third-middleware/pkg/client/verify"

	usercodemwcli "github.com/NpoolPlatform/basal-middleware/pkg/client/usercode"
	usercodemwpb "github.com/NpoolPlatform/message/npool/basal/mw/v1/usercode"

	commonpb "github.com/NpoolPlatform/message/npool"

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
			AppID:     &h.AppID,
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
			"AppID", h.AppID,
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
		h.AppID,
		*h.Account,
		*h.AccountType,
		*h.PasswordHash,
	)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("invalid user")
	}
	if _, err = uuid.Parse(info.ID); err != nil {
		return err
	}
	h.UserID = &info.ID
	h.User = info
	return nil
}

func (h *loginHandler) prepareMetadata(ctx context.Context) error {
	meta, err := MetadataFromContext(ctx)
	if err != nil {
		return err
	}
	meta.AppID = uuid.MustParse(h.AppID)
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

	if h.App.SigninVerifyEnable {
		h.User.LoginVerified = true
	}

	h.User.GoogleOTPAuth = fmt.Sprintf(
		"otpauth://totp/%s?secret=%s",
		*h.Account,
		h.User.GoogleSecret,
	)
}

func (h *loginHandler) getInvitationCode(ctx context.Context) error {
	code, err := ivcodemwcli.GetInvitationCodeOnly(
		ctx,
		&ivcodemwpb.Conds{
			AppID:  &commonpb.StringVal{Op: cruder.EQ, Value: h.AppID},
			UserID: &commonpb.StringVal{Op: cruder.EQ, Value: *h.UserID},
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

	if err := usercodemwcli.VerifyUserCode(
		ctx,
		&usercodemwpb.VerifyUserCodeRequest{
			Prefix:      basetypes.Prefix_PrefixUserCode.String(),
			AppID:       h.AppID,
			Account:     *h.Account,
			AccountType: *h.AccountType,
			UsedFor:     basetypes.UsedFor_Signin,
			Code:        *h.VerificationCode,
		},
	); err != nil {
		return err
	}

	h.Metadata.User.LoginVerified = true

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
	if !h.User.LoginVerified {
		return h.Metadata.User, nil
	}
	if err := verifyToken(h.Metadata, *h.Token); err != nil {
		return nil, err
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
