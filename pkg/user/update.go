//nolint:dupl
package user

import (
	"context"
	"fmt"

	appusertypes "github.com/NpoolPlatform/message/npool/basetypes/appuser/v1"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	recoverycodemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user/recoverycode"
	ivcodemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/invitationcode"

	usercodemwcli "github.com/NpoolPlatform/basal-middleware/pkg/client/usercode"
	usercodemwpb "github.com/NpoolPlatform/message/npool/basal/mw/v1/usercode"

	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	recoverycodemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user/recoverycode"
	ivcodemwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/invitation/invitationcode"
	regmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/invitation/registration"

	applangmwcli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"
	regmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/registration"
	applangmwpb "github.com/NpoolPlatform/message/npool/g11n/mw/v1/applang"

	sendmwpb "github.com/NpoolPlatform/message/npool/third/mw/v1/send"
	sendmwcli "github.com/NpoolPlatform/third-middleware/pkg/client/send"

	tmplmwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template"
	tmplmwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/template"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
)

type updateHandler struct {
	*Handler
	targetID     *uint32
	targetUser   *usermwpb.User
	origUser     *usermwpb.User
	recoveryCode *recoverycodemwpb.RecoveryCode
}

func (h *updateHandler) verifyOldPasswordHash(ctx context.Context) error {
	if h.OldPasswordHash == nil {
		return nil
	}

	if _, err := usermwcli.VerifyUser(
		ctx,
		*h.AppID,
		*h.UserID,
		*h.OldPasswordHash,
	); err != nil {
		return err
	}

	return nil
}

func (h *updateHandler) getUser(ctx context.Context) error {
	info, err := usermwcli.GetUser(ctx, *h.AppID, *h.UserID)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("update: invalid user: app_id=%v, user_id=%v", *h.AppID, *h.UserID)
	}

	h.User = info
	h.EntID = &info.EntID
	h.UserID = &info.EntID
	h.ID = &info.ID

	return nil
}

func (h *updateHandler) getTargetUser(ctx context.Context) error {
	info, err := usermwcli.GetUser(ctx, *h.AppID, *h.TargetUserID)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("update: invalid user: app_id=%v, user_id=%v", *h.AppID, *h.UserID)
	}
	h.targetID = &info.ID
	h.targetUser = info
	return nil
}

func (h *updateHandler) shouldVerifyNewCode() bool {
	if h.NewAccount != nil {
		return true
	}
	if h.NewAccountType != nil {
		switch *h.NewAccountType {
		case basetypes.SignMethod_Google:
			fallthrough //nolint
		case basetypes.SignMethod_Email:
			fallthrough //nolint
		case basetypes.SignMethod_Mobile:
			return true
		}
	}
	return false
}

func (h *updateHandler) shouldVerifyOldCode() bool {
	if h.shouldVerifyNewCode() {
		return true
	}
	if h.PasswordHash != nil {
		return true
	}
	return false
}

func (h *updateHandler) verifyOldAccountCode(ctx context.Context) error {
	if !h.shouldVerifyOldCode() {
		return nil
	}
	if h.VerificationCode == nil {
		return fmt.Errorf("invalid verification code")
	}
	if h.AccountType == nil {
		return fmt.Errorf("invalid account type")
	}
	account := ""
	if *h.AccountType == basetypes.SignMethod_Google {
		account = h.User.GoogleSecret
	} else if h.Account != nil {
		account = *h.Account
	}
	return usercodemwcli.VerifyUserCode(ctx, &usercodemwpb.VerifyUserCodeRequest{
		Prefix:      basetypes.Prefix_PrefixUserCode.String(),
		AppID:       *h.AppID,
		Account:     account,
		AccountType: *h.AccountType,
		UsedFor:     basetypes.UsedFor_Update,
		Code:        *h.VerificationCode,
	})
}

func (h *updateHandler) verifyNewAccountCode(ctx context.Context) error {
	if !h.shouldVerifyNewCode() {
		return nil
	}
	if h.NewAccountType == nil {
		return fmt.Errorf("invalid account type")
	}
	if h.NewVerificationCode == nil {
		return fmt.Errorf("invalid new verification code")
	}
	account := ""
	if *h.NewAccountType == basetypes.SignMethod_Google {
		account = h.User.GoogleSecret
	} else if h.NewAccount != nil {
		account = *h.NewAccount
	}
	return usercodemwcli.VerifyUserCode(ctx, &usercodemwpb.VerifyUserCodeRequest{
		Prefix:      basetypes.Prefix_PrefixUserCode.String(),
		AppID:       *h.AppID,
		Account:     account,
		AccountType: *h.NewAccountType,
		UsedFor:     basetypes.UsedFor_Update,
		Code:        *h.NewVerificationCode,
	})
}

