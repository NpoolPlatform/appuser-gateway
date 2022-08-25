package authing

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/appuser-gateway/pkg/user"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	authingmgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/authing/auth"
	historymgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/authing/history"
	authingmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/authing"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	authmgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/authing/auth"
	historymgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/authing/history"
	authmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/authing"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
)

func Authenticate(ctx context.Context, appID string, userID, token *string, resource, method string) (resp bool, err error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "Authenticate")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if userID != nil && token == nil {
		return false, fmt.Errorf("invalid token")
	}

	if userID != nil {
		span = commontracer.TraceInvoker(span, "auth", "user", "Logined")

		info, err := user.Logined(ctx, appID, *userID, *token)
		if err != nil {
			return false, err
		}
		if info == nil {
			return false, nil
		}
	}

	span = commontracer.TraceInvoker(span, "auth", "middleware", "ExistAuth")

	userIDStr := ""
	if userID != nil {
		userIDStr = *userID
	}

	allowed, err := authingmwcli.ExistAuth(ctx, appID, userIDStr, resource, method)
	if err != nil {
		return false, err
	}

	go func() {
		span = commontracer.TraceInvoker(span, "auth", "go_manager", "CreateHistory")

		_, err = historymgrcli.CreateHistory(context.Background(), &historymgrpb.HistoryReq{
			AppID:    &appID,
			UserID:   userID,
			Resource: &resource,
			Method:   &method,
			Allowed:  &allowed,
		})
		if err != nil {
			logger.Sugar().Errorw("CreateHistory", "err", err)
		}
	}()

	return allowed, nil
}

func CreateAuth(ctx context.Context, appID string, userID, roleID *string, resource, method string) (resp *authmwpb.Auth, err error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "Authenticate")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	conds := &authmgrpb.Conds{
		AppID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
		Resource: &npool.StringVal{
			Op:    cruder.EQ,
			Value: resource,
		},
		Method: &npool.StringVal{
			Op:    cruder.EQ,
			Value: method,
		},
	}

	if userID != nil {
		conds.UserID = &npool.StringVal{
			Op:    cruder.EQ,
			Value: *userID,
		}
	}

	if roleID != nil {
		conds.RoleID = &npool.StringVal{
			Op:    cruder.EQ,
			Value: *roleID,
		}
	}

	span = commontracer.TraceInvoker(span, "auth", "manager", "GetAuths")

	authInfo, _, err := authingmgrcli.GetAuths(ctx, conds, 0, 1)
	if err != nil {
		return nil, err
	}

	id := ""
	if len(authInfo) > 0 {
		id = authInfo[0].ID

		span = commontracer.TraceInvoker(span, "auth", "manager", "UpdateAuth")

		_, err := authingmgrcli.UpdateAuth(ctx, &authmgrpb.AuthReq{
			ID:       &authInfo[0].ID,
			AppID:    &appID,
			RoleID:   roleID,
			UserID:   userID,
			Resource: &resource,
			Method:   &method,
		})
		if err != nil {
			return nil, err
		}
	} else {
		span = commontracer.TraceInvoker(span, "auth", "manager", "CreateAuth")

		cAuth, err := authingmgrcli.CreateAuth(ctx, &authmgrpb.AuthReq{
			AppID:    &appID,
			RoleID:   roleID,
			UserID:   userID,
			Resource: &resource,
			Method:   &method,
		})
		if err != nil {
			return nil, err
		}
		id = cAuth.ID
	}

	span = commontracer.TraceInvoker(span, "auth", "middleware", "GetAuth")

	return authingmwcli.GetAuth(ctx, id)
}
