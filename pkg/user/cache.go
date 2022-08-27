package user

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redis2 "github.com/NpoolPlatform/go-service-framework/pkg/redis"
	signmethod "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

const (
	redisTimeout    = 5 * time.Second
	loginExpiration = 4 * time.Hour
)

func appAccountKey(appID uuid.UUID, account string, accountType signmethod.SignMethodType) string {
	return fmt.Sprintf("login-%v:%v:%v", appID, account, accountType)
}

func metaToAccountKey(meta *Metadata) string {
	return appAccountKey(meta.AppID, meta.Account, signmethod.SignMethodType(signmethod.SignMethodType_value[meta.AccountType]))
}

func appUserKey(appID, userID uuid.UUID) string {
	return fmt.Sprintf("login-%v:%v", appID, userID)
}

func metaToUserKey(meta *Metadata) string {
	return appUserKey(meta.AppID, meta.UserID)
}

type valAppUser struct {
	Account     string
	AccountType string
}

func createCache(ctx context.Context, meta *Metadata) error {
	cli, err := redis2.GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, redisTimeout)
	defer cancel()

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

func QueryAppAccount(ctx context.Context, appID uuid.UUID, account string, accountType signmethod.SignMethodType) (*Metadata, error) {
	cli, err := redis2.GetClient()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, redisTimeout)
	defer cancel()

	val, err := cli.Get(ctx, appAccountKey(appID, account, accountType)).Result()
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

func QueryAppUser(ctx context.Context, appID, userID uuid.UUID) (*usermwpb.User, error) {
	meta, err := queryAppUser(ctx, appID, userID)
	if err != nil {
		return nil, err
	}
	return meta.User, nil
}

func queryAppUser(ctx context.Context, appID, userID uuid.UUID) (*Metadata, error) {
	cli, err := redis2.GetClient()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, redisTimeout)
	defer cancel()

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
			signmethod.SignMethodType(signmethod.SignMethodType_value[appUser.AccountType]),
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

func deleteCache(ctx context.Context, meta *Metadata) error {
	cli, err := redis2.GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, redisTimeout)
	defer cancel()

	err = cli.Del(ctx, metaToUserKey(meta)).Err()
	if err != nil {
		return err
	}

	err = cli.Del(ctx, metaToAccountKey(meta)).Err()
	if err != nil {
		return err
	}

	return nil
}