func (h *updateHandler) updateUser(ctx context.Context) error {
	req := &usermwpb.UserReq{
		ID:                 h.ID,
		EntID:              h.UserID,
		AppID:              h.AppID,
		Username:           h.Username,
		AddressFields:      h.AddressFields,
		Gender:             h.Gender,
		PostalCode:         h.PostalCode,
		Age:                h.Age,
		Birthday:           h.Birthday,
		Avatar:             h.Avatar,
		Organization:       h.Organization,
		FirstName:          h.FirstName,
		LastName:           h.LastName,
		IDNumber:           h.IDNumber,
		SigninVerifyType:   h.SigninVerifyType,
		PasswordHash:       h.PasswordHash,
		KolConfirmed:       h.KolConfirmed,
		SelectedLangID:     h.SelectedLangID,
		GoogleSecret:       h.GoogleSecret,
		GoogleAuthVerified: h.GoogleAuthVerified,
		Banned:             h.Banned,
		BanMessage:         h.BanMessage,
		Kol:                h.Kol,
		EmailAddress:       h.EmailAddress,
	}
	if h.NewAccountType != nil {
		if *h.NewAccountType != basetypes.SignMethod_Google && h.NewAccount == nil {
			return fmt.Errorf("invalid account")
		}
		switch *h.NewAccountType {
		case basetypes.SignMethod_Google:
			verified := true
			req.GoogleAuthVerified = &verified
		case basetypes.SignMethod_Email:
			req.EmailAddress = h.NewAccount
		case basetypes.SignMethod_Mobile:
			req.PhoneNO = h.NewAccount
		}
	}

	info, err := usermwcli.UpdateUser(ctx, req)
	if err != nil {
		return err
	}

	h.User = info
	return nil
}

func (h *updateHandler) updateCache(ctx context.Context) error {
	if h.UpdateCacheMode == nil {
		return nil
	}

	meta, err := h.QueryCache(ctx)
	if err != nil {
		return err
	}

	switch *h.UpdateCacheMode {
	case RequiredUpdateCache:
	case DeleteCacheIfExist:
		if meta == nil {
			return nil
		}
		h.Metadata = meta
		if err := h.DeleteCache(ctx); err != nil {
			return err
		}
		return nil
	case UpdateCacheIfExist:
		if meta == nil {
			return nil
		}
	case DontUpdateCache:
		return nil
	default:
		return fmt.Errorf("invalid updatecachemode")
	}

	if err := h.UpdateCache(ctx); err != nil {
		return err
	}
	meta, err = h.QueryCache(ctx)
	if err != nil {
		return err
	}
	h.User = meta.User
	return nil
}

func (h *Handler) UpdateUser(ctx context.Context) (*usermwpb.User, error) {
	handler := &updateHandler{
		Handler: h,
	}

	if h.ID != nil && h.EntID != nil {
		if err := h.ExistUser(ctx); err != nil {
			return nil, err
		}
	}
	if h.UserID != nil && h.AppID != nil {
		if err := h.ExistUserInApp(ctx); err != nil {
			return nil, err
		}
	}

	notif1 := &notifHandler{
		Handler: h,
	}
	notif1.getUsedFor()

	if err := handler.CheckNewAccount(ctx); err != nil {
		return nil, err
	}
	if err := handler.verifyOldPasswordHash(ctx); err != nil {
		return nil, err
	}
	if err := handler.getUser(ctx); err != nil {
		return nil, err
	}
	if err := handler.verifyOldAccountCode(ctx); err != nil {
		return nil, err
	}
	if err := handler.verifyNewAccountCode(ctx); err != nil {
		return nil, err
	}

	handler.origUser = handler.User

	if err := handler.updateUser(ctx); err != nil {
		return nil, err
	}
	if err := handler.updateCache(ctx); err != nil {
		return nil, err
	}

	// Generate Notif
	notif1.generateNotif(ctx)

	return h.GetUser(ctx)
}

func (h *updateHandler) getAccountUser(ctx context.Context) error {
	conds := &usermwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}
	switch *h.AccountType {
	case basetypes.SignMethod_Email:
		conds.EmailAddress = &basetypes.StringVal{Op: cruder.EQ, Value: *h.Account}
	case basetypes.SignMethod_Mobile:
		conds.PhoneNO = &basetypes.StringVal{Op: cruder.EQ, Value: *h.Account}
	default:
		return fmt.Errorf("invalid account type")
	}

	info, err := usermwcli.GetUserOnly(ctx, conds)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("invalid user: conds=%v", conds)
	}

	h.UserID = &info.EntID
	h.User = info
	h.ID = &info.ID
	h.EntID = &info.EntID

	return nil
}

