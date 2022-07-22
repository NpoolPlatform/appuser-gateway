package user

import (
	"context"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	"github.com/NpoolPlatform/message/npool/appusergw/appuser"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateUser(ctx context.Context, in *appuser.UpdateUserRequest) (*appuser.UpdateUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateUser")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	err = validate(in.GetInfo())
	if err != nil {
		return &appuser.UpdateUserResponse{}, err
	}

	if _, err := uuid.Parse(in.GetInfo().GetID()); err != nil {
		logger.Sugar().Error("ID is invalid")
		return &appuser.UpdateUserResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgUserIDInvalid)
	}

	span.AddEvent("call grpc CreateAppUserV2")
	resp, err := grpc.UpdateAppUserV2(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("fail create ban app: %v", err)
		return &appuser.UpdateUserResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &appuser.UpdateUserResponse{
		Info: resp,
	}, nil
}
