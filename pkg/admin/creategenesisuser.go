package admin

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

type createGenesisUserHandler struct {
	*Handler
	roles []*rolemwpb.Role
	user  *usermwpb.User
}

func (h *createGenesisUserHandler) getGenesisRoles(ctx context.Context) error {
	const maxGenesisRoles = int32(20)
	infos, _, err := rolemwcli.GetRoles(ctx, &rolemwpb.Conds{
		AppID:   &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		Genesis: &basetypes.BoolVal{Op: cruder.EQ, Value: true},
	}, 0, maxGenesisRoles)
	if err != nil {
		return err
	}
	if len(infos) == 0 {
		return fmt.Errorf("invalid genesis role")
	}
	h.roles = infos
	return nil
}

func (h *createGenesisUserHandler) createUser(ctx context.Context) error {
	userID := uuid.NewString()
	roleIDs := []string{}
	for _, _role := range h.roles {
		roleIDs = append(roleIDs, _role.EntID)
	}

	info, err := usermwcli.CreateUser(ctx, &usermwpb.UserReq{
		EntID:        &userID,
		AppID:        h.AppID,
		EmailAddress: h.EmailAddress,
		PasswordHash: h.PasswordHash,
		RoleIDs:      roleIDs,
	})
	if err != nil {
		return err
	}

	h.user = info

	return nil
}

func (h *Handler) CreateGenesisUser(ctx context.Context) (*usermwpb.User, error) {
	handler := &createGenesisUserHandler{
		Handler: h,
	}
	if err := handler.getGenesisRoles(ctx); err != nil {
		return nil, err
	}
	if err := handler.createUser(ctx); err != nil {
		return nil, err
	}
	return handler.user, nil
}
