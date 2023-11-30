package ga

import (
	"context"
	"fmt"

	user1 "github.com/NpoolPlatform/appuser-gateway/pkg/user"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	usercodemwcli "github.com/NpoolPlatform/basal-middleware/pkg/client/usercode"
	usercodemwpb "github.com/NpoolPlatform/message/npool/basal/mw/v1/usercode"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) VerifyGoogleAuth(ctx context.Context) (*usermwpb.User, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(h.AppID, true),
		user1.WithUserID(h.UserID, true),
	)
	if err != nil {
		return nil, err
	}

	user, err := handler.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid user")
	}
	if !user.HasGoogleSecret {
		return nil, fmt.Errorf("invalid google secret")
	}

	if err := usercodemwcli.VerifyUserCode(ctx, &usercodemwpb.VerifyUserCodeRequest{
		Prefix:      basetypes.Prefix_PrefixUserCode.String(),
		AppID:       *h.AppID,
		Account:     user.GoogleSecret,
		AccountType: basetypes.SignMethod_Google,
		UsedFor:     basetypes.UsedFor_Update,
		Code:        *h.Code,
	}); err != nil {
		return nil, err
	}

	verified := true
	handler.GoogleAuthVerified = &verified
	user, err = handler.UpdateUser(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}
