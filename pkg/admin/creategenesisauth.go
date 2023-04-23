package admin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	tracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer/admin"
	appmgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/app"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	appmrgpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/app"
	rolemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role"

	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/admin"
	"github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/const"
	servicename "github.com/NpoolPlatform/appuser-gateway/pkg/servicename"
	servicename2 "github.com/NpoolPlatform/appuser-manager/pkg/servicename"

	approlemgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approle"
	roleusermgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/approleuser"
	authmgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/authing/auth"
	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	authmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/authing"
	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role"
	roleusermgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approleuser"

	approlepb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/approle"
	authmgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/authing/auth"
	appmw "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
	authmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/authing/auth"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
)

func createGenesisAuths(ctx context.Context, appID string) (infos []*authmwpb.Auth, total uint32, err error) {
	roleUsers, _, err := roleusermgrcli.GetAppRoleUsers(ctx, &roleusermgrpb.Conds{
		AppID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
	}, 0, 0)

	auths := []*authmgrpb.AuthReq{}

	for _, val := range roleUsers {
		for _, _api := range apis {
			api := _api
			auths = append(auths, &authmgrpb.AuthReq{
				AppID:    &val.AppID,
				Resource: &api.Path,
				Method:   &api.Method,
				UserID:   &val.UserID,
				RoleID:   &val.RoleID,
			})
		}
	}

	_, err = authmgrcli.CreateAuths(ctx, auths)
	if err != nil {
		return nil, 0, err
	}

	return authmwcli.GetAuths(ctx, appID, 0, 0)
}

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
		servicename2.ServiceDomain,
		constant.KeyGenesisAuthingAPIs,
	)
	if err := json.Unmarshal([]byte(apisJSON), &h.urls); err != nil {
		return err
	}
	if len(h.urls) == 0 {
		return fmt.Errorf("invalid genesis auths")
	}
	return nil
}

func (h *createGenesisAuthHandler) createGenesisAuths(ctx context.Context) error {
	for _, _app := range h.GenesisApps {
		authInfos, _, err := createGenesisAuths(ctx, val.GetID())
		if err != nil {
			return err
		}
		infos = append(infos, authInfos...)
	}
	return nil
}

func (h *Handler) AuthorizeGenesis(ctx context.Context) (infos []*authmwpb.Auth, err error) {
	if err := h.GetGenesisAppConfig(); err != nil {
		return nil, err
	}
	if err := h.GetGenesisApps(ctx); err != nil {
		return err
	}

	handler := &createGenesisAuthHandler{
		Handler: h,
	}
	if err := handler.getGenesisUrls(); err != nil {
		return nil, err
	}
	if err := handler.createGenesisAuths(ctx); err != nil {
		return nil, err
	}
	return h.auths, nil
}
