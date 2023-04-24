package ga

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	usercodemwcli "github.com/NpoolPlatform/basal-middleware/pkg/client/usercode"
	usercodemwpb "github.com/NpoolPlatform/message/npool/basal/mw/v1/usercode"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) VerifyGoogleAuth(ctx context.Context) (*usermwpb.User, error) {
	user, err := usermwcli.GetUser(ctx, h.AppID, h.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid user")
	}

	if err := usercodemwcli.VerifyUserCode(
		ctx,
		&usercodemwpb.VerifyUserCodeRequest{
			Prefix:      basetypes.Prefix_PrefixUserCode.String(),
			AppID:       h.AppID,
			Account:     user.GoogleSecret,
			AccountType: basetypes.SignMethod_Google,
			UsedFor:     basetypes.UsedFor_Update,
			Code:        h.Code,
		},
	); err != nil {
		return nil, err
	}

	verified := true

	user, err = usermwcli.UpdateUser(ctx, &usermwpb.UserReq{
		ID:                 &h.UserID,
		AppID:              &h.AppID,
		GoogleAuthVerified: &verified,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
