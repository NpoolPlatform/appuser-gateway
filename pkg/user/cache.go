package user

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redis2 "github.com/NpoolPlatform/go-service-framework/pkg/redis"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

const (
	redisTimeout    = 5 * time.Second
	loginExpiration = 4 * time.Hour
)

func appAccountKey(appID uuid.UUID, account string, accountType basetypes.SignMethod) string {
	return fmt.Sprintf(
		"%v:%v:%v:%v",
		basetypes.Prefix_PrefixUserLogin,
		appID,
		account,
		accountType,
	)
}

func metaToAccountKey(meta *Metadata) string {
	return appAccountKey(
		meta.AppID,
		meta.Account,
		basetypes.SignMethod(basetypes.SignMethod_value[meta.AccountType]),
	)
}

func appUserKey(appID, userID uuid.UUID) string {
	return fmt.Sprintf(
		"%v:%v:%v",
		basetypes.Prefix_PrefixUserLogin,
		appID,
		userID,
	)
}

func metaToUserKey(meta *Metadata) string {
	return appUserKey(meta.AppID, meta.UserID)
}

type valAppUser struct {
	Account     string
	AccountType string
}

func (h *Handler) CreateCache(ctx context.Context) error {
	if h.Metadata == nil {
		return fmt.Errorf("invalid metadata")
	}

	cli, err := redis2.GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, redisTimeout)
	defer cancel()

	meta := h.Metadata

	body, err := json.Marshal(meta)
	if err != nil {
		return err
	}

	err = cli.Set(ctx, metaToAccountKey(meta), body, loginExpiration).Err()
	if err != nil {
		return err
	}

	body, err = json.Marshal(&valAppUser{
		Account:     meta.Account,
		AccountType: meta.AccountType,
	})
	if err != nil {
		return err
	}

	err = cli.Set(ctx, metaToUserKey(meta), body, loginExpiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) QueryCache(ctx context.Context) (*Metadata, error) {
	cli, err := redis2.GetClient()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, redisTimeout)
	defer cancel()

	appID, err := uuid.Parse(*h.AppID)
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(*h.UserID)
	if err != nil {
		return nil, err
	}

	val, err := cli.Get(ctx, appUserKey(appID, userID)).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	appUser := valAppUser{}
	err = json.Unmarshal([]byte(val), &appUser)
	if err != nil {
		return nil, err
	}

	val, err = cli.Get(ctx,
		appAccountKey(
			appID, appUser.Account,
			basetypes.SignMethod(basetypes.SignMethod_value[appUser.AccountType]),
		),
	).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	meta := Metadata{}
	err = json.Unmarshal([]byte(val), &meta)
	if err != nil {
		return nil, err
	}

	return &meta, nil
}

func (h *Handler) DeleteCache(ctx context.Context) error {
	cli, err := redis2.GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, redisTimeout)
	defer cancel()

	err = cli.Del(ctx, metaToUserKey(h.Metadata)).Err()
	if err != nil {
		return err
	}

	err = cli.Del(ctx, metaToAccountKey(h.Metadata)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) UpdateCache(ctx context.Context) error {
	if h.User == nil {
		return fmt.Errorf("invalid user")
	}

	meta, err := h.QueryCache(ctx)
	if err != nil {
		return err
	}
	if meta == nil || meta.User == nil {
		return fmt.Errorf("cache: invalid user: app_id=%v, user_id=%v", *h.AppID, *h.UserID)
	}

	h.User.InvitationCode = meta.User.InvitationCode
	h.User.Logined = meta.User.Logined
	h.User.LoginAccount = meta.User.LoginAccount
	h.User.LoginAccountType = meta.User.LoginAccountType
	h.User.LoginToken = meta.User.LoginToken
	h.User.LoginClientIP = meta.User.LoginClientIP
	h.User.LoginClientUserAgent = meta.User.LoginClientUserAgent
	h.User.LoginVerified = meta.User.LoginVerified

	if h.User.GoogleOTPAuth == "" {
		h.User.GoogleOTPAuth = meta.User.GoogleOTPAuth
	}

	meta.User = h.User
	h.Metadata = meta

	return h.CreateCache(ctx)
}
