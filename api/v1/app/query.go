//nolint:nolintlint,dupl
package app

import (
	"context"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	"github.com/NpoolPlatform/message/npool/appuser/mgr/v2/recaptcha"
	"github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/app"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetApp(ctx context.Context, in *app.GetAppRequest) (*app.GetAppResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetApp")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span.SetAttributes(attribute.String("AppID", in.GetAppID()))

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("GetApp", "err", "AppID is invalid")
		return &app.GetAppResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	span = commontracer.TraceInvoker(span, "app", "middleware", "GetApp")

	resp, err := appmwcli.GetApp(ctx, in.GetAppID())
	if err != nil {
		logger.Sugar().Errorw("GetApp", "err", err)
		return &app.GetAppResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &app.GetAppResponse{
		Info: resp,
	}, nil
}

func (s *Server) GetApps(ctx context.Context, in *app.GetAppsRequest) (*app.GetAppsResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetApps")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	commontracer.TraceOffsetLimit(span, int(in.GetOffset()), int(in.GetLimit()))
	span = commontracer.TraceInvoker(span, "app", "middleware", "GetApps")

	resp, total, err := appmwcli.GetApps(ctx, in.GetLimit(), in.GetOffset())
	if err != nil {
		logger.Sugar().Errorw("GetApps", "err", err)
		return &app.GetAppsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &app.GetAppsResponse{
		Infos: resp,
		Total: total,
	}, nil
}

func (s *Server) GetUserApps(ctx context.Context, in *app.GetUserAppsRequest) (*app.GetUserAppsResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetUserApps")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	commontracer.TraceOffsetLimit(span, int(in.GetOffset()), int(in.GetLimit()))
	span.SetAttributes(attribute.String("TargetUserID", in.GetTargetUserID()))

	if _, err := uuid.Parse(in.GetTargetUserID()); err != nil {
		logger.Sugar().Errorw("GetUserApps", "err", "TargetUserID is invalid")
		return &app.GetUserAppsResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	span = commontracer.TraceInvoker(span, "app", "middleware", "GetUserApps")

	resp, total, err := appmwcli.GetUserApps(ctx, in.GetTargetUserID(), in.GetLimit(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetUserApps", "err", err)
		return &app.GetUserAppsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &app.GetUserAppsResponse{
		Infos: resp,
		Total: total,
	}, nil
}

func (s *Server) GetSignMethods(ctx context.Context, in *app.GetSignMethodsRequest) (*app.GetSignMethodsResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetSignMethods")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	signMethods := []*signmethod.SignMethod{}
	for _, val := range signmethod.SignMethodType_name {
		signMethods = append(signMethods, &signmethod.SignMethod{
			Method: val,
		})
	}
	return &app.GetSignMethodsResponse{
		Infos: signMethods,
	}, nil
}

func (s *Server) GetRecaptchas(ctx context.Context, in *app.GetRecaptchasRequest) (*app.GetRecaptchasResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetRecaptchas")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	recaptchas := []*recaptcha.Recaptcha{}
	for _, val := range recaptcha.RecaptchaType_name {
		recaptchas = append(recaptchas, &recaptcha.Recaptcha{
			Recaptcha: val,
		})
	}
	return &app.GetRecaptchasResponse{
		Infos: recaptchas,
	}, nil
}
