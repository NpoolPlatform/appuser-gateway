package user

import (
	"context"
	"fmt"

	recaptcha "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/recaptcha"
	signmethod "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	thirdgwcli "github.com/NpoolPlatform/third-gateway/pkg/client"

	"github.com/google/uuid"
)

func Login(
	ctx context.Context,
	appID, account, passwordHash string,
	accountType signmethod.SignMethodType,
	manMachineSpec, envSpec string,
) (
	*usermwpb.User, error,
) {
	app, err := appmwcli.GetApp(ctx, appID)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, fmt.Errorf("invalid app")
	}

	if app.RecaptchaMethod == recaptcha.RecaptchaType_GoogleRecaptchaV3 {
		err = thirdgwcli.VerifyGoogleRecaptchaV3(ctx, manMachineSpec)
		if err != nil {
			return nil, err
		}
	}

	user, err := usermwcli.VerifyUser(
		ctx,
		appID,
		account,
		accountType,
		passwordHash,
	)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid user")
	}

	meta, err := MetadataFromContext(ctx)
	if err != nil {
		return nil, err
	}

	meta.AppID = uuid.MustParse(appID)
	meta.Account = account
	meta.AccountType = accountType
	meta.UserID = uuid.MustParse(user.ID)

	token, err := createToken(meta)
	if err != nil {
		return nil, err
	}

	user.Logined = true
	user.LoginAccount = account
	user.LoginAccountType = accountType
	user.LoginToken = token
	user.LoginClientIP = meta.ClientIP.String()
	user.LoginClientUserAgent = meta.UserAgent
	meta.User = user

	if err := createCache(ctx, meta); err != nil {
		return nil, err
	}

	// TODO: add login history

	return user, nil
}

func Logined(ctx context.Context, appID, userID, token string) (*usermwpb.User, error) {
	meta, err := queryAppUser(ctx, uuid.MustParse(appID), uuid.MustParse(userID))
	if err != nil {
		return nil, err
	}
	if meta == nil {
		return nil, nil
	}

	if err := verifyToken(meta, token); err != nil {
		return nil, err
	}

	if err := createCache(ctx, meta); err != nil {
		return nil, err
	}

	return meta.User, nil
}

func Logout(ctx context.Context, appID, userID string) (*usermwpb.User, error) {
	meta, err := queryAppUser(ctx, uuid.MustParse(appID), uuid.MustParse(userID))
	if err != nil {
		return nil, err
	}

	if err := deleteCache(ctx, meta); err != nil {
		return nil, err
	}

	return meta.User, nil
}

func UpdateCache(ctx context.Context, user *usermwpb.User) error {
	meta, err := queryAppUser(ctx, uuid.MustParse(user.AppID), uuid.MustParse(user.ID))
	if err != nil {
		return err
	}

	meta.User = user
	if err := createCache(ctx, meta); err != nil {
		return err
	}

	return nil
}
