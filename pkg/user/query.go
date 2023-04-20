package user

import (
	"context"
	"fmt"

	servicename "github.com/NpoolPlatform/appuser-gateway/pkg/servicename"
	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	ivcodemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/invitationcode"
	ivcodemgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/invitation/invitationcode"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	scodes "go.opentelemetry.io/otel/codes"
)

func GetUsers(ctx context.Context, appID string, offset, limit int32) ([]*usermwpb.User, uint32, error) {
	var err error

	_, span := otel.Tracer(servicename.ServiceDomain).Start(ctx, "GetUsers")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.SetAttributes(attribute.String("AppID", appID))
	commontracer.TraceOffsetLimit(span, int(offset), int(limit))

	span = commontracer.TraceInvoker(span, "role", "middleware", "CreateUser")

	infos, total, err := usermwcli.GetUsers(ctx, &usermwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: appID},
	}, offset, limit)
	if err != nil {
		logger.Sugar().Errorw("GetUsers", "err", err)
		return nil, 0, err
	}
	if len(infos) == 0 {
		return []*usermwpb.User{}, 0, nil
	}

	userIDs := []string{}
	for _, val := range infos {
		userIDs = append(userIDs, val.ID)
	}

	codes, _, err := ivcodemwcli.GetInvitationCodes(ctx, &ivcodemgrpb.Conds{
		AppID:   &commonpb.StringVal{Op: cruder.EQ, Value: appID},
		UserIDs: &commonpb.StringSliceVal{Op: cruder.IN, Value: userIDs},
	}, int32(0), int32(len(userIDs)))
	if err != nil {
		logger.Sugar().Errorw("GetUsers", "err", err)
		return nil, 0, err
	}

	userCode := map[string]*ivcodemgrpb.InvitationCode{}

	for _, val := range codes {
		userCode[val.UserID] = val
	}

	for key, val := range infos {
		code, ok := userCode[val.ID]
		if ok {
			infos[key].InvitationCode = &code.InvitationCode
		}
	}

	return infos, total, nil
}

func (h *Handler) GetUser(ctx context.Context) (*usermwpb.User, error) {
	info, err := usermwcli.GetUser(ctx, h.AppID, h.UserID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid user %v:%v", h.AppID, h.UserID)
	}

	code, _ := ivcodemwcli.GetInvitationCodeOnly(ctx, &ivcodemgrpb.Conds{
		AppID:  &commonpb.StringVal{Op: cruder.EQ, Value: h.AppID},
		UserID: &commonpb.StringVal{Op: cruder.EQ, Value: h.UserID},
	})

	if code != nil {
		info.InvitationCode = &code.InvitationCode
	}
	return info, nil
}
