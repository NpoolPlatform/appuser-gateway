package ga

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	user1 "github.com/NpoolPlatform/appuser-gateway/pkg/user"
)

func SetupGoogleAuth(ctx context.Context, appID, userID string) (*usermwpb.User, error) {
	user, err := usermwcli.GetUser(ctx, appID, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid user")
	}

	if user.HasGoogleSecret {
		return user, nil
	}

	secret, err := GenerateSecret()
	if err != nil {
		return nil, err
	}

	user, err = usermwcli.UpdateUser(ctx, &usermwpb.UserReq{
		ID:           &userID,
		AppID:        &appID,
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
	_ = user1.UpdateCache(ctx, user)

	return user, nil
}

func VerifyGoogleAuth(ctx context.Context, appID, userID, code string) (*usermwpb.User, error) {
	user, err := usermwcli.GetUser(ctx, appID, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid user")
	}

	verified, err := VerifyCode(user.GoogleSecret, code)
	if err != nil {
		return nil, err
	}
	if !verified {
		return nil, fmt.Errorf("invalid code")
	}

	user, err = usermwcli.UpdateUser(ctx, &usermwpb.UserReq{
		ID:                 &userID,
		AppID:              &appID,
		GoogleAuthVerified: &verified,
	})
	if err != nil {
		return nil, err
	}

	_ = user1.UpdateCache(ctx, user)
	return user, nil
}
