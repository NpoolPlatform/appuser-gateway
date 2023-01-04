package app

import (
	"context"

	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	tracer "github.com/NpoolPlatform/appuser-middleware/pkg/tracer/app"
	appcrud "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/app"
	"google.golang.org/grpc/codes"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	appmgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/app"
	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/app"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateApp(ctx context.Context, in *app.UpdateAppRequest) (*app.UpdateAppResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateApp")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()
	if _, err := uuid.Parse(in.GetID()); err != nil {
		return &app.UpdateAppResponse{}, status.Error(codes.InvalidArgument, "ID is invalid")
	}

	appInfo := &appmwpb.AppReq{
		ID:                       &in.ID,
		Name:                     in.Name,
		Logo:                     in.Logo,
		Description:              in.Description,
		SignupMethods:            in.SignupMethods,
		ExtSigninMethods:         in.ExtSigninMethods,
		RecaptchaMethod:          in.RecaptchaMethod,
		KycEnable:                in.KycEnable,
		SigninVerifyEnable:       in.SigninVerifyEnable,
		InvitationCodeMust:       in.InvitationCodeMust,
		CreateInvitationCodeWhen: in.CreateInvitationCodeWhen,
		MaxTypedCouponsPerOrder:  in.MaxTypedCouponsPerOrder,
	}

	span = tracer.Trace(span, appInfo)
	span = commontracer.TraceInvoker(span, "admin", "middleware", "ExistAppConds")

	exist, err := appmgrcli.ExistAppConds(ctx, &appcrud.Conds{
		Name: &npool.StringVal{
			Value: in.GetName(),
			Op:    cruder.EQ,
		}})
	if err != nil {
		logger.Sugar().Errorw("UpdateApp", "err", err)
		return &app.UpdateAppResponse{}, status.Error(codes.Internal, err.Error())
	}

	if exist {
		logger.Sugar().Errorw("UpdateApp", "err", "app name already exists")
		return &app.UpdateAppResponse{}, status.Error(codes.AlreadyExists, "app name already exists")
	}

	span = commontracer.TraceInvoker(span, "admin", "middleware", "UpdateApp")

	info, err := appmwcli.UpdateApp(ctx, appInfo)
	if err != nil {
		logger.Sugar().Errorw("UpdateApp", "err", err)
		return &app.UpdateAppResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &app.UpdateAppResponse{
		Info: info,
	}, nil
}
