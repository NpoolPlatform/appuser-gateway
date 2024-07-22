package ga

import (
	"context"
	"fmt"

	user1 "github.com/NpoolPlatform/appuser-gateway/pkg/user"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/go-service-framework/pkg/pubsub"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	eventmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event"
)

type setupHandler struct {
	*Handler
}

func (h *setupHandler) rewardSet2FA() {
	if err := pubsub.WithPublisher(func(publisher *pubsub.Publisher) error {
		req := &eventmwpb.CalcluateEventRewardsRequest{
			AppID:       *h.AppID,
			UserID:      *h.UserID,
			EventType:   basetypes.UsedFor_Set2FA,
			Consecutive: 1,
		}
		return publisher.Update(
			basetypes.MsgID_CalculateEventRewardReq.String(),
			nil,
			nil,
			nil,
			req,
		)
	}); err != nil {
		logger.Sugar().Errorw(
			"rewardSet2FA",
			"AppID", *h.AppID,
			"UserID", h.UserID,
			"Error", err,
		)
	}
}

func (h *Handler) SetupGoogleAuth(ctx context.Context) (*usermwpb.User, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(h.AppID, true),
		user1.WithUserID(h.UserID, true),
	)
	if err != nil {
		return nil, err
	}

	user, err := handler.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid user")
	}

	if user.HasGoogleSecret {
		return user, nil
	}

	secret, err := generateSecret()
	if err != nil {
		return nil, err
	}

	handler.GoogleSecret = &secret
	user, err = handler.UpdateUser(ctx)
	if err != nil {
		return nil, err
	}

	account := user.EmailAddress
	if account == "" {
		account = user.PhoneNO
	}
	if account == "" {
		return nil, fmt.Errorf("invalid email and mobile")
	}

	user.GoogleOTPAuth = fmt.Sprintf("otpauth://totp/%s?secret=%s", account, user.GoogleSecret)

	handler2 := &setupHandler{
		Handler: h,
	}
	handler2.rewardSet2FA()

	return user, nil
}
