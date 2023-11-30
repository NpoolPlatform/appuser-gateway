package user

import (
	"context"

	usercodemwcli "github.com/NpoolPlatform/basal-middleware/pkg/client/usercode"
	usercodemwpb "github.com/NpoolPlatform/message/npool/basal/mw/v1/usercode"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) VerifyUserCode(ctx context.Context, usedFor basetypes.UsedFor) error {
	return usercodemwcli.VerifyUserCode(ctx, &usercodemwpb.VerifyUserCodeRequest{
		Prefix:      basetypes.Prefix_PrefixUserCode.String(),
		AppID:       *h.AppID,
		Account:     *h.Account,
		AccountType: *h.AccountType,
		UsedFor:     usedFor,
		Code:        *h.VerificationCode,
	})
}