func (h *updateHandler) verifyAccountCode(ctx context.Context) error {
	if h.VerificationCode == nil {
		return nil
	}
	if h.Account == nil || h.AccountType == nil {
		return fmt.Errorf("invalid account")
	}
	if h.VerificationCode == nil {
		return fmt.Errorf("invalid verification code")
	}
	return usercodemwcli.VerifyUserCode(ctx, &usercodemwpb.VerifyUserCodeRequest{
		Prefix:      basetypes.Prefix_PrefixUserCode.String(),
		AppID:       *h.AppID,
		Account:     *h.Account,
		AccountType: *h.AccountType,
		UsedFor:     basetypes.UsedFor_Update,
		Code:        *h.VerificationCode,
	})
}

func (h *updateHandler) verifyRecoveryCode(ctx context.Context) error {
	if h.RecoveryCode == nil {
		return nil
	}
	code, err := recoverycodemwcli.GetRecoveryCodeOnly(ctx, &recoverycodemwpb.Conds{
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		UserID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID},
		Code:   &basetypes.StringVal{Op: cruder.EQ, Value: *h.RecoveryCode},
		Used:   &basetypes.BoolVal{Op: cruder.EQ, Value: false},
	})
	if err != nil {
		return err
	}
	if code == nil {
		return fmt.Errorf("invalid recovery code")
	}
	h.recoveryCode = code
	return nil
}

func (h *updateHandler) expireRecoveryCode(ctx context.Context) error {
	if h.recoveryCode == nil {
		return nil
	}
	used := true
	if _, err := recoverycodemwcli.UpdateRecoveryCode(ctx, &recoverycodemwpb.RecoveryCodeReq{
		ID:   &h.recoveryCode.ID,
		Used: &used,
	}); err != nil {
		return err
	}
	return nil
}

func (h *Handler) ResetUser(ctx context.Context) error {
	if h.VerificationCode == nil && h.RecoveryCode == nil {
		return fmt.Errorf("need verification code or recovery code")
	}
	handler := &updateHandler{
		Handler: h,
	}
	if err := handler.getAccountUser(ctx); err != nil {
		return err
	}
	if err := h.VerifyResetUserLink(ctx); err != nil {
		return err
	}
	if err := handler.verifyAccountCode(ctx); err != nil {
		return err
	}
	if err := handler.verifyRecoveryCode(ctx); err != nil {
		return err
	}
	if _, err := usermwcli.UpdateUser(ctx, &usermwpb.UserReq{
		ID:           h.ID,
		EntID:        h.UserID,
		AppID:        h.AppID,
		PasswordHash: h.PasswordHash,
	}); err != nil {
		return err
	}
	if err := handler.expireRecoveryCode(ctx); err != nil {
		return err
	}
	updateCacheMode := DeleteCacheIfExist
	h.UpdateCacheMode = &updateCacheMode
	if err := handler.updateCache(ctx); err != nil {
		return err
	}
	if err := h.DeleteResetUserLink(ctx); err != nil {
		return err
	}
	return nil
}

func (h *updateHandler) verifyRegistrationInvitation(ctx context.Context) error {
	if h.TargetUserID == nil {
		return fmt.Errorf("invalid target userid")
	}

	reg, err := regmwcli.GetRegistrationOnly(ctx, &regmwpb.Conds{
		AppID:     &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		InviterID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID},
		InviteeID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.TargetUserID},
	})
	if err != nil {
		return nil
	}
	if reg == nil {
		return fmt.Errorf("invalid registration invitation")
	}
	return nil
}

func (h *updateHandler) tryCreateInvitationCode(ctx context.Context) error {
	info, err := ivcodemwcli.GetInvitationCodeOnly(ctx, &ivcodemwpb.Conds{
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		UserID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.TargetUserID},
	},
	)
	if err != nil {
		return err
	}
	if info != nil {
		return nil
	}

	_, err = ivcodemwcli.CreateInvitationCode(ctx, &ivcodemwpb.InvitationCodeReq{
		AppID:  h.AppID,
		UserID: h.TargetUserID,
	})
	if err != nil {
		return err
	}
	return nil
}

