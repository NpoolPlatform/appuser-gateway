package user

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	"github.com/NpoolPlatform/message/npool/notif/mw/v1/notif"
	"github.com/NpoolPlatform/message/npool/notif/mw/v1/template"
	notifmwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif"
)

type notifHandler struct {
	*Handler
	UsedFor basetypes.UsedFor
}

func (h *notifHandler) GetUsedFor() {
	h.UsedFor = basetypes.UsedFor_DefaultUsedFor

	if h.NewAccountType != nil {
		switch *h.NewAccountType {
		case basetypes.SignMethod_Email:
			h.UsedFor = basetypes.UsedFor_UpdateEmail
		case basetypes.SignMethod_Mobile:
			h.UsedFor = basetypes.UsedFor_UpdateMobile
		}
	}

	if h.OldPasswordHash != nil {
		h.UsedFor = basetypes.UsedFor_UpdatePassword
	}

	if h.NewAccountType != nil && h.NewVerificationCode != nil {
		if *h.NewAccountType == basetypes.SignMethod_Google {
			h.UsedFor = basetypes.UsedFor_UpdateGoogleAuth
		}
	}
}

func (h *notifHandler) GenerateNotif(ctx context.Context) {
	if h.UsedFor == basetypes.UsedFor_DefaultUsedFor {
		logger.Sugar().Errorf("no notif situation matched")
		return
	}

	_, err := notifmwcli.GenerateNotifs(ctx, &notif.GenerateNotifsRequest{
		AppID:     h.AppID,
		UserID:    *h.UserID,
		EventType: h.UsedFor,
		Vars:      &template.TemplateVars{},
		NotifType: basetypes.NotifType_NotifUnicast,
	})
	if err != nil {
		logger.Sugar().Errorf("send notif error %v", err)
	}
}
