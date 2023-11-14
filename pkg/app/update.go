package app

import (
	"context"
	"fmt"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) UpdateApp(ctx context.Context) (*appmwpb.App, error) {
	exist, err := appmwcli.ExistAppConds(ctx, &appmwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
	})
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("invalid app")
	}

	if h.Name != nil {
		exist, err := appmwcli.ExistAppConds(ctx, &appmwpb.Conds{
			Name: &basetypes.StringVal{Op: cruder.EQ, Value: *h.Name},
		})
		if err != nil {
			return nil, err
		}
		if exist {
			return nil, fmt.Errorf("appname exist")
		}
	}

	if err := h.ExistApp(ctx); err != nil {
		return nil, err
	}

	return appmwcli.UpdateApp(ctx, &appmwpb.AppReq{
		ID:                       h.ID,
		EntID:                    h.NewEntID,
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
		CommitButtonTargets:      h.CommitButtonTargets,
		Banned:                   h.Banned,
		BanMessage:               h.BanMessage,
	})
}
