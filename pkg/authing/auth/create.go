package auth

import (
	"context"
	"fmt"

	authmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/authing/auth"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	authmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/authing/auth"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func (h *Handler) CreateAuth(ctx context.Context) (*authmwpb.Auth, error) {
	conds := &authmwpb.Conds{
		AppID:    &basetypes.StringVal{Op: cruder.EQ, Value: h.AppID},
		Resource: &basetypes.StringVal{Op: cruder.EQ, Value: h.Resource},
		Method:   &basetypes.StringVal{Op: cruder.EQ, Value: h.Method},
	}
	exist, err := authmwcli.ExistAuthConds(ctx, conds)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, fmt.Errorf("auth already exist")
	}

	if h.RoleID != nil {
		conds.RoleID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.RoleID}
	}
	exist, err = authmwcli.ExistAuthConds(ctx, conds)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, fmt.Errorf("auth already exist")
	}

	if h.UserID != nil {
		conds.UserID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID}
	}
	exist, err = authmwcli.ExistAuthConds(ctx, conds)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, fmt.Errorf("auth already exist")
	}

	return authmwcli.CreateAuth(ctx, &authmwpb.AuthReq{
		AppID:    &h.AppID,
		RoleID:   h.RoleID,
		UserID:   h.UserID,
		Resource: &h.Resource,
		Method:   &h.Method,
	})
}
