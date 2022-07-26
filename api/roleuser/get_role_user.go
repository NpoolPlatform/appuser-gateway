package roleuser

import (
	"context"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	"github.com/NpoolPlatform/message/npool/appusergw/approleuser"
	approleusercrud "github.com/NpoolPlatform/message/npool/appusermgrv2/approleuser"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetRoleUser(ctx context.Context, in *approleuser.GetRoleUserRequest) (*approleuser.GetRoleUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetRoleUser")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Error("ID is invalid")
		return &approleuser.GetRoleUserResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgIDInvalid)
	}

	span.AddEvent("call grpc GetAppRoleUserV2")
	resp, err := grpc.GetAppRoleUserV2(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Error("fail get role user:%v", err)
		return &approleuser.GetRoleUserResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &approleuser.GetRoleUserResponse{
		Info: resp,
	}, nil
}

func (s *Server) GetRoleUserByUsers(ctx context.Context, in *approleuser.GetRoleUserByUsersRequest) (*approleuser.GetRoleUserByUsersResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetRoleUserByUsers")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Error("AppID is invalid")
		return &approleuser.GetRoleUserByUsersResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgAppIDInvalid)
	}

	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Error("UserID is invalid")
		return &approleuser.GetRoleUserByUsersResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgUserIDInvalid)
	}

	span.AddEvent("call grpc GetAppRoleUsersV2")
	resp, err := grpc.GetAppRoleUsersV2(ctx, &approleusercrud.Conds{
		UserID: &npool.StringVal{
			Value: in.GetUserID(),
			Op:    cruder.EQ,
		},
		AppID: &npool.StringVal{
			Value: in.GetAppID(),
			Op:    cruder.EQ,
		},
	})
	if err != nil {
		logger.Sugar().Error("fail get role user:%v", err)
		return &approleuser.GetRoleUserByUsersResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &approleuser.GetRoleUserByUsersResponse{
		Infos: resp,
	}, nil
}

func (s *Server) GetRoleUsersByRole(ctx context.Context, in *approleuser.GetRoleUsersByRoleRequest) (*approleuser.GetRoleUsersByRoleResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetRoleUsersByRole")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Error("AppID is invalid")
		return &approleuser.GetRoleUsersByRoleResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgAppIDInvalid)
	}

	if _, err := uuid.Parse(in.GetRoleID()); err != nil {
		logger.Sugar().Error("RoleID is invalid")
		return &approleuser.GetRoleUsersByRoleResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgUserIDInvalid)
	}

	span.AddEvent("call grpc GetAppRoleUsersV2")
	resp, err := grpc.GetAppRoleUsersV2(ctx, &approleusercrud.Conds{
		RoleID: &npool.StringVal{
			Value: in.GetRoleID(),
			Op:    cruder.EQ,
		},
		AppID: &npool.StringVal{
			Value: in.GetAppID(),
			Op:    cruder.EQ,
		},
	})
	if err != nil {
		logger.Sugar().Error("fail get role user:%v", err)
		return &approleuser.GetRoleUsersByRoleResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &approleuser.GetRoleUsersByRoleResponse{
		Infos: resp,
	}, nil
}

func (s *Server) GetAppRoleUsersByRole(ctx context.Context, in *approleuser.GetAppRoleUsersByRoleRequest) (*approleuser.GetAppRoleUsersByRoleResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppRoleUsersByRole")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Error("AppID is invalid")
		return &approleuser.GetAppRoleUsersByRoleResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgAppIDInvalid)
	}

	if _, err := uuid.Parse(in.GetRoleID()); err != nil {
		logger.Sugar().Error("RoleID is invalid")
		return &approleuser.GetAppRoleUsersByRoleResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgUserIDInvalid)
	}

	span.AddEvent("call grpc GetAppRoleUsersV2")
	resp, err := grpc.GetAppRoleUsersV2(ctx, &approleusercrud.Conds{
		RoleID: &npool.StringVal{
			Value: in.GetRoleID(),
			Op:    cruder.EQ,
		},
		AppID: &npool.StringVal{
			Value: in.GetTargetAppID(),
			Op:    cruder.EQ,
		},
	})
	if err != nil {
		logger.Sugar().Error("fail get role users:%v", err)
		return &approleuser.GetAppRoleUsersByRoleResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &approleuser.GetAppRoleUsersByRoleResponse{
		Infos: resp,
	}, nil
}

func (s *Server) GetRoleUsers(ctx context.Context, in *approleuser.GetRoleUsersRequest) (*approleuser.GetRoleUsersResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetRoleUsers")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Error("AppID is invalid")
		return &approleuser.GetRoleUsersResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgAppIDInvalid)
	}

	span.AddEvent("call grpc GetAppRoleUsersV2")
	resp, err := grpc.GetAppRoleUsersV2(ctx, &approleusercrud.Conds{
		AppID: &npool.StringVal{
			Value: in.GetAppID(),
			Op:    cruder.EQ,
		},
	})
	if err != nil {
		logger.Sugar().Error("fail get role users:%v", err)
		return &approleuser.GetRoleUsersResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &approleuser.GetRoleUsersResponse{
		Infos: resp,
	}, nil
}

func (s *Server) GetAppRoleUsers(ctx context.Context, in *approleuser.GetAppRoleUsersRequest) (*approleuser.GetAppRoleUsersResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppRoleUsers")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Error("AppID is invalid")
		return &approleuser.GetAppRoleUsersResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgAppIDInvalid)
	}

	span.AddEvent("call grpc GetAppRoleUsersV2")
	resp, err := grpc.GetAppRoleUsersV2(ctx, &approleusercrud.Conds{
		AppID: &npool.StringVal{
			Value: in.GetTargetAppID(),
			Op:    cruder.EQ,
		},
	})
	if err != nil {
		logger.Sugar().Error("fail get role users:%v", err)
		return &approleuser.GetAppRoleUsersResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &approleuser.GetAppRoleUsersResponse{
		Infos: resp,
	}, nil
}
