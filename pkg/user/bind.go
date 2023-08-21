package user

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	oauthmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/oauth/appoauththirdparty"
	usercodemwcli "github.com/NpoolPlatform/basal-middleware/pkg/client/usercode"
	oauthmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/oauth/appoauththirdparty"
	usercodemwpb "github.com/NpoolPlatform/message/npool/basal/mw/v1/usercode"
)

type bindHandler struct {
	*Handler
	oldUserInfo   *usermwpb.User
	thirdUserInfo *usermwpb.User
	oauthConf     *oauthmwpb.OAuthThirdParty
}

//nolint:gocyclo
func (h *bindHandler) validate() error {
	if h.UserID == nil {
		return fmt.Errorf("invalid userid")
	}
	if h.AppID == "" {
		return fmt.Errorf("invalid appid")
	}
	if h.Account == nil {
		return fmt.Errorf("invalid account")
	}
	if h.AccountType == nil {
		return fmt.Errorf("invalid accounttype")
	}
	switch *h.AccountType {
	case basetypes.SignMethod_Github:
	case basetypes.SignMethod_Google:
	case basetypes.SignMethod_Facebook:
	case basetypes.SignMethod_Twitter:
	case basetypes.SignMethod_Linkedin:
	case basetypes.SignMethod_Wechat:
	default:
		return fmt.Errorf("invalid accounttype")
	}
	if h.NewAccount == nil {
		return fmt.Errorf("invalid newaccount")
	}
	if h.NewAccountType == nil {
		return fmt.Errorf("invalid newaccounttype")
	}
	switch *h.NewAccountType {
	case basetypes.SignMethod_Email:
	case basetypes.SignMethod_Mobile:
	default:
		return fmt.Errorf("invalid newaccounttype")
	}
	if h.NewVerificationCode == nil {
		return fmt.Errorf("invalid newverificationcode")
	}
	return nil
}

func (h *bindHandler) validUnbindOAuth() error {
	if h.UserID == nil {
		return fmt.Errorf("invalid userid")
	}
	if h.AppID == "" {
		return fmt.Errorf("invalid appid")
	}
	if h.Account == nil {
		return fmt.Errorf("invalid account")
	}
	if h.AccountType == nil {
		return fmt.Errorf("invalid accounttype")
	}
	switch *h.AccountType {
	case basetypes.SignMethod_Github:
	case basetypes.SignMethod_Google:
	case basetypes.SignMethod_Facebook:
	case basetypes.SignMethod_Twitter:
	case basetypes.SignMethod_Linkedin:
	case basetypes.SignMethod_Wechat:
	default:
		return fmt.Errorf("invalid accounttype")
	}
	return nil
}

func (h *bindHandler) getUser(ctx context.Context) error {
	info, err := usermwcli.GetUser(ctx, h.AppID, *h.UserID)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("bind: invalid user: app_id=%v, user_id=%v", h.AppID, *h.UserID)
	}

	h.User = info
	h.oldUserInfo = info
	return nil
}

func (h *bindHandler) getThirdUserInfo(ctx context.Context) error {
	info, err := usermwcli.GetUserOnly(
		ctx,
		&usermwpb.Conds{
			AppID:            &basetypes.StringVal{Op: cruder.EQ, Value: h.AppID},
			ThirdPartyUserID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.Account},
			ThirdPartyID:     &basetypes.StringVal{Op: cruder.EQ, Value: h.oauthConf.ThirdPartyID},
		},
	)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("appuserthirdparty is empty")
	}

	h.thirdUserInfo = info

	return nil
}

func (h *bindHandler) getThirdPartyConf(ctx context.Context) error {
	info, err := oauthmwcli.GetOAuthThirdPartyOnly(
		ctx,
		&oauthmwpb.Conds{
			AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: h.AppID},
			ClientName: &basetypes.Int32Val{Op: cruder.EQ, Value: int32(*h.AccountType)},
		},
	)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("invalid accounttype")
	}
	h.oauthConf = info

	return nil
}

func (h *bindHandler) verifyNewAccountCode(ctx context.Context) error {
	if h.NewAccountType == nil {
		return fmt.Errorf("invalid account type")
	}
	if h.NewVerificationCode == nil {
		return fmt.Errorf("invalid new verification code")
	}
	account := ""
	if h.NewAccount != nil {
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
			Code:        *h.NewVerificationCode,
		},
	)
}

func (h *bindHandler) verifyNewAccount(ctx context.Context) error {
	conds := &usermwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: h.AppID},
	}
	switch *h.NewAccountType {
	case basetypes.SignMethod_Email:
		conds.EmailAddress = &basetypes.StringVal{Op: cruder.EQ, Value: *h.NewAccount}
	case basetypes.SignMethod_Mobile:
		conds.PhoneNO = &basetypes.StringVal{Op: cruder.EQ, Value: *h.NewAccount}
	}

	info, err := usermwcli.GetUserOnly(ctx, conds)
	if err != nil {
		return err
	}
	if info != nil && info.ID != *h.UserID {
		return fmt.Errorf("invalid account")
	}
	return nil
}

