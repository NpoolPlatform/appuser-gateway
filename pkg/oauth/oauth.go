package oauth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	user1 "github.com/NpoolPlatform/appuser-gateway/pkg/user"
	oauthmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/oauth/appoauththirdparty"
	rolemwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/role"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	redis2 "github.com/NpoolPlatform/go-service-framework/pkg/redis"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	oauthmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/oauth/appoauththirdparty"
	rolemwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/role"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	thirdmwpb "github.com/NpoolPlatform/message/npool/third/mw/v1/oauth"
	thirdmwcli "github.com/NpoolPlatform/third-middleware/pkg/client/oauth"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type ThirdUserInfo struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	Message   string `json:"message"`
}

type oauthHandler struct {
	*Handler
	accessTokenInfo *thirdmwpb.AccessTokenInfo
	oauthConf       *oauthmwpb.OAuthThirdParty
	thirdUserInfo   *thirdmwpb.ThirdUserInfo
	userInfo        *usermwpb.User
}

func (h *Handler) GetOAuthURL(ctx context.Context) (string, error) {
	info, err := oauthmwcli.GetOAuthThirdPartyOnly(ctx, &oauthmwpb.Conds{
		AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		ClientName: &basetypes.Int32Val{Op: cruder.EQ, Value: int32(*h.ClientName)},
	})
	if err != nil {
		return "", err
	}
	if info == nil {
		return "", fmt.Errorf("unsupport oauth")
	}

	clientNameStr := h.ClientName.String()
	state := fmt.Sprintf("%v-%v", clientNameStr, uuid.NewString())
	stateKey := fmt.Sprintf("%v:%v:%v", basetypes.Prefix_PrefixOAuthLogin, *h.AppID, state)
	const expireTime = 5 * time.Minute
	cli, err := redis2.GetClient()
	if err != nil {
		return "", err
	}
	err = cli.Set(ctx, stateKey, clientNameStr, expireTime).Err()
	if err != nil {
		return "", err
	}

	redirectURL := fmt.Sprintf(
		"%s?client_id=%s&scope=%s&redirect_uri=%s&response_type=%s&state=%s",
		info.ClientOAuthURL, info.ClientID, info.Scope, info.CallbackURL, info.ResponseType, state,
	)
	_, err = cli.Get(ctx, stateKey).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return redirectURL, nil
}

func (h *oauthHandler) validate(ctx context.Context) error {
	if h.Code == nil || *h.Code == "" {
		return fmt.Errorf("code is empty")
	}
	if h.State == nil || *h.State == "" {
		return fmt.Errorf("state is empty")
	}
	cli, err := redis2.GetClient()
	if err != nil {
		return err
	}
	stateKey := fmt.Sprintf("%v:%v:%v", basetypes.Prefix_PrefixOAuthLogin, *h.AppID, *h.State)
	clientNameStr, err := cli.Get(ctx, stateKey).Result()
	if err != nil {
		return fmt.Errorf("invalid state")
	}
	clientName := basetypes.SignMethod(basetypes.SignMethod_value[clientNameStr])
	h.ClientName = &clientName

	return nil
}

func (h *oauthHandler) getThirdPartyConf(ctx context.Context) error {
	info, err := oauthmwcli.GetOAuthThirdPartyOnly(
		ctx,
		&oauthmwpb.Conds{
			AppID:         &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			ClientName:    &basetypes.Int32Val{Op: cruder.EQ, Value: int32(*h.ClientName)},
			DecryptSecret: &basetypes.BoolVal{Op: cruder.EQ, Value: true},
		},
	)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("invalid oauth method")
	}
	h.oauthConf = info

	return nil
}

func (h *oauthHandler) getAccessToken(ctx context.Context) error {
	accessTokenInfo, err := thirdmwcli.GetOAuthAccessToken(ctx, *h.ClientName, h.oauthConf.ClientID, h.oauthConf.ClientSecret, *h.Code, h.oauthConf.CallbackURL)
	if err != nil {
		return err
	}
	h.accessTokenInfo = accessTokenInfo
	return nil
}

