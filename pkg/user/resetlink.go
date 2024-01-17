package user

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	appusertypes "github.com/NpoolPlatform/message/npool/basetypes/appuser/v1"

	redis2 "github.com/NpoolPlatform/go-service-framework/pkg/redis"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	"github.com/go-redis/redis/v8"
)

type valResetUser struct {
	AppID       string
	UserID      string
	Account     string
	AccountType string
	StartAt     uint32
}

func resetUserKey(val valResetUser) string {
	return fmt.Sprintf(
		"%v:%v:%v:%v:%v:%v",
		basetypes.Prefix_PrefixCreateResetUserLink,
		val.AppID,
		val.UserID,
		val.Account,
		val.AccountType,
		val.StartAt,
	)
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

	key := resetUserKey(val)
	err = cli.Set(ctx, key, body, resetUserExpiration).Err()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString([]byte(key)), nil
}

func (h *Handler) VerifyResetUserLink(ctx context.Context) error {
	if h.App.ResetUserMethod != appusertypes.ResetUserMethod_Link {
		return nil
	}
	if h.ResetToken == nil {
		return fmt.Errorf("invalid reset token")
	}
	cli, err := redis2.GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, redisTimeout)
	defer cancel()

	tokenBytes, err := base64.StdEncoding.DecodeString(*h.ResetToken)
	if err != nil {
		return err
	}
	key := string(tokenBytes[:]) //nolint

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
	if resetUser.UserID != *h.UserID ||
		resetUser.Account != *h.Account ||
		resetUser.AccountType != h.AccountType.String() ||
		resetUser.AppID != *h.AppID {
		return fmt.Errorf("reset token invalid")
	}
	return nil
}

func (h *Handler) DeleteResetUserLink(ctx context.Context) error {
	if h.ResetToken == nil {
		return nil
	}
	tokenBytes, err := base64.StdEncoding.DecodeString(*h.ResetToken)
	if err != nil {
		return err
	}
	key := string(tokenBytes[:]) //nolint

	cli, err := redis2.GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, redisTimeout)
	defer cancel()

	err = cli.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
