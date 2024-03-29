package app

import (
	"context"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"

	"github.com/google/uuid"
)

func (h *Handler) CreateApp(ctx context.Context) (*appmwpb.App, error) {
	id := uuid.NewString()
	if h.EntID == nil {
		h.EntID = &id
	}

	return appmwcli.CreateApp(ctx, &appmwpb.AppReq{
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
	})
}
