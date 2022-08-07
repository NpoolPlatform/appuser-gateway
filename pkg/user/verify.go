package user

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"

	"github.com/NpoolPlatform/appuser-manager/pkg/client/appuser"
	appusermgrconst "github.com/NpoolPlatform/appuser-manager/pkg/const"
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
		case appusermgrconst.SignupByMobile:
			if user.GetPhoneNo() != account {
				return fmt.Errorf("invalid mobile")
			}
		case appusermgrconst.SignupByEmail:
			if user.EmailAddress != account {
				return fmt.Errorf("invalid email")
			}
		}
	}

	switch accountType {
	case appusermgrconst.SignupByMobile:
		err = verifyByMobile(ctx, appID, account, code, usedFor)
	case appusermgrconst.SignupByEmail:
		err = verifyByEmail(ctx, appID, account, code, usedFor)
	default:
		err = verifyByGoogle(ctx, appID, userID, code)
	}

	if err != nil {
		return fmt.Errorf("fail verify code: %v", err)
	}

	return nil
}