func (h *bindHandler) updateUser(ctx context.Context) error {
	req := &usermwpb.UserReq{
		ID:               h.UserID,
		AppID:            &h.AppID,
		EmailAddress:     h.EmailAddress,
		PhoneNO:          h.PhoneNO,
		ThirdPartyID:     h.thirdUserInfo.ThirdPartyID,
		ThirdPartyUserID: h.thirdUserInfo.ThirdPartyUserID,
	}

	if h.NewAccountType != nil {
		if *h.NewAccountType != basetypes.SignMethod_Google && h.NewAccount == nil {
			return fmt.Errorf("invalid account")
		}
		switch *h.NewAccountType {
		case basetypes.SignMethod_Email:
			req.EmailAddress = h.NewAccount
		case basetypes.SignMethod_Mobile:
			req.PhoneNO = h.NewAccount
		}
	}

	_, err := usermwcli.UpdateUser(ctx, req)
	if err != nil {
		return err
	}

	if err := h.getThirdUserInfo(ctx); err != nil {
		return err
	}

	h.User = h.thirdUserInfo
	return nil
}

func (h *bindHandler) updateCache(ctx context.Context) error {
	err := h.UpdateCache(ctx)
	if err != nil {
		return err
	}
	meta, err := h.QueryCache(ctx)
	if err != nil {
		return err
	}
	h.Metadata = meta

	return nil
}

func (h *bindHandler) deleteCache(ctx context.Context) error {
	meta, err := h.QueryCache(ctx)
	if err != nil {
		return err
	}
	if meta == nil || meta.User == nil {
		return fmt.Errorf("cache: invalid user: app_id=%v, user_id=%v", h.AppID, *h.UserID)
	}
	h.Metadata = meta
	if err := h.DeleteCache(ctx); err != nil {
		return err
	}
	return nil
}

func (h *bindHandler) deleteThirdUserInfo(ctx context.Context) error {
	_, err := usermwcli.DeleteThirdUser(ctx, h.AppID, *h.UserID, h.oauthConf.ThirdPartyID, *h.Account)
	if err != nil {
		return nil
	}
	if err := h.deleteCache(ctx); err != nil {
		return err
	}

	return nil
}

func (h *bindHandler) deleteUser(ctx context.Context) error {
	_, err := usermwcli.DeleteUser(ctx, h.AppID, *h.UserID)
	if err != nil {
		return nil
	}
	if err := h.deleteCache(ctx); err != nil {
		return err
	}

	return nil
}

func (h *Handler) BindUser(ctx context.Context) (*usermwpb.User, error) {
	handler := &bindHandler{
		Handler: h,
	}

	if err := handler.validate(); err != nil {
		return nil, err
	}

	if err := handler.verifyNewAccountCode(ctx); err != nil {
		return nil, err
	}
	if err := handler.verifyNewAccount(ctx); err != nil {
		return nil, err
	}
	if err := handler.getThirdPartyConf(ctx); err != nil {
		return nil, err
	}
	if err := handler.getUser(ctx); err != nil {
		return nil, err
	}
	if err := handler.getThirdUserInfo(ctx); err != nil {
		return nil, err
	}
	if handler.thirdUserInfo.ID != handler.User.ID {
		return nil, fmt.Errorf("invalid userid")
	}

	if err := handler.updateUser(ctx); err != nil {
		return nil, err
	}

	if !h.ShouldUpdateCache {
		return h.User, nil
	}

	if err := handler.updateCache(ctx); err != nil {
		return nil, err
	}

	return h.Metadata.User, nil
}

func (h *Handler) UnbindOAuth(ctx context.Context) error {
	handler := &bindHandler{
		Handler: h,
	}

	if err := handler.validUnbindOAuth(); err != nil {
		return err
	}
	if err := handler.getThirdPartyConf(ctx); err != nil {
		return err
	}
	if err := handler.getUser(ctx); err != nil {
		return err
	}
	if err := handler.getThirdUserInfo(ctx); err != nil {
		return err
	}
	if handler.thirdUserInfo.ID != handler.User.ID {
		return fmt.Errorf("invalid userid")
	}
	if handler.User.EmailAddress != "" || handler.User.PhoneNO != "" {
		err := handler.deleteThirdUserInfo(ctx)
		if err != nil {
			return err
		}
		return nil
	}

	if err := handler.deleteUser(ctx); err != nil {
		return err
	}

	return nil
}
