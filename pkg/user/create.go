package user

import (
	"context"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	"github.com/google/uuid"
)

func (h *Handler) CreateUser(ctx context.Context) (*usermwpb.User, error) {
	id := uuid.NewString()
	if h.UserID == nil {
		h.UserID = &id
	}

	return usermwcli.CreateUser(
		ctx,
		&usermwpb.UserReq{
			ID:           h.UserID,
			AppID:        &h.AppID,
			EmailAddress: h.EmailAddress,
			PhoneNO:      h.PhoneNO,
			PasswordHash: h.PasswordHash,
		},
	)
}
