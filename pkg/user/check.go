package user

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) CheckUser(ctx context.Context) error {
	conds := &usermwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: h.AppID},
	}
	if h.EmailAddress != nil {
		conds.EmailAddress = &basetypes.StringVal{Op: cruder.EQ, Value: *h.EmailAddress}
	}
	if h.PhoneNO != nil {
		conds.PhoneNO = &basetypes.StringVal{Op: cruder.EQ, Value: *h.PhoneNO}
	}

	exist, err := usermwcli.ExistUserConds(ctx, conds)
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("user already exist")
	}
	return nil
}

func (h *Handler) CheckNewAccount(ctx context.Context) error {
	conds := &usermwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: h.AppID},
	}
	if h.NewAccountType == nil {
		return nil
	}
	switch *h.NewAccountType {
	case basetypes.SignMethod_Email:
		if h.EmailAddress == nil {
			return fmt.Errorf("invalid email address")
		}
		conds.EmailAddress = &basetypes.StringVal{Op: cruder.EQ, Value: *h.EmailAddress}
	case basetypes.SignMethod_Mobile:
		if h.PhoneNO == nil {
			return fmt.Errorf("invalid phone no")
		}
		conds.PhoneNO = &basetypes.StringVal{Op: cruder.EQ, Value: *h.PhoneNO}
	default:
		return fmt.Errorf("invalid account type")
	}

	exist, err := usermwcli.ExistUserConds(ctx, conds)
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("user already exist")
	}
	return nil
}
