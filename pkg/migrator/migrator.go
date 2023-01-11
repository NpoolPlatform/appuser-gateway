//nolint:nolintlint
package migrator

import (
	"context"

	"github.com/NpoolPlatform/appuser-manager/pkg/db"
	"github.com/NpoolPlatform/appuser-manager/pkg/db/ent"

	appctrlmgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/appcontrol"
)

func Migrate(ctx context.Context) error {
	if err := db.Init(); err != nil {
		return err
	}

	return db.WithClient(ctx, func(_ctx context.Context, cli *ent.Client) error {
		_, err := cli.
			AppControl.
			Update().
			SetMaxTypedCouponsPerOrder(1).
			SetCreateInvitationCodeWhen(appctrlmgrpb.CreateInvitationCodeWhen_Registration.String()).
			Save(_ctx)
		return err
	})
}
