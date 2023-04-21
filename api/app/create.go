package app

import (
	"context"

	app1 "github.com/NpoolPlatform/appuser-gateway/pkg/app"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/app"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateApp(ctx context.Context, in *npool.CreateAppRequest) (*npool.CreateAppResponse, error) {
	handler, err := app1.NewHandler(
		ctx,
		app1.WithCreatedBy(&in.CreatedBy),
		app1.WithName(&in.Name),
		app1.WithLogo(&in.Logo),
		app1.WithDescription(&in.Description),
		app1.WithSignupMethods(in.GetSignupMethods()),
		app1.WithExtSigninMethods(in.GetExtSigninMethods()),
		app1.WithRecaptchaMethod(&in.RecaptchaMethod),
		app1.WithKycEnable(&in.KycEnable),
		app1.WithSigninVerifyEnable(&in.SigninVerifyEnable),
		app1.WithInvitationCodeMust(&in.InvitationCodeMust),
		app1.WithCreateInvitationCodeWhen(&in.CreateInvitationCodeWhen),
		app1.WithMaxTypedCouponsPerOrder(&in.MaxTypedCouponsPerOrder),
		app1.WithMaintaining(&in.Maintaining),
		app1.WithCommitButtonTargets(in.GetCommitButtonTargets()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateApp",
			"In", in,
			"Error", err,
		)
		return &npool.CreateAppResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateApp(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateApp",
			"In", in,
			"Error", err,
		)
		return &npool.CreateAppResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateAppResponse{
		Info: info,
	}, nil
}
