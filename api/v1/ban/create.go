package ban

import (
	"context"

	banappuserpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/banappuser"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	banappmgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/banapp"
	banappusermgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/banappuser"
	tracer "github.com/NpoolPlatform/appuser-manager/pkg/tracer/banapp"
	tracerbanuser "github.com/NpoolPlatform/appuser-manager/pkg/tracer/banappuser"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	"google.golang.org/grpc/codes"

	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/ban"
	banapppb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/banapp"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateBanApp(ctx context.Context, in *ban.CreateBanAppRequest) (*ban.CreateBanAppResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateBanApp")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	banApp := &banapppb.BanAppReq{
		AppID:   &in.TargetAppID,
		Message: &in.Message,
	}
	err = validate(banApp)
	if err != nil {
		logger.Sugar().Errorw("CreateBanApp", "err", err)
		return nil, err
	}

	span = tracer.Trace(span, banApp)
	span = commontracer.TraceInvoker(span, "banapp", "middleware", "ExistBanAppConds")

	exist, err := banappmgrcli.ExistBanAppConds(ctx, &banapppb.Conds{
		AppID: &npool.StringVal{
			Value: in.TargetAppID,
			Op:    cruder.EQ,
		}})
	if err != nil {
		logger.Sugar().Errorw("CreateBanApp", "err", err)
		return &ban.CreateBanAppResponse{}, status.Error(codes.Internal, err.Error())
	}

	if exist {
		logger.Sugar().Errorw("CreateBanApp", "err", "ban app already exists")
		return &ban.CreateBanAppResponse{}, status.Error(codes.AlreadyExists, "ban app already exists")
	}

	span = commontracer.TraceInvoker(span, "banapp", "manager", "CreateBanApp")

	_, err = banappmgrcli.CreateBanApp(ctx, banApp)
	if err != nil {
		logger.Sugar().Errorw("CreateBanApp", "err", err)
		return &ban.CreateBanAppResponse{}, status.Error(codes.Internal, err.Error())
	}

	info, err := appmwcli.GetApp(ctx, in.TargetAppID)
	if err != nil {
		logger.Sugar().Errorw("CreateBanApp", "err", err)
		return &ban.CreateBanAppResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &ban.CreateBanAppResponse{
		Info: info,
	}, nil
}

func (s *Server) CreateBanUser(ctx context.Context,
	in *ban.CreateBanUserRequest) (*ban.CreateBanUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateBanAppUser")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	banUser := &banappuserpb.BanAppUserReq{
		AppID:   &in.AppID,
		UserID:  &in.TargetUserID,
		Message: &in.Message,
	}

	span = tracerbanuser.Trace(span, banUser)

	err = validateBanUser(banUser)
	if err != nil {
		return nil, err
	}

	span = commontracer.TraceInvoker(span, "banappuser", "middleware", "CreateBanAppUser")

	_, err = banappusermgrcli.CreateBanAppUser(ctx, banUser)
	if err != nil {
		logger.Sugar().Errorw("CreateBanUser", "err", err)
		return &ban.CreateBanUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := usermwcli.GetUser(ctx, in.GetAppID(), in.GetTargetUserID())
	if err != nil {
		logger.Sugar().Errorw("CreateBanApp", "err", err)
		return &ban.CreateBanUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &ban.CreateBanUserResponse{
		Info: resp,
	}, nil
}

func (s *Server) CreateAppBanUser(ctx context.Context,
	in *ban.CreateAppBanUserRequest) (*ban.CreateAppBanUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateAppBanUser")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	banUser := &banappuserpb.BanAppUserReq{
		AppID:   &in.TargetAppID,
		UserID:  &in.TargetUserID,
		Message: &in.Message,
	}

	span = tracerbanuser.Trace(span, banUser)

	err = validateBanUser(banUser)
	if err != nil {
		logger.Sugar().Errorw("CreateAppBanUser", "err", err)
		return nil, err
	}

	span = commontracer.TraceInvoker(span, "banappuser", "middleware", "CreateBanAppUser")

	_, err = banappusermgrcli.CreateBanAppUser(ctx, banUser)
	if err != nil {
		logger.Sugar().Errorw("CreateAppBanUser", "err", err)
		return &ban.CreateAppBanUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := usermwcli.GetUser(ctx, in.GetTargetAppID(), in.GetTargetUserID())
	if err != nil {
		logger.Sugar().Errorw("CreateBanApp", "err", err)
		return &ban.CreateAppBanUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &ban.CreateAppBanUserResponse{
		Info: info,
	}, nil
}
