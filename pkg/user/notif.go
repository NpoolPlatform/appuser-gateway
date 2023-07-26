package user

import (
	"context"
	"time"

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
	logger.Sugar().Infof("generate notifs start")
	if h.UsedFor == basetypes.UsedFor_DefaultUsedFor {
		logger.Sugar().Errorf("no notif situation matched")
		return
	}

	templateVars := &template.TemplateVars{}
	if h.UsedFor == basetypes.UsedFor_NewDeviceDetected {
		clientIP := h.Metadata.ClientIP.String()
		location := h.Metadata.UserAgent
		templateVars.IP = &clientIP
		templateVars.Location = &location
		now := uint32(time.Now().Unix())
		templateVars.Timestamp = &now
	}

	logger.Sugar().Infof(
		"generate notif",
		"AppID", h.AppID,
		"UserID", *h.UserID,
		"EventType", h.UsedFor,
		"UserAgent", h.Metadata.UserAgent,
	)
	_, err := notifmwcli.GenerateNotifs(ctx, &notif.GenerateNotifsRequest{
		AppID:     h.AppID,
		UserID:    *h.UserID,
		EventType: h.UsedFor,
		Vars:      templateVars,
		NotifType: basetypes.NotifType_NotifUnicast,
	})
	if err != nil {
		logger.Sugar().Errorf(
			"send notif error %v", err,
			"AppID", h.AppID,
			"UserID", h.UserID,
			"Vars", templateVars,
		)
	}
}
