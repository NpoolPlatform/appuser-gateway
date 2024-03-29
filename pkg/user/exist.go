package user

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) ExistUser(ctx context.Context) error {
	exist, err := usermwcli.ExistUserConds(ctx, &usermwpb.Conds{
		ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
	})
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid user id=%v, ent_id=%v", *h.ID, *h.EntID)
	}
	return nil
}

func (h *Handler) ExistUserInApp(ctx context.Context) error {
	exist, err := usermwcli.ExistUserConds(ctx, &usermwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID},
	})
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid user app_id=%v, user_id=%v", *h.AppID, *h.UserID)
	}
	return nil
}
