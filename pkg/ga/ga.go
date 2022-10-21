package ga

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	"github.com/NpoolPlatform/message/npool/third/mgr/v1/usedfor"
	thirdmwcli "github.com/NpoolPlatform/third-middleware/pkg/client/verify"
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

	secret, err := generateSecret()
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

	if err := thirdmwcli.VerifyCode(
		ctx,
		appID,
		user.GoogleSecret,
		code,
		signmethod.SignMethodType_Google,
		usedfor.UsedFor_Update,
	); err != nil {
		return nil, err
	}

	verified := true

	user, err = usermwcli.UpdateUser(ctx, &usermwpb.UserReq{
		ID:                 &userID,
		AppID:              &appID,
		GoogleAuthVerified: &verified,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
