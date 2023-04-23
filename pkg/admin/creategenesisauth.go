package admin

import (
	"context"
	"encoding/json"
	"fmt"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	authmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/authing/auth"
	roleusermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role/user"
	servicename "github.com/NpoolPlatform/appuser-middleware/pkg/servicename"
	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	authmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/authing/auth"
	roleusermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

type genesisURL struct {
	Path   string
	Method string
}

type createGenesisAuthHandler struct {
	*Handler
	auths     []*authmwpb.Auth
	urls      []*genesisURL
	roleusers []*roleusermwpb.User
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
	for _, _user := range h.roleusers {
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

func (h *createGenesisAuthHandler) getGenesisUsers(ctx context.Context) error {
	appIDs := []string{}
	for _, _app := range h.GenesisApps {
		appIDs = append(appIDs, _app.ID)
	}

	const maxGenesisUsers = int32(100)
	infos, _, err := roleusermwcli.GetUsers(ctx, &roleusermwpb.Conds{
		AppIDs:  &basetypes.StringSliceVal{Op: cruder.EQ, Value: appIDs},
		Genesis: &basetypes.BoolVal{Op: cruder.EQ, Value: true},
	}, 0, maxGenesisUsers)
	if err != nil {
		return err
	}
	if len(infos) == 0 {
		return fmt.Errorf("invalid genesis user")
	}

	h.roleusers = infos
	return nil
}

func (h *Handler) AuthorizeGenesis(ctx context.Context) (infos []*authmwpb.Auth, err error) {
	if err := h.GetGenesisAppConfig(); err != nil {
		return nil, err
	}
	created, err := h.GetGenesisApps(ctx)
	if err != nil {
		return nil, err
	}
	if !created {
		return nil, fmt.Errorf("genesis app not created")
	}

	handler := &createGenesisAuthHandler{
		Handler: h,
	}
	if err := handler.getGenesisUrls(); err != nil {
		return nil, err
	}
	if err := handler.getGenesisUsers(ctx); err != nil {
		return nil, err
	}
	if err := handler.createGenesisAuths(ctx); err != nil {
		return nil, err
	}
	return handler.auths, nil
}
