package user

import (
	"context"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	inspirecli "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/client"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	scodes "go.opentelemetry.io/otel/codes"
)

func GetUsers(ctx context.Context, appID string, offset, limit int32) ([]*user.User, uint32, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetUsers")
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

	infos, total, err := usermwcli.GetUsers(ctx, appID, offset, limit)
	if err != nil {
		logger.Sugar().Errorw("GetUsers", "err", err)
		return nil, 0, err
	}

	userIDs := []string{}
	for _, val := range infos {
		userIDs = append(userIDs, val.ID)
	}

	codes, err := inspirecli.GetManyUserInvitationCodes(ctx, userIDs)
	if err != nil {
		logger.Sugar().Errorw("GetUsers", "err", err)
		return nil, 0, err
	}

	userCode := map[string]*inspirepb.UserInvitationCode{}

	for _, val := range codes {
		userCode[val.UserID] = val
	}

	for key, val := range infos {
		code, ok := userCode[val.ID]
		if ok {
			infos[key].InvitationCode = &code.InvitationCode
			infos[key].InvitationCodeID = &code.ID
			infos[key].InvitationCodeConfirmed = code.Confirmed
		}
	}

	return infos, total, nil
}