func (h *oauthHandler) getOAuthUserInfo(ctx context.Context) error {
	thirdUserInfo, err := thirdmwcli.GetOAuthUserInfo(ctx, *h.ClientName, h.accessTokenInfo.AccessToken)
	if err != nil {
		return err
	}

	h.thirdUserInfo = thirdUserInfo
	return nil
}

func (h *oauthHandler) getUserInfo(ctx context.Context) (*usermwpb.User, error) {
	info, err := usermwcli.GetUserOnly(
		ctx,
		&usermwpb.Conds{
			AppID:            &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			ThirdPartyUserID: &basetypes.StringVal{Op: cruder.EQ, Value: h.thirdUserInfo.ID},
			ThirdPartyID:     &basetypes.StringVal{Op: cruder.EQ, Value: h.oauthConf.ThirdPartyID},
		},
	)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}
	h.userInfo = info
	return info, nil
}

func encryptPassword(pwd string) string {
	hash := sha256.New()
	hash.Write([]byte(pwd))
	bytes := hash.Sum(nil)
	hashCode := hex.EncodeToString(bytes)
	return hashCode
}

func (h *oauthHandler) createUserInfo(ctx context.Context) (*usermwpb.User, error) {
	passwordStr := uuid.NewString()
	passwordHash := encryptPassword(passwordStr)
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(h.AppID, true),
		user1.WithPasswordHash(&passwordHash, true),
		user1.WithAccount(&h.thirdUserInfo.ID, true),
		user1.WithAccountType(&h.oauthConf.ClientName, true),
	)

	if err != nil {
		return nil, err
	}

	id := uuid.NewString()
	if handler.UserID == nil {
		handler.UserID = &id
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

	info, err := usermwcli.CreateUser(ctx, &usermwpb.UserReq{
		EntID:              handler.UserID,
		AppID:              h.AppID,
		PasswordHash:       handler.PasswordHash,
		RoleIDs:            []string{role.EntID},
		ThirdPartyID:       &h.oauthConf.ThirdPartyID,
		ThirdPartyUserID:   &h.thirdUserInfo.ID,
		ThirdPartyUsername: &h.thirdUserInfo.Login,
		ThirdPartyAvatar:   &h.thirdUserInfo.AvatarURL,
	})
	if err != nil {
		return nil, err
	}
	h.userInfo = info
	return info, nil
}

func (h *oauthHandler) login(ctx context.Context) (info *usermwpb.User, err error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(h.AppID, true),
		user1.WithUserID(&h.userInfo.EntID, true),
		user1.WithAccount(&h.thirdUserInfo.ID, true),
		user1.WithAccountType(h.ClientName, true),
		user1.WithThirdPartyID(&h.oauthConf.ThirdPartyID, true),
	)
	if err != nil {
		return nil, err
	}
	info, err = handler.ThirdLogin(ctx)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid oauth2")
	}

	return info, nil
}

func (h *Handler) OAuthLogin(ctx context.Context) (info *usermwpb.User, err error) {
	handler := &oauthHandler{
		Handler: h,
	}

	if err := handler.validate(ctx); err != nil {
		return nil, err
	}

	if err := handler.getThirdPartyConf(ctx); err != nil {
		return nil, err
	}

	if err := handler.getAccessToken(ctx); err != nil {
		return nil, err
	}

	if err := handler.getOAuthUserInfo(ctx); err != nil {
		return nil, err
	}

	userInfo, err := handler.getUserInfo(ctx)
	if err != nil {
		return nil, err
	}

	if userInfo == nil {
		_, err = handler.createUserInfo(ctx)
		if err != nil {
			return nil, err
		}
	}

	info, err = handler.login(ctx)
	if err != nil {
		return nil, err
	}

	return info, nil
}
