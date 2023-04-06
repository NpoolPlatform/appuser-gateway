package user

import (
	"context"
	"fmt"

	ivcodemwpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/invitationcode"

	ivcodemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/invitationcode"
	commonpb "github.com/NpoolPlatform/message/npool"

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
		InvitationCode: &commonpb.StringVal{Op: cruder.EQ, Value: *h.InvitationCode},
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
