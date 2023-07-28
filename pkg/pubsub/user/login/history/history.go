package history

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	historymwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user/login/history"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	"github.com/NpoolPlatform/message/npool/notif/mw/v1/notif"
	"github.com/NpoolPlatform/message/npool/notif/mw/v1/template"
	notifmwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif"
)

func Prepare(body string) (interface{}, error) {
	req := historymwpb.HistoryReq{}
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		return nil, err
	}
	return &req, nil
}

func Apply(ctx context.Context, req interface{}) error {
	in, ok := req.(*historymwpb.HistoryReq)
	if !ok {
		return fmt.Errorf("invalid request in apply")
	}

	now := uint32(time.Now().Unix())
	_, err := notifmwcli.GenerateNotifs(ctx, &notif.GenerateNotifsRequest{
		AppID:     *in.AppID,
		UserID:    *in.UserID,
		EventType: basetypes.UsedFor_NewLogin,
		Vars: &template.TemplateVars{
			IP:        in.ClientIP,
			Location:  in.Location,
			UserAgent: in.UserAgent,
			Timestamp: &now,
		},
		NotifType: basetypes.NotifType_NotifUnicast,
	})

	if err != nil {
		logger.Sugar().Errorf(
			"send notif error %v", err,
			"AppID", *in.AppID,
			"UserID", *in.UserID,
			"EventType", basetypes.UsedFor_NewLogin,
			"req", in,
		)
	}

	return err
}
