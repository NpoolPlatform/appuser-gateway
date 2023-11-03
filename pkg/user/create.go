package user

import (
	"context"
	"fmt"

	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	rolemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/google/uuid"
)

func (h *Handler) CreateUser(ctx context.Context) (*usermwpb.User, error) {
	if err := h.CheckUser(ctx); err != nil {
		return nil, err
	}

	id := uuid.NewString()
	if h.UserID == nil {
		h.UserID = &id
	}

	role, err := rolemwcli.GetRoleOnly(ctx, &rolemwpb.Conds{
		AppID:   &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		Default: &basetypes.BoolVal{Op: cruder.EQ, Value: true},
	})
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, fmt.Errorf("invalid default role")
	}

	return usermwcli.CreateUser(ctx, &usermwpb.UserReq{
		EntID:        h.UserID,
		AppID:        h.AppID,
		EmailAddress: h.EmailAddress,
		PhoneNO:      h.PhoneNO,
		PasswordHash: h.PasswordHash,
		RoleIDs:      []string{role.EntID},
	})
}
