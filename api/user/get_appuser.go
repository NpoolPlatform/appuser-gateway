package user

import (
	"context"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	appusermw "github.com/NpoolPlatform/appuser-gateway/pkg/middleware/user"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	"github.com/NpoolPlatform/message/npool/appusergw/appuser"
	appusercrud "github.com/NpoolPlatform/message/npool/appusermgrv2/appuser"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetUser(ctx context.Context, in *appuser.GetUserRequest) (*appuser.GetUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppUser")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Error("UserID is invalid")
		return &appuser.GetUserResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgUserIDInvalid)
	}

	span.AddEvent("call grpc GetAppUserV2")
	resp, err := grpc.GetAppUserV2(ctx, in.GetUserID())
	if err != nil {
		return &appuser.GetUserResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &appuser.GetUserResponse{
		Info: resp,
	}, nil
}

func (s *Server) GetUserByAccount(ctx context.Context, in *appuser.GetUserByAccountRequest) (*appuser.GetUserByAccountResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppUser")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Error("ID is invalid")
		return &appuser.GetUserByAccountResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgAppIDInvalid)
	}

	if in.GetAccount() == "" {
		logger.Sugar().Error("Account empty")
		return &appuser.GetUserByAccountResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgAccountEmpty)
	}

	span.AddEvent("call grpc GetAppUserOnlyV2")
	resp, err := appusermw.GetUserByAccount(ctx, in.GetAppID(), in.GetAccount())
	if err != nil {
		return &appuser.GetUserByAccountResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &appuser.GetUserByAccountResponse{
		Info: resp,
	}, nil
}

func (s *Server) GetAppUser(ctx context.Context, in *appuser.GetAppUserRequest) (*appuser.GetAppUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppUser")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Error("ID is invalid")
		return &appuser.GetAppUserResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgUserIDInvalid)
	}

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Error("ID is invalid")
		return &appuser.GetAppUserResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgAppIDInvalid)
	}

	span.AddEvent("call grpc GetAppUserOnlyV2")
	resp, err := grpc.GetAppUserOnlyV2(ctx, &appusercrud.Conds{
		ID: &npool.StringVal{
			Value: in.GetUserID(),
			Op:    cruder.EQ,
		},
		AppID: &npool.StringVal{
			Value: in.GetTargetAppID(),
			Op:    cruder.EQ,
		},
	})
	if err != nil {
		return &appuser.GetAppUserResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &appuser.GetAppUserResponse{
		Info: resp,
	}, nil
}

func (s *Server) GetUserRolesByUser(ctx context.Context, in *appuser.GetUserRolesByUserRequest) (*appuser.GetUserRolesByUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppUser")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Error("AppID is invalid")
		return &appuser.GetUserRolesByUserResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgAppIDInvalid)
	}

	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Error("UserID is invalid")
		return &appuser.GetUserRolesByUserResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgUserIDInvalid)
	}

	span.AddEvent("call grpc GetUserRolesByUser")
	resp, total, err := appusermw.GetUserRolesByUser(ctx, in.GetAppID(), in.GetUserID(), in.GetLimit(), in.GetOffset())
	if err != nil {
		return &appuser.GetUserRolesByUserResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &appuser.GetUserRolesByUserResponse{
		Infos: resp,
		Total: total,
	}, nil
}
