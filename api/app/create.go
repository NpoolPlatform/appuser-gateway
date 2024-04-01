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
		app1.WithCreatedBy(&in.CreatedBy, true),
		app1.WithName(&in.Name, true),
		app1.WithLogo(&in.Logo, false),
		app1.WithDescription(&in.Description, false),
		app1.WithSignupMethods(in.GetSignupMethods(), false),
		app1.WithExtSigninMethods(in.GetExtSigninMethods(), false),
		app1.WithRecaptchaMethod(in.RecaptchaMethod, false),
		app1.WithKycEnable(in.KycEnable, false),
		app1.WithSigninVerifyEnable(in.SigninVerifyEnable, false),
		app1.WithInvitationCodeMust(in.InvitationCodeMust, false),
		app1.WithCreateInvitationCodeWhen(in.CreateInvitationCodeWhen, false),
		app1.WithMaxTypedCouponsPerOrder(in.MaxTypedCouponsPerOrder, false),
		app1.WithMaintaining(in.Maintaining, false),
		app1.WithCouponWithdrawEnable(in.CouponWithdrawEnable, false),
		app1.WithCommitButtonTargets(in.GetCommitButtonTargets(), false),
		app1.WithResetUserMethod(in.ResetUserMethod, false),
		app1.WithSettleMode(&in.SettleMode, true),
		app1.WithSettleAmountType(&in.SettleAmountType, true),
		app1.WithSettleInterval(&in.SettleInterval, true),
		app1.WithCommissionType(&in.CommissionType, true),
		app1.WithSettleBenefit(in.SettleBenefit, false),
		app1.WithStartAt(in.StartAt, false),
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
