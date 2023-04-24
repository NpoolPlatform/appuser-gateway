package ga

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
)

func (h *Handler) SetupGoogleAuth(ctx context.Context) (*usermwpb.User, error) {
	user, err := usermwcli.GetUser(ctx, h.AppID, h.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid user")
	}

	if user.HasGoogleSecret {
		return user, nil
	}

	secret, err := generateSecret()
	if err != nil {
		return nil, err
	}

	user, err = usermwcli.UpdateUser(ctx, &usermwpb.UserReq{
		ID:           &h.UserID,
		AppID:        &h.AppID,
		GoogleSecret: &secret,
	})
	if err != nil {
		return nil, err
	}

	account := user.EmailAddress
	if account == "" {
		account = user.PhoneNO
	}
	if account == "" {
		return nil, fmt.Errorf("invalid email and mobile")
	}

	user.GoogleOTPAuth = fmt.Sprintf("otpauth://totp/%s?secret=%s", account, user.GoogleSecret)

	return user, nil
}
