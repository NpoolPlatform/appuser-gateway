package user

import (
	"context"
	"fmt"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	ivcodemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/invitationcode"

	usercodemwcli "github.com/NpoolPlatform/basal-middleware/pkg/client/usercode"
	usercodemwpb "github.com/NpoolPlatform/message/npool/basal/mw/v1/usercode"

	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	ivcodemwpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/invitationcode"

	applangmwcli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"
	applangmgrpb "github.com/NpoolPlatform/message/npool/g11n/mgr/v1/applang"

	chanmgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/channel"

	sendmwpb "github.com/NpoolPlatform/message/npool/third/mw/v1/send"
	sendmwcli "github.com/NpoolPlatform/third-middleware/pkg/client/send"

	tmplmwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template"
	tmplmwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/template"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"
)

type updateHandler struct {
	*Handler
}

func (h *updateHandler) verifyOldPasswordHash(ctx context.Context) error {
	if h.OldPasswordHash == nil {
		return nil
	}

	if _, err := usermwcli.VerifyUser(
		ctx,
		h.AppID,
		h.UserID,
		*h.OldPasswordHash,
	); err != nil {
		return err
	}

	return nil
}

func (h *updateHandler) getUser(ctx context.Context) error {
	info, err := usermwcli.GetUser(ctx, h.AppID, h.UserID)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("invalid user")
	}

	h.User = info
	return nil
}

func (h *updateHandler) shouldVerifyNewCode(ctx context.Context) bool {
	if h.NewAccount != nil {
		return true
	}
	if h.NewAccountType != nil && *h.NewAccountType == basetypes.SignMethod_Google {
		return true
	}
	return false
}

func (h *updateHandler) shouldVerifyOldCode(ctx context.Context) bool {
	if h.shouldVerifyNewCode(ctx) {
		return true
	}
	if h.PasswordHash != nil {
		return true
	}
	return false
}

