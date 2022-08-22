//nolint:nolintlint
package migrator

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/NpoolPlatform/appuser-manager/pkg/db"
	"github.com/NpoolPlatform/appuser-manager/pkg/db/ent"
	authingent "github.com/NpoolPlatform/authing-gateway/pkg/db/ent"
	authconstant "github.com/NpoolPlatform/authing-gateway/pkg/message/const"
	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	constant "github.com/NpoolPlatform/go-service-framework/pkg/mysql/const"
)

func Migrate(ctx context.Context) error {
	return migrationAuthingGateway(ctx)
}

const (
	keyUsername = "username"
	keyPassword = "password"
	keyDBName   = "database_name"
	maxOpen     = 10
	maxIdle     = 10
)

func dsn(hostname string) (string, error) {
	username := config.GetStringValueWithNameSpace(constant.MysqlServiceName, keyUsername)
	password := config.GetStringValueWithNameSpace(constant.MysqlServiceName, keyPassword)
	dbname := config.GetStringValueWithNameSpace(hostname, keyDBName)

	svc, err := config.PeekService(constant.MysqlServiceName)
	if err != nil {
		logger.Sugar().Warnw("dsb", "error", err)
		return "", err
	}

	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&interpolateParams=true",
		username, password,
		svc.Address,
		svc.Port,
		dbname,
	), nil
}

func open(hostname string) (conn *sql.DB, err error) {
	hdsn, err := dsn(hostname)
	if err != nil {
		return nil, err
	}

	conn, err = sql.Open("mysql", hdsn)
	if err != nil {
		return nil, err
	}

	// https://github.com/go-sql-driver/mysql
	// See "Important settings" section.

	conn.SetConnMaxLifetime(time.Minute * 3)
	conn.SetMaxOpenConns(maxOpen)
	conn.SetMaxIdleConns(maxIdle)

	return conn, nil
}

func migrationAuthingGateway(ctx context.Context) (err error) {
	cli, err := db.Client()
	if err != nil {
		return err
	}

	authInfos, err := cli.Auth.Query().All(ctx)
	if err != nil {
		return err
	}

	if len(authInfos) > 0 {
		return nil
	}

	logger.Sugar().Infow("Migrate", "Start", "...")
	defer func() {
		logger.Sugar().Infow("Migrate", "Done", "...", "error", err)
	}()

	auth, err := open(authconstant.ServiceName)
	if err != nil {
		return err
	}

	defer auth.Close()

	authCli := authingent.NewClient(authingent.Driver(entsql.OpenDB(dialect.MySQL, auth)))

	appAuths, err := authCli.
		AppAuth.
		Query().
		All(ctx)
	if err != nil {
		return err
	}

	appRoleAuths, err := authCli.
		AppRoleAuth.
		Query().
		All(ctx)
	if err != nil {
		return err
	}

	appUserAuths, err := authCli.
		AppUserAuth.
		Query().
		All(ctx)
	if err != nil {
		return err
	}

	auths := []*ent.Auth{}

	for _, val := range appAuths {
		auths = append(auths, &ent.Auth{
			AppID:     val.AppID,
			Resource:  val.Resource,
			Method:    val.Method,
			CreatedAt: val.CreateAt,
			UpdatedAt: val.UpdateAt,
		})
	}

	for _, val := range appRoleAuths {
		auths = append(auths, &ent.Auth{
			AppID:     val.AppID,
			RoleID:    val.RoleID,
			Resource:  val.Resource,
			Method:    val.Method,
			CreatedAt: val.CreateAt,
			UpdatedAt: val.UpdateAt,
		})
	}

	for _, val := range appUserAuths {
		auths = append(auths, &ent.Auth{
			AppID:     val.AppID,
			UserID:    val.UserID,
			Resource:  val.Resource,
			Method:    val.Method,
			CreatedAt: val.CreateAt,
			UpdatedAt: val.UpdateAt,
		})
	}

	err = db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		bulk := make([]*ent.AuthCreate, len(auths))
		for i, val := range auths {
			bulk[i] = tx.Auth.
				Create().
				SetAppID(val.AppID).
				SetRoleID(val.RoleID).
				SetUserID(val.UserID).
				SetResource(val.Resource).
				SetMethod(val.Method).
				SetCreatedAt(val.CreatedAt).
				SetUpdatedAt(val.UpdatedAt)
		}
		_, err = tx.Auth.CreateBulk(bulk...).Save(_ctx)
		return err
	})
	if err != nil {
		return err
	}

	return nil
}
