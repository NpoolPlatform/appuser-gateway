package ga

import (
	"context"
	"fmt"

	user1 "github.com/NpoolPlatform/appuser-gateway/pkg/user"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
)

func (h *Handler) SetupGoogleAuth(ctx context.Context) (*usermwpb.User, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(h.AppID),
		user1.WithUserID(&h.UserID),
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

	if user.HasGoogleSecret {
		return user, nil
	}

	secret, err := generateSecret()
	if err != nil {
		return nil, err
	}

	handler.GoogleSecret = &secret
	user, err = handler.UpdateUser(ctx)
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
