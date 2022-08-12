package user

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"

	"github.com/NpoolPlatform/appuser-manager/pkg/client/appuser"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/appuser"
	thirdgwpb "github.com/NpoolPlatform/message/npool/thirdgateway"
	thindcli "github.com/NpoolPlatform/third-gateway/pkg/client"
)

func verifyByMobile(ctx context.Context, appID, phoneNO, code, usedFor string) error {
	return thindcli.VerifySMSCode(ctx, &thirdgwpb.VerifySMSCodeRequest{
		AppID:   appID,
		PhoneNO: phoneNO,
		UsedFor: usedFor,
		Code:    code,
	})
}

func verifyByEmail(ctx context.Context, appID, emailAddr, code, usedFor string) error {
	return thindcli.VerifyEmailCode(ctx, &thirdgwpb.VerifyEmailCodeRequest{
		AppID:        appID,
		EmailAddress: emailAddr,
		UsedFor:      usedFor,
		Code:         code,
	})
}

func verifyByGoogle(ctx context.Context, appID, userID, code string) error {
	return thindcli.VerifyGoogleAuthentication(ctx, &thirdgwpb.VerifyGoogleAuthenticationRequest{
		AppID:  appID,
		UserID: userID,
		Code:   code,
	})
}

func VerifyCode(ctx context.Context, appID, userID, account, accountType, code, usedFor string, accountMatch bool) error {
	var err error

	if accountMatch {
		user, err := appuser.GetAppUserOnly(ctx, &appusermgrpb.Conds{
			ID: &npool.StringVal{
				Op:    cruder.EQ,
				Value: userID,
			},
			AppID: &npool.StringVal{
				Op:    cruder.EQ,
				Value: appID,
			},
		})
		if err != nil {
			return fmt.Errorf("fail get app user: %v", err)
		}

		if user == nil {
			return fmt.Errorf("fail get app user ")
		}

		switch accountType {
		case signmethod.SignMethodType_Mobile.String():
			if user.GetPhoneNo() != account {
				return fmt.Errorf("invalid mobile")
			}
		case signmethod.SignMethodType_Email.String():
			if user.EmailAddress != account {
				return fmt.Errorf("invalid email")
			}
		}
	}

	switch accountType {
	case signmethod.SignMethodType_Mobile.String():
		err = verifyByMobile(ctx, appID, account, code, usedFor)
	case signmethod.SignMethodType_Email.String():
		err = verifyByEmail(ctx, appID, account, code, usedFor)
	default:
		err = verifyByGoogle(ctx, appID, userID, code)
	}

	if err != nil {
		return fmt.Errorf("fail verify code: %v", err)
	}

	return nil
}
