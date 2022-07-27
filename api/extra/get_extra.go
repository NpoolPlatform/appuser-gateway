//nolint:nolintlint,dupl
package extra

import (
	"context"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	"github.com/NpoolPlatform/message/npool/appusergw/appuserextra"
	appuserextracrud "github.com/NpoolPlatform/message/npool/appusermgrv2/appuserextra"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetExtra(ctx context.Context,
	in *appuserextra.GetExtraRequest) (*appuserextra.GetExtraResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetExtra")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Error("ID is invalid")
		return &appuserextra.GetExtraResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgIDInvalid)
	}

	span.AddEvent("call grpc GetAppUserExtraV2")
	resp, err := grpc.GetAppUserExtraV2(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Error("fail get extra")
		return &appuserextra.GetExtraResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &appuserextra.GetExtraResponse{
		Info: resp,
	}, nil
}

func (s *Server) GetAppUserExtra(ctx context.Context,
	in *appuserextra.GetAppUserExtraRequest) (*appuserextra.GetAppUserExtraResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppUserExtra")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Error("UserID is invalid")
		return &appuserextra.GetAppUserExtraResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgUserIDInvalid)
	}

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Error("TargetAppID is invalid")
		return &appuserextra.GetAppUserExtraResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgAppIDInvalid)
	}

	span.AddEvent("call grpc GetAppUserExtraOnlyV2")
	resp, err := grpc.GetAppUserExtraOnlyV2(ctx, &appuserextracrud.Conds{
		AppID: &npool.StringVal{
			Value: in.GetTargetAppID(),
			Op:    cruder.EQ,
		},
		UserID: &npool.StringVal{
			Value: in.GetUserID(),
			Op:    cruder.EQ,
		},
	})
	if err != nil {
		logger.Sugar().Error("fail get extra")
		return &appuserextra.GetAppUserExtraResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &appuserextra.GetAppUserExtraResponse{
		Info: resp,
	}, nil
}
