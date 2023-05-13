package auth

import (
	"context"
	"fmt"

	user1 "github.com/NpoolPlatform/appuser-gateway/pkg/user"
	authmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/authing/auth"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/go-service-framework/pkg/pubsub"
	authhismwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/authing/history"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

type authenticateHandler struct {
	*Handler
}

func (h *authenticateHandler) logined(ctx context.Context) (bool, error) {
	if h.UserID == nil {
		return true, nil
	}
	if h.Token == nil {
		return false, fmt.Errorf("invalid token")
	}

	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(h.AppID),
		user1.WithUserID(h.UserID),
		user1.WithToken(*h.Token),
	)
	if err != nil {
		return false, err
	}

	info, err := handler.Logined(ctx)
	if err != nil {
		return false, err
	}
	if info == nil {
		return false, fmt.Errorf("user not logined")
	}
	return true, nil
}

func (h *authenticateHandler) notifyAuthenticate(allowed bool) {
	if err := pubsub.WithPublisher(func(publisher *pubsub.Publisher) error {
		req := &authhismwpb.HistoryReq{
			AppID:    &h.AppID,
			UserID:   h.UserID,
			Resource: &h.Resource,
			Method:   &h.Method,
			Allowed:  &allowed,
		}
		return publisher.Update(
			basetypes.MsgID_CreateAuthHistoryReq.String(),
			nil,
			nil,
			nil,
			req,
		)
	}); err != nil {
		logger.Sugar().Errorw(
			"notifyAuthenticate",
			"AppID", h.AppID,
			"UserID", h.UserID,
			"Resource", h.Resource,
			"Method", h.Method,
			"Error", err,
		)
	}
}

func (h *Handler) Authenticate(ctx context.Context) (bool, error) {
	handler := &authenticateHandler{
		Handler: h,
	}

	_allowed := false
	defer handler.notifyAuthenticate(_allowed)

	_logined, err := handler.logined(ctx)
	if err != nil {
		logger.Sugar().Warnw(
			"Authenticate",
			"AppID", h.AppID,
			"UserID", h.UserID,
			"Resource", h.Resource,
			"Method", h.Method,
			"Error", err,
		)
		return false, err
	}
	if !_logined {
		logger.Sugar().Warnw(
			"Authenticate",
			"AppID", h.AppID,
			"UserID", h.UserID,
			"Resource", h.Resource,
			"Method", h.Method,
		)
		return false, nil
	}

	_allowed, err = authmwcli.ExistAuth(ctx, h.AppID, h.UserID, h.Resource, h.Method)
	if err != nil {
		logger.Sugar().Warnw(
			"Authenticate",
			"AppID", h.AppID,
			"UserID", h.UserID,
			"Resource", h.Resource,
			"Method", h.Method,
			"Error", err,
		)
		return false, err
	}
	if !_allowed {
		logger.Sugar().Warnw(
			"Authenticate",
			"AppID", h.AppID,
			"UserID", h.UserID,
			"Resource", h.Resource,
			"Method", h.Method,
			"Allowed", _allowed,
		)
	}

	return _allowed, nil
}
