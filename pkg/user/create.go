package user

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	tracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer/user"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"

	rolemgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approle"
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
	span = commontracer.TraceInvoker(span, "user", "middleware", "GetApp")

	app, err := appmwcli.GetApp(ctx, in.GetAppID())
	if err != nil {
		logger.Sugar().Errorw("Signup", "err", err)
		return nil, err
	}
	if app == nil {
		logger.Sugar().Errorw("Signup", "err", "fail get app")
		return nil, fmt.Errorf("fail get app")
	}

	if app.InvitationCodeMust {
		if in.GetInvitationCode() == "" {
			logger.Sugar().Errorw("Signup", "err", "invitation code is must")
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
			logger.Sugar().Errorw("Signup", "err", err)
			return nil, err
		}

		if code == nil {
			if app.InvitationCodeMust {
				logger.Sugar().Errorw("Signup", "invitation code is must")
				return nil, fmt.Errorf("fail get invitation code")
			}
		} else {
			if code.AppID != in.GetAppID() {
				logger.Sugar().Errorw("Signup", "err", "invalid invitation code for app")
				return nil, fmt.Errorf("invalid invitation code for app")
			}
			inviterID = code.UserID
		}
	}

	span = commontracer.TraceInvoker(span, "user", "Signup", "VerifyCode")

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
		logger.Sugar().Errorw("Signup", "err", err)
		return nil, err
	}

	emailAddress := ""
	phoneNO := ""

	if in.GetAccountType() == signmethod.SignMethodType_Mobile.String() {
		phoneNO = in.GetAccount()
	} else if in.GetAccountType() == signmethod.SignMethodType_Email.String() {
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
		logger.Sugar().Errorw("Signup", "err", err)
		return nil, err
	}
	if role == nil {
		logger.Sugar().Errorw("Signup", "err", "fail get role")
		return nil, fmt.Errorf("fail get role")
	}

	var userInfo *usermwp.User

	if in.GetInvitationCode() != "" && inviterID != "" {
		actions := []*dtm.Action{
			{
				ServiceName: appusermwsvconst.ServiceName,
				Action:      appusermwconst.CreateUser,
				Revert:      appusermwconst.CreateUserRevert,
				Param: &usermwp.CreateUserRequest{
					Info: &usermwp.UserReq{
						ID:                &userID,
						AppID:             &in.AppID,
						EmailAddress:      &emailAddress,
						PhoneNO:           &phoneNO,
						ImportedFromAppID: &importedFromAppID,
						Username:          &in.Username,
						PasswordHash:      &in.PasswordHash,
						RoleIDs:           []string{role.ID},
					},
				}},
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
		dispose.TransOptions.TimeoutToFail = 2
		dispose.Actions = actions

		span = commontracer.TraceInvoker(span, "user", "dtm", "WithSaga")

		err = dtm.WithSaga(ctx, &dispose, nil, func(ctx context.Context) error {
			userInfo, err = usermwcli.GetUser(ctx, in.GetAppID(), userID)
			if err != nil {
				logger.Sugar().Errorw("Signup", "err", err)
				return err
			}
			return nil
		})
		if err != nil {
			logger.Sugar().Errorw("Signup", "err", err)
			return nil, err
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
			RoleIDs:           []string{role.ID},
		})
		if err != nil {
			logger.Sugar().Errorw("Signup", "err", err)
			return nil, err
		}
	}

	return userInfo, err
}
