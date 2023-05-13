package user

import (
	"context"
	"fmt"

	usercodemwcli "github.com/NpoolPlatform/basal-middleware/pkg/client/usercode"
	usercodemwpb "github.com/NpoolPlatform/message/npool/basal/mw/v1/usercode"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) VerifyUserCode(ctx context.Context, usedFor basetypes.UsedFor) error {
	if h.Account == nil || h.AccountType == nil {
		return fmt.Errorf("invalid account type")
	}
	if h.VerificationCode == nil {
		return fmt.Errorf("invalid verification code")
	}
	return usercodemwcli.VerifyUserCode(
		ctx,
		&usercodemwpb.VerifyUserCodeRequest{
			Prefix:      basetypes.Prefix_PrefixUserCode.String(),
			AppID:       h.AppID,
			Account:     *h.Account,
			AccountType: *h.AccountType,
			UsedFor:     usedFor,
			Code:        *h.VerificationCode,
		},
	)
}
