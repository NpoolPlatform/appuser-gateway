package user

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	redis2 "github.com/NpoolPlatform/go-service-framework/pkg/redis"
)

type resetLinkHandler struct {
	*Handler
}

type valResetUser struct {
	AppID       string
	UserID      string
	Account     string
	AccountType string
	StartAt     uint32
}

func metaToResetKey(val valResetUser) string {
	key := fmt.Sprintf(
		"%v:%v:%v:%v:%v",
		val.AppID,
		val.UserID,
		val.Account,
		val.AccountType,
		val.StartAt,
	)
	return base64.StdEncoding.EncodeToString([]byte(key))
}

func (h *Handler) CreateResetUserLink(ctx context.Context) (string, error) {
	cli, err := redis2.GetClient()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(ctx, redisTimeout)
	defer cancel()

	val := valResetUser{
		AppID:       *h.AppID,
		UserID:      h.User.EntID,
		Account:     *h.Account,
		AccountType: h.AccountType.String(),
		StartAt:     uint32(time.Now().Unix()),
	}

	body, err := json.Marshal(val)
	if err != nil {
		return "", err
	}

	key := metaToResetKey(val)
	err = cli.Set(ctx, key, body, resetUserExpiration).Err()
	if err != nil {
		return "", err
	}
	return key, nil
}

func (h *Handler) VerifyResetUserLink(ctx context.Context, key string) error {
	cli, err := redis2.GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, redisTimeout)
	defer cancel()

	val, err := cli.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return err
	}

	resetUser := valResetUser{}
	err = json.Unmarshal([]byte(val), &resetUser)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) DeleteResetUserLink(ctx context.Context, link string) error {
	return nil
}
