package user

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
)

func (h *Handler) DeleteUser(ctx context.Context) (*usermwpb.User, error) {
	if h.UserID == nil {
		return nil, fmt.Errorf("invalid userid")
	}
	return usermwcli.DeleteUser(ctx, h.AppID, *h.UserID)
}
