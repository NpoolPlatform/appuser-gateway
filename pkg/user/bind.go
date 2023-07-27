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
	oldUserInfo *usermwpb.User
}

func (h *bindHandler) checkThirdNewAccount(ctx context.Context) error {
	conds := &usermwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: h.AppID},
	}
	if h.NewAccountType == nil {
		return nil
	}
	switch *h.NewAccountType {
	case basetypes.SignMethod_Email:
		fallthrough //nolint
	case basetypes.SignMethod_Mobile:
		if h.NewAccount == nil {
			return fmt.Errorf("invalid new account")
		}
	default:
		return fmt.Errorf("invalid account type")
	}
	switch *h.NewAccountType {
	case basetypes.SignMethod_Email:
		conds.EmailAddress = &basetypes.StringVal{Op: cruder.EQ, Value: *h.NewAccount}
	case basetypes.SignMethod_Mobile:
		conds.PhoneNO = &basetypes.StringVal{Op: cruder.EQ, Value: *h.NewAccount}
	default:
		return fmt.Errorf("invalid account type")
	}

	info, err := usermwcli.GetUserOnly(ctx, conds)
	if err != nil {
		return err
	}
	if info != nil {
		h.oldUserInfo = info
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
	return nil
}

func (h *bindHandler) getThirdUserInfo(ctx context.Context) (*usermwpb.User, error) {
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
		return nil, err
	}
	if infos == nil {
		return nil, nil
	}
	if len(infos) == 0 {
		return nil, fmt.Errorf("appuserthirdparty is empty")
	}
	if len(infos) > 1 {
		return nil, fmt.Errorf("oauth user too many")
	}

	return infos[0], nil
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
	thirdUserInfo, err := h.getThirdUserInfo(ctx)
	if err != nil {
		return err
	}
	req := &usermwpb.UserReq{
		ID:               h.UserID,
		AppID:            &h.AppID,
		EmailAddress:     h.EmailAddress,
		PhoneNO:          h.PhoneNO,
		ThirdPartyID:     thirdUserInfo.ThirdPartyID,
		ThirdPartyUserID: thirdUserInfo.ThirdPartyUserID,
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

func (h *Handler) BindUser(ctx context.Context) (*usermwpb.User, error) {
	handler := &bindHandler{
		Handler: h,
	}

	if h.UserID == nil {
		return nil, fmt.Errorf("invalid userid")
	}

	if err := handler.checkThirdNewAccount(ctx); err != nil {
		return nil, err
	}
	if err := handler.getUser(ctx); err != nil {
		return nil, err
	}
	if err := handler.verifyNewAccountCode(ctx); err != nil {
		return nil, err
	}
	if err := handler.updateUser(ctx); err != nil {
		return nil, err
	}

	if !h.ShouldUpdateCache {
		return h.User, nil
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
