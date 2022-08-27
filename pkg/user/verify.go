package user

import (
	"context"
	"fmt"

	signmethod "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"

	ga "github.com/NpoolPlatform/appuser-gateway/pkg/ga"

	thirdgwpb "github.com/NpoolPlatform/message/npool/thirdgateway"
	thirdgwcli "github.com/NpoolPlatform/third-gateway/pkg/client"
)

func verifyByMobile(ctx context.Context, appID, phoneNO, code, usedFor string) error {
	return thirdgwcli.VerifySMSCode(ctx, &thirdgwpb.VerifySMSCodeRequest{
		AppID:   appID,
		PhoneNO: phoneNO,
		UsedFor: usedFor,
		Code:    code,
	})
}

func verifyByEmail(ctx context.Context, appID, emailAddr, code, usedFor string) error {
	return thirdgwcli.VerifyEmailCode(ctx, &thirdgwpb.VerifyEmailCodeRequest{
		AppID:        appID,
		EmailAddress: emailAddr,
		UsedFor:      usedFor,
		Code:         code,
	})
}

func VerifyCode(
	ctx context.Context,
	appID, userID, account string,
	accountType signmethod.SignMethodType,
	code, usedFor string,
	accountMatch bool,
) error {
	return verifyCode(ctx, appID, userID, account, accountType, code, usedFor, accountMatch)
}

func verifyCode(
	ctx context.Context,
	appID, userID, account string,
	accountType signmethod.SignMethodType,
	code, usedFor string,
	accountMatch bool,
) error {
	var err error

	if accountMatch {
		user, err := usermwcli.GetUser(ctx, appID, userID)
		if err != nil {
			return err
		}

		if user == nil {
			return fmt.Errorf("fail get user ")
		}

		switch accountType {
		case signmethod.SignMethodType_Mobile:
			if user.GetPhoneNO() != account {
				return fmt.Errorf("invalid mobile")
			}
		case signmethod.SignMethodType_Email:
			if user.EmailAddress != account {
				return fmt.Errorf("invalid email")
			}
		}
	}

	switch accountType {
	case signmethod.SignMethodType_Mobile:
		err = verifyByMobile(ctx, appID, account, code, usedFor)
	case signmethod.SignMethodType_Email:
		err = verifyByEmail(ctx, appID, account, code, usedFor)
	default:
		_, err = ga.VerifyGoogleAuth(ctx, appID, userID, code)
	}

	if err != nil {
		return fmt.Errorf("fail verify code: %v", err)
	}

	return nil
}
