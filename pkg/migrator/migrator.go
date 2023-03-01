//nolint:nolintlint
package migrator

import (
	"context"

	"github.com/NpoolPlatform/appuser-manager/pkg/db"
	"github.com/NpoolPlatform/appuser-manager/pkg/db/ent"
)

func Migrate(ctx context.Context) error {
	return db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		_, err := tx.
			ExecContext(
				ctx,
				"update app_user_extras set action_credits='0' where action_credits is NULL",
			)
		return err
	})
}
