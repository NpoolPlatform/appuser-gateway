package user

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	ivcodemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/invitationcode"
	ivcodemgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/invitationcode"

	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

type queryHandler struct {
	*Handler
	infos []*usermwpb.User
	total uint32
}

func (h *queryHandler) getInvitationCodes(ctx context.Context) error {
	ids := []string{}
	for _, info := range h.infos {
		ids = append(ids, info.ID)
	}
	codes, _, err := ivcodemwcli.GetInvitationCodes(
		ctx,
		&ivcodemgrpb.Conds{
			AppID:   &commonpb.StringVal{Op: cruder.EQ, Value: h.AppID},
			UserIDs: &commonpb.StringSliceVal{Op: cruder.IN, Value: ids},
		},
		int32(0),
		int32(len(ids)),
	)
	if err != nil {
		return err
	}
	userCode := map[string]*ivcodemgrpb.InvitationCode{}
	for _, code := range codes {
		userCode[code.UserID] = code
	}
	for _, info := range h.infos {
		code, ok := userCode[info.ID]
		if ok {
			info.InvitationCode = &code.InvitationCode
		}
	}
	return nil
}

func (h *Handler) GetUsers(ctx context.Context) ([]*usermwpb.User, uint32, error) {
	handler := &queryHandler{
		Handler: h,
	}
	infos, total, err := usermwcli.GetUsers(
		ctx,
		&usermwpb.Conds{
			AppID: &basetypes.StringVal{Op: cruder.EQ, Value: h.AppID},
		},
		h.Offset,
		h.Limit,
	)
	if err != nil {
		return nil, 0, err
	}
	handler.total = total
	handler.infos = infos
	if err := handler.getInvitationCodes(ctx); err != nil {
		return nil, 0, err
	}
	return handler.infos, handler.total, nil
}

func (h *Handler) GetUser(ctx context.Context) (*usermwpb.User, error) {
	info, err := usermwcli.GetUser(ctx, h.AppID, h.UserID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid user %v:%v", h.AppID, h.UserID)
	}
	code, _ := ivcodemwcli.GetInvitationCodeOnly(ctx, &ivcodemgrpb.Conds{
		AppID:  &commonpb.StringVal{Op: cruder.EQ, Value: h.AppID},
		UserID: &commonpb.StringVal{Op: cruder.EQ, Value: h.UserID},
	})
	if code != nil {
		info.InvitationCode = &code.InvitationCode
	}
	return info, nil
}