// TODO: move this to pubsub message
func (h *updateHandler) sendKolNotification(ctx context.Context) {
	if !*h.Kol {
		return
	}

	lang, err := applangmwcli.GetLangOnly(ctx, &applangmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		Main:  &basetypes.BoolVal{Op: cruder.EQ, Value: true},
	})
	if err != nil {
		logger.Sugar().Errorw("sendKolNotification", "Error", err)
		return
	}
	if lang == nil {
		logger.Sugar().Warnw("sendKolNotification", "Error", "Main AppLang not exist")
		return
	}

	info, err := tmplmwcli.GenerateText(ctx, &tmplmwpb.GenerateTextRequest{
		AppID:     *h.AppID,
		LangID:    lang.LangID,
		Channel:   basetypes.NotifChannel_ChannelEmail,
		EventType: basetypes.UsedFor_CreateInvitationCode,
	})
	if err != nil {
		logger.Sugar().Errorw("sendKolNotification", "Error", err)
		return
	}
	if info == nil {
		logger.Sugar().Warnw("sendKolNotification", "Error", "Cannot generate text")
		return
	}

	err = sendmwcli.SendMessage(ctx, &sendmwpb.SendMessageRequest{
		Subject:     info.Subject,
		Content:     info.Content,
		From:        info.From,
		To:          h.targetUser.EmailAddress,
		ToCCs:       info.ToCCs,
		ReplyTos:    info.ReplyTos,
		AccountType: basetypes.SignMethod_Email,
	})
	if err != nil {
		logger.Sugar().Errorw("sendKolNotification", "Error", err)
	}
}

func (h *Handler) UpdateUserKol(ctx context.Context) (*usermwpb.User, error) {
	handler := &updateHandler{
		Handler: h,
	}

	if h.CheckInvitation != nil && *h.CheckInvitation {
		if h.UserID == nil {
			return nil, fmt.Errorf("invalid userid")
		}
		if err := handler.verifyRegistrationInvitation(ctx); err != nil {
			return nil, err
		}
	}
	if err := handler.getUser(ctx); err != nil {
		return nil, err
	}
	if err := handler.getTargetUser(ctx); err != nil {
		return nil, err
	}
	req := &usermwpb.UserReq{
		ID:    handler.targetID,
		EntID: h.TargetUserID,
		AppID: h.AppID,
		Kol:   h.Kol,
	}
	if _, err := usermwcli.UpdateUser(ctx, req); err != nil {
		return nil, err
	}

	if h.Kol != nil && *h.Kol {
		if err := handler.tryCreateInvitationCode(ctx); err != nil {
			return nil, err
		}
		handler.sendKolNotification(ctx)
	}

	return h.GetUser(ctx)
}

func (h *updateHandler) getApp(ctx context.Context) error {
	app, err := appmwcli.GetApp(ctx, *h.AppID)
	if err != nil {
		return err
	}
	if app == nil {
		return fmt.Errorf("invalid app")
	}
	if app.ResetUserMethod != appusertypes.ResetUserMethod_Link {
		return fmt.Errorf("permission denied")
	}
	return nil
}

func (h *Handler) PreResetUser(ctx context.Context) error {
	handler := &updateHandler{
		Handler: h,
	}
	if err := handler.getApp(ctx); err != nil {
		return err
	}
	if err := handler.getAccountUser(ctx); err != nil {
		return err
	}

	link, err := h.CreateResetUserLink(ctx)
	if err != nil {
		return err
	}

	lang, err := applangmwcli.GetLangOnly(ctx, &applangmwpb.Conds{
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		LangID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.LangID},
	})
	if err != nil {
		return nil
	}
	if lang == nil {
		return fmt.Errorf("invalid langid")
	}

	var channel basetypes.NotifChannel
	switch *h.AccountType {
	case basetypes.SignMethod_Email:
		channel = basetypes.NotifChannel_ChannelEmail
	case basetypes.SignMethod_Mobile:
		channel = basetypes.NotifChannel_ChannelSMS
	default:
		return fmt.Errorf("invalid channel")
	}

	info, err := tmplmwcli.GenerateText(ctx, &tmplmwpb.GenerateTextRequest{
		AppID:     *h.AppID,
		LangID:    lang.LangID,
		Channel:   channel,
		EventType: basetypes.UsedFor_ResetPassword,
		Vars: &tmplmwpb.TemplateVars{
			Message: &link,
		},
	})
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("generate text failed")
	}

	return sendmwcli.SendMessage(ctx, &sendmwpb.SendMessageRequest{
		Subject:     info.Subject,
		Content:     info.Content,
		From:        info.From,
		To:          *h.Account,
		ToCCs:       info.ToCCs,
		ReplyTos:    info.ReplyTos,
		AccountType: *h.AccountType,
	})
}
