package user

import (
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
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

func (h *notifHandler) GenerateNotif() {
	if h.UsedFor == basetypes.UsedFor_DefaultUsedFor {
		return
	}
}
