package user

import (
	"context"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	inspirecli "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	scodes "go.opentelemetry.io/otel/codes"
)

func GetUsers(ctx context.Context, appID string, offset, limit int32) ([]*user.User, error) {
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

	resp, _, err := usermwcli.GetUsers(ctx, appID, offset, limit)
	if err != nil {
		logger.Sugar().Errorw("GetUsers", "err", err)
		return nil, err
	}

	for key, val := range resp {
		code, err := inspirecli.GetUserInvitationCodeByAppUser(ctx, appID, val.ID)
		if err != nil {
			logger.Sugar().Errorw("GetUsers", "err", err)
			return nil, err
		}
		resp[key].InvitationCode = code.GetInvitationCode()
	}
	return resp, nil
}
