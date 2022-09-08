package app

import (
	"context"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	tracer "github.com/NpoolPlatform/appuser-middleware/pkg/tracer/app"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"

	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/app"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateApp(ctx context.Context, in *app.CreateAppRequest) (*app.CreateAppResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateApp")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()
	appID := uuid.NewString()

	appInfo := &appmwpb.AppReq{
		ID:                 &appID,
		CreatedBy:          &in.CreatedBy,
		Name:               &in.Name,
		Logo:               &in.Logo,
		Description:        &in.Description,
		SignupMethods:      in.SignupMethods,
		ExtSigninMethods:   in.ExtSigninMethods,
		RecaptchaMethod:    &in.RecaptchaMethod,
		KycEnable:          &in.KycEnable,
		SigninVerifyEnable: &in.SigninVerifyEnable,
		InvitationCodeMust: &in.InvitationCodeMust,
	}
	span = tracer.Trace(span, appInfo)

	err = validate(ctx, appInfo)
	if err != nil {
		return nil, err
	}

	span = commontracer.TraceInvoker(span, "app", "middleware", "CreateApp")

	info, err := appmwcli.CreateApp(ctx, appInfo)
	if err != nil {
		logger.Sugar().Errorw("CreateApp", "err", err)
		return &app.CreateAppResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &app.CreateAppResponse{
		Info: info,
	}, nil
}
