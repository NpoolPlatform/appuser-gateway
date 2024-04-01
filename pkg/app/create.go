package app

import (
	"context"

	dtmcli "github.com/NpoolPlatform/dtm-cluster/pkg/dtm"
	inspiremwsvcname "github.com/NpoolPlatform/inspire-middleware/pkg/servicename"
	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
	appconfigmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/app/config"
	"github.com/dtm-labs/dtm/client/dtmcli/dtmimp"

	"github.com/google/uuid"
)

type createHandler struct {
	*Handler
}

func (h *createHandler) withCreateInspireAppConfig(dispose *dtmcli.SagaDispose) {
	id := uuid.NewString()
	req := &appconfigmwpb.AppConfigReq{
		EntID:            &id,
		AppID:            h.EntID,
		CommissionType:   h.CommissionType,
		SettleMode:       h.SettleMode,
		SettleAmountType: h.SettleAmountType,
		SettleInterval:   h.SettleInterval,
		StartAt:          h.StartAt,
		SettleBenefit:    h.SettleBenefit,
	}
	dispose.Add(
		inspiremwsvcname.ServiceDomain,
		"inspire.middleware.app.config.v1.Middleware/CreateAppConfig",
		"inspire.middleware.app.config.v1.Middleware/DeleteAppConfig",
		&appconfigmwpb.CreateAppConfigRequest{
			Info: req,
		},
	)
}

func (h *createHandler) withCreateApp(dispose *dtmcli.SagaDispose) {
	req := &appmwpb.AppReq{
		EntID:                    h.EntID,
		CreatedBy:                h.CreatedBy,
		Name:                     h.Name,
		Logo:                     h.Logo,
		Description:              h.Description,
		SignupMethods:            h.SignupMethods,
		ExtSigninMethods:         h.ExtSigninMethods,
		RecaptchaMethod:          h.RecaptchaMethod,
		KycEnable:                h.KycEnable,
		SigninVerifyEnable:       h.SigninVerifyEnable,
		InvitationCodeMust:       h.InvitationCodeMust,
		CreateInvitationCodeWhen: h.CreateInvitationCodeWhen,
		MaxTypedCouponsPerOrder:  h.MaxTypedCouponsPerOrder,
		Maintaining:              h.Maintaining,
		CouponWithdrawEnable:     h.CouponWithdrawEnable,
		CommitButtonTargets:      h.CommitButtonTargets,
		ResetUserMethod:          h.ResetUserMethod,
	}
	dispose.Add(
		inspiremwsvcname.ServiceDomain,
		"appuser.middleware.app.v1.Middleware/CreateApp",
		"appuser.middleware.app.v1.Middleware/DeleteApp",
		&appmwpb.CreateAppRequest{
			Info: req,
		},
	)
}

func (h *Handler) CreateApp(ctx context.Context) (*appmwpb.App, error) {
	id := uuid.NewString()
	if h.EntID == nil {
		h.EntID = &id
	}
	handler := &createHandler{
		Handler: h,
	}

	sagaDispose := dtmcli.NewSagaDispose(dtmimp.TransOptions{
		WaitResult:     true,
		RequestTimeout: *handler.RequestTimeoutSeconds,
	})

	handler.withCreateApp(sagaDispose)
	handler.withCreateInspireAppConfig(sagaDispose)

	if err := dtmcli.WithSaga(ctx, sagaDispose); err != nil {
		return nil, err
	}

	return handler.GetApp(ctx)
}
