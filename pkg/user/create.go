package user

import (
	"context"
	"fmt"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	tracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer/user"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"

	rolemgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approle"
	appusermgrconst "github.com/NpoolPlatform/appuser-manager/pkg/const"
	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	appusermwconst "github.com/NpoolPlatform/appuser-middleware/pkg/const"
	appusermwsvconst "github.com/NpoolPlatform/appuser-middleware/pkg/message/const"
	inspirecli "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/client"
	inspireconst "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/const"
	inspiresvcswconst "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/message/const"
	"github.com/NpoolPlatform/dtm-cluster/pkg/dtm"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"
	rolemgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approle"
	usermwp "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"
	thirdgwconst "github.com/NpoolPlatform/third-gateway/pkg/const"
	"github.com/google/uuid"
)

//nolint:gocyclo
func Signup(ctx context.Context, in *user.SignupRequest) (*usermwp.User, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "Signup")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = tracer.Trace(span, in)
	span = commontracer.TraceInvoker(span, "user", "middleware", "Signup")
	app, err := appmwcli.GetApp(ctx, in.GetAppID())
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, fmt.Errorf("fail get app")
	}

	if app.InvitationCodeMust {
		if in.GetInvitationCode() == "" {
			return nil, fmt.Errorf("invitation code is must")
		}
	}

	inviterID := ""

	if in.GetInvitationCode() != "" {
		span = commontracer.TraceInvoker(span, "user", "manager", "GetUserInvitationCodeOnly")
		code, err := inspirecli.GetUserInvitationCodeOnly(ctx, &inspirepb.Conds{
			InvitationCode: &npool.StringVal{
				Op:    cruder.EQ,
				Value: in.GetInvitationCode(),
			},
		})
		if err != nil {
			logger.Sugar().Errorw("validate", "err", err)
			return nil, err
		}

		if code == nil {
			if app.InvitationCodeMust {
				logger.Sugar().Errorw("validate", "invitation code is must")
				return nil, fmt.Errorf("fail get invitation code")
			}
		} else {
			if code.AppID != in.GetAppID() {
				return nil, fmt.Errorf("invalid invitation code for app")
			}
			inviterID = code.UserID
		}
	}

	span = commontracer.TraceInvoker(span, "user", "verify", "VerifyCode")
	err = VerifyCode(
		ctx,
		in.GetAppID(),
		"",
		in.GetAccount(),
		in.GetAccountType(),
		in.GetVerificationCode(),
		thirdgwconst.UsedForSignup,
		false,
	)
	if err != nil {
		return nil, err
	}

	emailAddress := ""
	phoneNO := ""

	if in.GetAccountType() == appusermgrconst.SignupByMobile {
		phoneNO = in.GetAccount()
	} else if in.GetAccountType() == appusermgrconst.SignupByEmail {
		emailAddress = in.GetAccount()
	}

	userID := uuid.NewString()
	importedFromAppID := uuid.UUID{}.String()

	span = commontracer.TraceInvoker(span, "user", "manager", "GetAppRoleOnly")

	role, err := rolemgrcli.GetAppRoleOnly(ctx, &rolemgrpb.Conds{
		AppID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
		Default: &npool.BoolVal{
			Op:    cruder.EQ,
			Value: true,
		},
	})
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, fmt.Errorf("fail get role")
	}

	var userInfo *usermwp.User

	if in.GetInvitationCode() != "" && inviterID != "" {
		actions := []*dtm.Action{
			{
				ServiceName: appusermwsvconst.ServiceName,
				Action:      appusermwconst.CreateUser,
				Revert:      appusermwconst.CreateUserRevert,
				Param: &usermwp.UserReq{
					ID:                &userID,
					AppID:             &in.AppID,
					EmailAddress:      &emailAddress,
					PhoneNO:           &phoneNO,
					ImportedFromAppID: &importedFromAppID,
					Username:          &in.Username,
					PasswordHash:      &in.PasswordHash,
					RoleID:            &role.ID,
				},
			},
			{
				ServiceName: inspiresvcswconst.ServiceName,
				Action:      inspireconst.CreateRegistrationInvitation,
				Revert:      inspireconst.CreateRegistrationInvitationRevert,
				Param: &inspirepb.CreateRegistrationInvitationRequest{
					Info: &inspirepb.RegistrationInvitation{
						AppID:     in.GetAppID(),
						InviterID: inviterID,
						InviteeID: userID,
					},
				},
			},
		}
		dispose := dtm.SagaDispose{}
		dispose.TransOptions.WaitResult = true
		dispose.TransOptions.TimeoutToFail = 10
		dispose.Actions = actions

		span = commontracer.TraceInvoker(span, "user", "dtm", "WithSaga")

		err = dtm.WithSaga(ctx, &dispose, nil, func(ctx context.Context) error {
			userInfo, err = usermwcli.GetUser(ctx, in.GetAppID(), userID)
			if err != nil {
				return fmt.Errorf("fail create registration invitation: %v", err)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("fail dtm: %v", err)
		}
	} else {

		span = commontracer.TraceInvoker(span, "user", "middleware", "CreateUser")
		userInfo, err = usermwcli.CreateUser(ctx, &usermwp.UserReq{
			ID:                &userID,
			AppID:             &in.AppID,
			EmailAddress:      &emailAddress,
			PhoneNO:           &phoneNO,
			ImportedFromAppID: &importedFromAppID,
			Username:          &in.Username,
			PasswordHash:      &in.PasswordHash,
			RoleID:            &role.ID,
		})
		if err != nil {
			return nil, err
		}
	}

	return userInfo, err
}
