package user

import (
	"context"
	"fmt"
	"time"

	loginhispb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/login/history"
	recaptcha "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/recaptcha"
	signmethod "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	loginhiscli "github.com/NpoolPlatform/appuser-manager/pkg/client/login/history"
	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	thirdgwcli "github.com/NpoolPlatform/third-gateway/pkg/client"

	thirdgwconst "github.com/NpoolPlatform/third-gateway/pkg/const"

	commonpb "github.com/NpoolPlatform/message/npool"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	"github.com/go-resty/resty/v2"

	"github.com/google/uuid"
)

func addHistory(appID, userID, clientIP, userAgent string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //nolint
	defer cancel()

	location := ""
	histories, _, err := loginhiscli.GetHistories(ctx, &loginhispb.Conds{
		ClientIP: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: clientIP,
		},
		Location: &commonpb.StringVal{
			Op:    cruder.NEQ,
			Value: "",
		},
	}, 0, 1)
	if err != nil {
		return
	}
	if len(histories) > 0 {
		location = histories[0].Location
	} else {
		type resp struct {
			Error   bool   `json:"error"`
			City    string `json:"city"`
			Country string `json:"country_name"`
			IP      string `json:"ip"`
			Reason  string `json:"reason"`
		}

		r, err := resty.
			New().
			R().
			SetResult(&resp{}).
			Get(fmt.Sprintf("https://ipapi.co/%v/json", clientIP))
		if err == nil {
			rc, ok := r.Result().(*resp)
			if ok && !rc.Error {
				location = fmt.Sprintf("%v, %v", rc.City, rc.Country)
			}
		}
	}

	_, _ = loginhiscli.CreateHistory(
		ctx,
		&loginhispb.HistoryReq{
			AppID:     &appID,
			UserID:    &userID,
			ClientIP:  &clientIP,
			UserAgent: &userAgent,
			Location:  &location,
		},
	)
}

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

	go addHistory(appID, user.ID, meta.ClientIP.String(), meta.UserAgent)

	return user, nil
}

func LoginVerify(ctx context.Context, appID, userID, token, code string) (*usermwpb.User, error) {
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

	account := meta.User.EmailAddress
	accountType := meta.User.SigninVerifyType

	switch meta.User.SigninVerifyType {
	case signmethod.SignMethodType_Email:
	case signmethod.SignMethodType_Mobile:
		account = meta.User.PhoneNO
	case signmethod.SignMethodType_Google:
	default:
		return nil, fmt.Errorf("not supported")
	}
	if err := verifyCode(
		ctx,
		appID, userID,
		account, accountType,
		code,
		thirdgwconst.UsedForSignin,
		true,
	); err != nil {
		return nil, err
	}

	meta.User.LoginVerified = true
	if err := createCache(ctx, meta); err != nil {
		return nil, err
	}

	return meta.User, nil
}

func Logined(ctx context.Context, appID, userID, token string) (*usermwpb.User, error) {
	meta, err := queryAppUser(ctx, uuid.MustParse(appID), uuid.MustParse(userID))
	if err != nil {
		return nil, err
	}
	if meta == nil {
		return nil, nil
	}
	if !meta.User.LoginVerified {
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