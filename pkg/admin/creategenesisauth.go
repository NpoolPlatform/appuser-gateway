package admin

import (
	"context"
	"encoding/json"
	"fmt"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	authmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/authing/auth"
	servicename "github.com/NpoolPlatform/appuser-middleware/pkg/servicename"
	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	authmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/authing/auth"
)

type genesisURL struct {
	Path   string
	Method string
}

type createGenesisAuthHandler struct {
	*Handler
	auths []*authmwpb.Auth
	urls  []*genesisURL
}

func (h *createGenesisAuthHandler) getGenesisUrls() error {
	str := config.GetStringValueWithNameSpace(
		servicename.ServiceDomain,
		constant.KeyGenesisAuthingAPIs,
	)
	if err := json.Unmarshal([]byte(str), &h.urls); err != nil {
		return err
	}
	if len(h.urls) == 0 {
		return fmt.Errorf("invalid genesis auths")
	}
	return nil
}

func (h *createGenesisAuthHandler) createGenesisAuths(ctx context.Context) error {
	reqs := []*authmwpb.AuthReq{}
	for _, _user := range h.GenesisRoleUsers {
		for _, _url := range h.urls {
			reqs = append(reqs, &authmwpb.AuthReq{
				AppID:    &_user.AppID,
				Resource: &_url.Path,
				Method:   &_url.Method,
				UserID:   &_user.UserID,
				RoleID:   &_user.RoleID,
			})
		}
	}
	auths, err := authmwcli.CreateAuths(ctx, reqs)
	if err != nil {
		return err
	}

	h.auths = auths

	return nil
}

func (h *Handler) AuthorizeGenesis(ctx context.Context) (infos []*authmwpb.Auth, err error) {
	_apps, err := h.GetGenesisApps(ctx)
	if err != nil {
		return nil, err
	}
	if len(_apps) == 0 {
		return nil, fmt.Errorf("genesis app not created")
	}

	handler := &createGenesisAuthHandler{
		Handler: h,
	}
	if err := handler.getGenesisUrls(); err != nil {
		return nil, err
	}

	_roleusers, err := h.GetGenesisRoleUsers(ctx)
	if err != nil {
		return nil, err
	}
	h.GenesisRoleUsers = _roleusers

	_users, err := h.GetGenesisUsers(ctx)
	if err != nil {
		return nil, err
	}
	if len(_users) == 0 {
		return nil, fmt.Errorf("genesis user not created")
	}
	h.GenesisUsers = _users

	if err := handler.createGenesisAuths(ctx); err != nil {
		return nil, err
	}
	return handler.auths, nil
}
