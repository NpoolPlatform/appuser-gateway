package user

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	ivcodemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/invitationcode"
	ivcodemwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/invitation/invitationcode"

	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
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
		ids = append(ids, info.EntID)
	}
	codes, _, err := ivcodemwcli.GetInvitationCodes(ctx, &ivcodemwpb.Conds{
		AppID:   &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		UserIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: ids},
	}, int32(0), int32(len(ids)))
	if err != nil {
		return err
	}
	userCode := map[string]*ivcodemwpb.InvitationCode{}
	for _, code := range codes {
		userCode[code.UserID] = code
	}
	for _, info := range h.infos {
		code, ok := userCode[info.EntID]
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
			AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
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
	if h.UserID == nil {
		return nil, fmt.Errorf("invalid userid")
	}
	info, err := usermwcli.GetUser(ctx, *h.AppID, *h.UserID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("query: invalid user app_id=%v, user_id=%v", *h.AppID, *h.UserID)
	}
	code, _ := ivcodemwcli.GetInvitationCodeOnly(ctx, &ivcodemwpb.Conds{
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		UserID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID},
	})
	if code != nil {
		info.InvitationCode = &code.InvitationCode
	}
	return info, nil
}
