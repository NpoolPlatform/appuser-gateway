package roleuser

import (
	"context"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	mw "github.com/NpoolPlatform/appuser-gateway/pkg/middleware/roleuser"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw/approleuser"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateRoleUser(ctx context.Context, in *approleuser.CreateRoleUserRequest) (*approleuser.CreateRoleUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateRoleUser")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	err = validate(in.GetInfo())
	if err != nil {
		return nil, err
	}

	span.AddEvent("call middleware CreateRoleUser")
	resp, err := mw.CreateRoleUser(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("fail create app role user: %v", err)
		return &approleuser.CreateRoleUserResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &approleuser.CreateRoleUserResponse{
		Info: resp,
	}, nil
}

func (s *Server) CreateAppUserRoleUser(ctx context.Context, in *approleuser.CreateAppUserRoleUserRequest) (*approleuser.CreateAppUserRoleUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateRoleUser")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	err = validate(in.GetInfo())
	if err != nil {
		return nil, err
	}

	info := in.GetInfo()
	appID := in.GetTargetAppID()
	userID := in.GetTargetUserID()
	info.UserID = &appID
	info.AppID = &userID

	span.AddEvent("call middleware CreateRoleUser")
	resp, err := mw.CreateRoleUser(ctx, info)
	if err != nil {
		logger.Sugar().Errorw("fail create app role user: %v", err)
		return &approleuser.CreateAppUserRoleUserResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &approleuser.CreateAppUserRoleUserResponse{
		Info: resp,
	}, nil
}

func (s *Server) CreateUserRoleUser(ctx context.Context, in *approleuser.CreateUserRoleUserRequest) (*approleuser.CreateUserRoleUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateUserRoleUser")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	err = validate(in.GetInfo())
	if err != nil {
		return nil, err
	}

	info := in.GetInfo()
	userID := in.GetTargetUserID()
	info.AppID = &userID

	span.AddEvent("call middleware CreateRoleUser")
	resp, err := mw.CreateRoleUser(ctx, info)
	if err != nil {
		logger.Sugar().Errorw("fail create app role user: %v", err)
		return &approleuser.CreateUserRoleUserResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &approleuser.CreateUserRoleUserResponse{
		Info: resp,
	}, nil
}
