package authing

import (
	"context"
	"fmt"

	user "github.com/NpoolPlatform/appuser-gateway/pkg/user"
	authingcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/authing"
)

func Authenticate(ctx context.Context, appID string, userID, token *string, resource, method string) (bool, error) {
	if userID != nil && token == nil {
		return false, fmt.Errorf("invalid token")
	}

	if userID != nil {
		info, err := user.Logined(ctx, appID, *userID, *token)
		if err != nil {
			return false, err
		}
		if info == nil {
			return false, nil
		}
	}

	return authingcli.ExistAuth(ctx, appID, userID, resource, method)
}
