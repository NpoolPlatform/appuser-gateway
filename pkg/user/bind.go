package user

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	usercodemwcli "github.com/NpoolPlatform/basal-middleware/pkg/client/usercode"
	usercodemwpb "github.com/NpoolPlatform/message/npool/basal/mw/v1/usercode"
)

type bindHandler struct {
	*Handler
	oldUserInfo   *usermwpb.User
	thirdUserInfo *usermwpb.User
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

func (h *bindHandler) getUser(ctx context.Context) error {
	info, err := usermwcli.GetUser(ctx, h.AppID, *h.UserID)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("bind: invalid user: app_id=%v, user_id=%v", h.AppID, *h.UserID)
	}
	if info.EmailAddress != "" || info.PhoneNO != "" {
		return fmt.Errorf("bind: invalid user: account has been bound")
	}

	h.User = info
	h.oldUserInfo = info
	return nil
}

func (h *bindHandler) getThirdUserInfo(ctx context.Context) error {
	const maxlimit = 2
	infos, _, err := usermwcli.GetThirdUsers(
		ctx,
		&usermwpb.Conds{
			ThirdPartyUserID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.Account},
		},
		0,
		maxlimit,
	)
	if err != nil {
		return err
	}
	if infos == nil {
		return fmt.Errorf("appuserthirdparty is empty")
	}
	if len(infos) == 0 {
		return fmt.Errorf("appuserthirdparty is empty")
	}
	if len(infos) > 1 {
		return fmt.Errorf("oauth user too many")
	}
	h.thirdUserInfo = infos[0]

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

func (h *bindHandler) updateUser(ctx context.Context) error {
	req := &usermwpb.UserReq{
		ID:               h.UserID,
		AppID:            &h.AppID,
		EmailAddress:     h.EmailAddress,
		PhoneNO:          h.PhoneNO,
		ThirdPartyID:     h.thirdUserInfo.ThirdPartyID,
		ThirdPartyUserID: h.thirdUserInfo.ThirdPartyUserID,
	}
	fmt.Printf("new_account_type=%v, new_account=%v\n", h.NewAccountType, h.NewAccount)
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

	info, err := usermwcli.UpdateUser(ctx, req)
	if err != nil {
		return err
	}

	h.User = info
	return nil
}

func (h *bindHandler) updateCache(ctx context.Context) error {
	if h.oldUserInfo.ID != h.User.ID {
		meta, err := h.QueryCache(ctx)
		if err != nil {
			return err
		}
		if meta == nil || meta.User == nil {
			return fmt.Errorf("cache: invalid user: app_id=%v, user_id=%v", h.AppID, *h.UserID)
		}
		h.Metadata = meta
		h.UserID = &h.oldUserInfo.ID
		if err := h.DeleteCache(ctx); err != nil {
			return err
		}

		h.UserID = &h.User.ID
		handler := &loginHandler{
			Handler: h.Handler,
		}
		if err := handler.prepareMetadata(ctx); err != nil {
			return err
		}
		token, err := createToken(h.Metadata)
		if err != nil {
			return err
		}
		h.Token = &token
		handler.formalizeUser()
		if err := h.CreateCache(ctx); err != nil {
			return err
		}
	} else {
		if err := h.UpdateCache(ctx); err != nil {
			return err
		}
	}
	meta, err := h.QueryCache(ctx)
	if err != nil {
		return err
	}
	h.Metadata = meta

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