func (h *updateHandler) verifyOldAccountCode(ctx context.Context) error {
	if !h.shouldVerifyOldCode(ctx) {
		return nil
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
	return usercodemwcli.VerifyUserCode(
		ctx,
		&usercodemwpb.VerifyUserCodeRequest{
			Prefix:      basetypes.Prefix_PrefixUserCode.String(),
			AppID:       h.AppID,
			Account:     account,
			AccountType: *h.AccountType,
			UsedFor:     basetypes.UsedFor_Update,
			Code:        h.VerificationCode,
		},
	)
}

func (h *updateHandler) verifyNewAccountCode(ctx context.Context) error {
	if !h.shouldVerifyNewCode(ctx) {
		return nil
	}
	if h.NewAccountType == nil {
		return fmt.Errorf("invalid account type")
	}
	account := ""
	if *h.NewAccountType == basetypes.SignMethod_Google {
		account = h.User.GoogleSecret
	} else if h.NewAccount != nil {
		account = *h.NewAccount
	}
	return usercodemwcli.VerifyUserCode(
		ctx,
		&usercodemwpb.VerifyUserCodeRequest{
			Prefix:      basetypes.Prefix_PrefixUserCode.String(),
			AppID:       h.AppID,
			Account:     account,
			AccountType: *h.NewAccountType,
			UsedFor:     basetypes.UsedFor_Update,
			Code:        h.NewVerificationCode,
		},
	)
}

func (h *updateHandler) apply(ctx context.Context) error {
	req := &usermwpb.UserReq{
		ID:               &h.UserID,
		AppID:            &h.AppID,
		Username:         h.Username,
		AddressFields:    h.AddressFields,
		Gender:           h.Gender,
		PostalCode:       h.PostalCode,
		Age:              h.Age,
		Birthday:         h.Birthday,
		Avatar:           h.Avatar,
		Organization:     h.Organization,
		FirstName:        h.FirstName,
		LastName:         h.LastName,
		IDNumber:         h.IDNumber,
		SigninVerifyType: h.SigninVerifyType,
		PasswordHash:     h.PasswordHash,
		KolConfirmed:     h.KolConfirmed,
	}
	if h.NewAccountType != nil {
		if h.NewAccount == nil {
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

func (h *Handler) UpdateUser(ctx context.Context) (*usermwpb.User, error) {
	handler := &updateHandler{
		Handler: h,
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
	if err := handler.apply(ctx); err != nil {
		return nil, err
	}

	if err := h.UpdateCache(ctx); err != nil {
		return nil, err
	}
	meta, err := h.QueryCache(ctx)
	if err != nil {
		return nil, err
	}

	h.Metadata = meta

	return h.Metadata.User, nil
}

func ResetUser(ctx context.Context, in *npool.ResetUserRequest) error {
	conds := &usermwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: in.GetAppID()},
	}

	switch in.GetAccountType() {
	case basetypes.SignMethod_Email:
		conds.EmailAddress = &basetypes.StringVal{Op: cruder.EQ, Value: in.GetAccount()}
	case basetypes.SignMethod_Mobile:
		conds.PhoneNO = &basetypes.StringVal{Op: cruder.EQ, Value: in.GetAccount()}
	default:
		return fmt.Errorf("invalid account type")
	}

	auser, err := usermwcli.GetUserOnly(ctx, conds)
	if err != nil {
		return err
	}
	if auser == nil {
		return fmt.Errorf("invalid user")
	}

	if err := usercodemwcli.VerifyUserCode(ctx, &usercodemwpb.VerifyUserCodeRequest{
		Prefix:      basetypes.Prefix_PrefixUserCode.String(),
		AppID:       in.GetAppID(),
		Account:     in.GetAccount(),
		AccountType: in.GetAccountType(),
		UsedFor:     basetypes.UsedFor_Update,
		Code:        in.GetVerificationCode(),
	}); err != nil {
		return err
	}

	_, err = usermwcli.UpdateUser(ctx, &usermwpb.UserReq{
		ID:           &auser.ID,
		AppID:        &in.AppID,
		PasswordHash: in.PasswordHash,
	})

	return err
}

func UpdateUserKol(ctx context.Context, in *npool.UpdateUserKolRequest) (*usermwpb.User, error) {
	app, err := appmwcli.GetApp(ctx, in.GetAppID())
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, fmt.Errorf("invalid app")
	}

	req := &usermwpb.UserReq{
		ID:    &in.TargetUserID,
		AppID: &in.AppID,
		Kol:   &in.Kol,
	}

	info, err := usermwcli.UpdateUser(ctx, req)
	if err != nil {
		logger.Sugar().Errorw("UpdateUserKol", "err", err)
		return nil, err
	}

	code, err := ivcodemwcli.GetInvitationCodeOnly(ctx, &ivcodemwpb.Conds{
		AppID:  &commonpb.StringVal{Op: cruder.EQ, Value: in.GetAppID()},
		UserID: &commonpb.StringVal{Op: cruder.EQ, Value: in.GetTargetUserID()},
	})
	if err != nil {
		logger.Sugar().Errorw("UpdateUserKol", "err", err)
		return nil, err
	}
	if code == nil {
		code, err = ivcodemwcli.CreateInvitationCode(ctx, &ivcodemwpb.InvitationCodeReq{
			AppID:  &info.AppID,
			UserID: &info.ID,
		})
		if err != nil {
			return nil, err
		}

		info.InvitationCode = &code.InvitationCode
	}

	if in.GetKol() {
		lang, err := applangmwcli.GetLangOnly(ctx, &applangmgrpb.Conds{
			AppID: &commonpb.StringVal{Op: cruder.EQ, Value: info.AppID},
			Main:  &commonpb.BoolVal{Op: cruder.EQ, Value: true},
		})
		if err != nil {
			logger.Sugar().Errorw("UpdateUserKol", "Error", err)
			return info, nil
		}
		if lang == nil {
			logger.Sugar().Warnw("UpdateUserKol", "Error", "Main AppLang not exist")
			return info, nil
		}

		info1, err := tmplmwcli.GenerateText(ctx, &tmplmwpb.GenerateTextRequest{
			AppID:     info.AppID,
			LangID:    lang.LangID,
			Channel:   chanmgrpb.NotifChannel_ChannelEmail,
			EventType: basetypes.UsedFor_CreateInvitationCode,
		})
		if err != nil {
			logger.Sugar().Errorw("UpdateUserKol", "Error", err)
			return info, nil
		}
		if info1 == nil {
			logger.Sugar().Warnw("UpdateUserKol", "Error", "Cannot generate text")
			return info, nil
		}

		err = sendmwcli.SendMessage(ctx, &sendmwpb.SendMessageRequest{
			Subject:     info1.Subject,
			Content:     info1.Content,
			From:        info1.From,
			To:          info.EmailAddress,
			ToCCs:       info1.ToCCs,
			ReplyTos:    info1.ReplyTos,
			AccountType: basetypes.SignMethod_Email,
		})
		if err != nil {
			logger.Sugar().Errorw("UpdateUserKol", "Error", err)
			return info, nil
		}
	}

	return info, nil
}
