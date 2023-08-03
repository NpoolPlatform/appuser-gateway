package user

import (
	"context"
	"fmt"

	ivcodemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/invitationcode"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	ivcodemwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/invitation/invitationcode"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
)

func (h *Handler) CheckInvitationCode(ctx context.Context) (*string, error) {
	if h.App.InvitationCodeMust {
		if h.InvitationCode == nil || *h.InvitationCode == "" {
			return nil, fmt.Errorf("invitation code is must")
		}
	}

	if h.InvitationCode == nil || *h.InvitationCode == "" {
		return nil, nil
	}

	ivc, err := ivcodemwcli.GetInvitationCodeOnly(ctx, &ivcodemwpb.Conds{
		InvitationCode: &basetypes.StringVal{Op: cruder.EQ, Value: *h.InvitationCode},
	})
	if err != nil {
		return nil, err
	}
	if ivc == nil {
		return nil, fmt.Errorf("invalid code")
	}

	if ivc.AppID != h.AppID {
		return nil, fmt.Errorf("invalid invitation code")
	}

	return &ivc.UserID, nil
}
