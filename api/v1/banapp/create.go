package banapp

import (
	"context"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	tracer "github.com/NpoolPlatform/appuser-manager/pkg/tracer/banapp"
	"google.golang.org/grpc/codes"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	banappmgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/banapp"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"

	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/banapp"
	appcrud "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/banapp"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateBanApp(ctx context.Context, in *banapp.CreateBanAppRequest) (*banapp.CreateBanAppResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateBanApp")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	err = validate(in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("CreateBanApp", "err", err)
		return nil, err
	}

	span = tracer.Trace(span, in.GetInfo())
	span = commontracer.TraceInvoker(span, "banapp", "middleware", "ExistBanAppConds")

	exist, err := banappmgrcli.ExistBanAppConds(ctx, &appcrud.Conds{
		AppID: &npool.StringVal{
			Value: in.GetInfo().GetAppID(),
			Op:    cruder.EQ,
		}})
	if err != nil {
		logger.Sugar().Errorw("CreateBanApp", "err", err)
		return &banapp.CreateBanAppResponse{}, status.Error(codes.Internal, err.Error())
	}

	if exist {
		logger.Sugar().Errorw("CreateBanApp", "err", "ban app already exists")
		return &banapp.CreateBanAppResponse{}, status.Error(codes.AlreadyExists, "ban app already exists")
	}

	span = commontracer.TraceInvoker(span, "banapp", "manager", "CreateBanApp")

	resp, err := banappmgrcli.CreateBanApp(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("CreateBanApp", "err", err)
		return &banapp.CreateBanAppResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &banapp.CreateBanAppResponse{
		Info: resp,
	}, nil
}
